// Comando jobsearch: scarica le offerte, le filtra rispetto al profilo e
// produce docs/index.html, docs/offerte.xlsx e docs/offerte.json.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/rirolli/work_search/internal/config"
	"github.com/rirolli/work_search/internal/fetch"
	"github.com/rirolli/work_search/internal/model"
	"github.com/rirolli/work_search/internal/report"
	"github.com/rirolli/work_search/internal/score"
	"github.com/rirolli/work_search/internal/sources"
	"github.com/rirolli/work_search/internal/state"
)

func main() {
	cfgDir := flag.String("config", "config", "cartella con companies.json, searches.json, profile.json")
	outDir := flag.String("out", "docs", "cartella di output (HTML, XLSX, JSON)")
	statePath := flag.String("state", "data/state.json", "file di stato delle offerte già viste")
	only := flag.String("only", "", "esegui solo le aziende/ricerche il cui nome contiene questo testo (debug)")
	parallel := flag.Int("parallel", 6, "richieste in parallelo")
	readme := flag.String("readme", "README.md", "README in cui aggiornare la sezione Risultati (vuoto = non toccarlo)")
	flag.Parse()

	cfg, err := config.Load(*cfgDir)
	if err != nil {
		log.Fatalf("configurazione: %v", err)
	}
	st, err := state.Load(*statePath)
	if err != nil {
		log.Fatalf("stato: %v", err)
	}

	now := time.Now().UTC()
	client := fetch.New(45 * time.Second)
	raw, statuses := collect(context.Background(), client, cfg, *only, *parallel)

	kept, discarded := process(raw, cfg, now)
	log.Printf("offerte scaricate: %d, tenute: %d, scartate: %d", len(raw), len(kept), discarded)

	failed := map[string]bool{}
	for _, s := range statuses {
		if !s.OK && !s.Skipped {
			failed[s.Name] = true
		}
	}
	st.Merge(kept, now, time.Duration(cfg.Profile.NewWithinHours)*time.Hour,
		failed, time.Duration(cfg.Profile.StateRetentionDays)*24*time.Hour)

	// Conteggio delle offerte tenute per fonte.
	keptBy := map[string]int{}
	for _, j := range kept {
		keptBy[j.Origin]++
	}
	for i := range statuses {
		statuses[i].Kept = keptBy[statuses[i].Name]
	}

	data := report.Data{
		Generated: now,
		Jobs:      kept,
		Sources:   statuses,
		Profile:   cfg.Profile,
		Fetched:   len(raw),
	}
	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		log.Fatal(err)
	}
	must(report.WriteHTML(filepath.Join(*outDir, "index.html"), data))
	must(report.WriteXLSX(filepath.Join(*outDir, "offerte.xlsx"), data))
	must(writeJSON(filepath.Join(*outDir, "offerte.json"), data))
	must(st.Save(*statePath))
	if *readme != "" {
		must(report.UpdateReadme(*readme, data))
	}

	log.Printf("fonti: %d, in errore: %d. Report in %s", len(statuses), len(failed), *outDir)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type task struct {
	name, typ string
	verified  bool
	run       func(ctx context.Context) ([]model.Job, error)
	meta      *config.Company
}

func collect(ctx context.Context, c *fetch.Client, cfg *config.Config, only string, parallel int) ([]model.Job, []model.SourceStatus) {
	var tasks []task
	for i := range cfg.Companies {
		comp := &cfg.Companies[i]
		if comp.Exclude || (only != "" && !strings.Contains(strings.ToLower(comp.Name), strings.ToLower(only))) {
			continue
		}
		for _, src := range comp.Sources {
			src := src
			f, ok := sources.ATS[src.Type]
			if !ok {
				log.Printf("tipo di fonte sconosciuto %q per %s", src.Type, comp.Name)
				continue
			}
			name := comp.Name
			tasks = append(tasks, task{
				name: name + " · " + sourceLabel(src.Type), typ: src.Type, verified: src.Verified, meta: comp,
				run: func(ctx context.Context) ([]model.Job, error) { return f(ctx, c, src, name) },
			})
		}
	}
	for _, s := range cfg.Searches {
		s := s
		if s.Disabled {
			continue
		}
		label := strings.TrimSpace(s.Query + " @ " + s.Where + " " + s.Country)
		if only != "" && !strings.Contains(strings.ToLower(label), strings.ToLower(only)) {
			continue
		}
		f, ok := sources.Aggregators[s.Type]
		if !ok {
			log.Printf("aggregatore sconosciuto %q", s.Type)
			continue
		}
		tasks = append(tasks, task{
			name: sourceLabel(s.Type) + ": " + label, typ: s.Type, verified: true,
			run: func(ctx context.Context) ([]model.Job, error) { return f(ctx, c, s) },
		})
	}

	results := make([][]model.Job, len(tasks))
	statuses := make([]model.SourceStatus, len(tasks))
	sem := make(chan struct{}, max(1, parallel))
	var wg sync.WaitGroup
	for i, t := range tasks {
		wg.Add(1)
		go func(i int, t task) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			tctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
			defer cancel()
			start := time.Now()
			jobs, err := t.run(tctx)
			st := model.SourceStatus{Name: t.name, Type: t.typ, Verified: t.verified, Fetched: len(jobs), Millis: time.Since(start).Milliseconds()}
			switch {
			case errors.Is(err, sources.ErrMissingKey):
				st.Skipped, st.Error = true, err.Error()
			case err != nil:
				st.Error = err.Error()
				st.OK = len(jobs) > 0
			default:
				st.OK = true
			}
			for k := range jobs {
				jobs[k].Origin = t.name
				if t.meta != nil {
					applyMeta(&jobs[k], t.meta)
				}
			}
			results[i], statuses[i] = jobs, st
			log.Printf("%-60s %4d offerte %s", t.name, len(jobs), st.Error)
		}(i, t)
	}
	wg.Wait()

	var all []model.Job
	for _, r := range results {
		all = append(all, r...)
	}
	return all, statuses
}

// Il nome della fonte nel conteggio deve coincidere con model.Job.Source.
func sourceLabel(t string) string {
	switch t {
	case "smartrecruiters":
		return "SmartRecruiters"
	default:
		return strings.ToUpper(t[:1]) + t[1:]
	}
}

func applyMeta(j *model.Job, c *config.Company) {
	j.Company = c.Name
	j.Sector, j.MacroSector, j.CompanyRating, j.Country = c.Sector, c.MacroSector, c.Rating, c.Country
	if c.Aerospace {
		j.Aerospace = "Sì"
	} else {
		j.Aerospace = "No"
	}
}

// matchCompany riconosce l'azienda di un'offerta da aggregatore.
func matchCompany(name string, companies []config.Company) *config.Company {
	n := score.Norm(name)
	if strings.TrimSpace(n) == "" {
		return nil
	}
	for i := range companies {
		c := &companies[i]
		for _, a := range append([]string{c.Name}, c.Aliases...) {
			if score.Has(n, a) {
				return c
			}
		}
	}
	return nil
}

func process(raw []model.Job, cfg *config.Config, now time.Time) ([]model.Job, int) {
	seenKey := map[string]bool{}
	byTitle := map[string]int{} // azienda+titolo -> indice in kept
	var kept []model.Job
	discarded := 0
	for _, j := range raw {
		if j.Title == "" || j.URL == "" || seenKey[j.Key] {
			continue
		}
		seenKey[j.Key] = true
		if j.Sector == "" { // offerta da aggregatore
			if c := matchCompany(j.Company, cfg.Companies); c != nil {
				if c.Exclude {
					discarded++
					continue
				}
				applyMeta(&j, c)
			}
		}
		ok, _ := score.Evaluate(&j, cfg.Profile, now)
		if !ok {
			discarded++
			continue
		}
		// Stessa azienda e stesso titolo (più sedi o più fonti): una sola riga,
		// con le sedi unite. Vince la prima vista (le fonti ATS precedono gli aggregatori).
		dk := score.Norm(j.Company) + "|" + score.Norm(j.Title)
		if i, dup := byTitle[dk]; dup {
			k := &kept[i]
			if j.Location != "" && !strings.Contains(k.Location, j.Location) && len(k.Location) < 200 {
				k.Location += " / " + j.Location
			}
			if j.Score > k.Score {
				k.Score, k.Priority, k.Matched = j.Score, j.Priority, j.Matched
			}
			continue
		}
		byTitle[dk] = len(kept)
		kept = append(kept, j)
	}
	score.Sort(kept)
	return kept, discarded
}

func writeJSON(path string, d report.Data) error {
	b, err := json.MarshalIndent(struct {
		Generated time.Time            `json:"generated"`
		Jobs      []model.Job          `json:"jobs"`
		Sources   []model.SourceStatus `json:"sources"`
	}{d.Generated, d.Jobs, d.Sources}, "", " ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, b, 0o644); err != nil {
		return fmt.Errorf("scrittura %s: %w", path, err)
	}
	return nil
}
