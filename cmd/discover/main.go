// Comando discover: cerca le API pubbliche delle pagine carriere di un elenco di
// aziende provando i nomi probabili su ogni tipo di ATS supportato.
// Stampa un report markdown con le fonti funzionanti e lo snippet da copiare
// in config/companies.json.
//
//	go run ./cmd/discover -candidates config/candidates.json > discover.md
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/rirolli/work_search/internal/config"
	"github.com/rirolli/work_search/internal/fetch"
	"github.com/rirolli/work_search/internal/sources"
)

type candidate struct {
	Company string   `json:"company"`
	Slugs   []string `json:"slugs,omitempty"`   // in aggiunta a quelli ricavati dal nome
	Domains []string `json:"domains,omitempty"` // per teamtailor, pinpoint, successfactors
	Workday bool     `json:"workday,omitempty"` // prova anche *.myworkdayjobs.com
	Note    string   `json:"note,omitempty"`
}

var nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

func slugs(c candidate) []string {
	seen := map[string]bool{}
	var out []string
	add := func(s string) {
		if s != "" && !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	for _, s := range c.Slugs {
		add(s)
	}
	n := strings.ToLower(c.Company)
	add(nonAlnum.ReplaceAllString(n, ""))
	add(strings.Trim(nonAlnum.ReplaceAllString(n, "-"), "-"))
	return out
}

func trials(c candidate) []config.Source {
	var out []config.Source
	for _, s := range slugs(c) {
		out = append(out,
			config.Source{Type: "greenhouse", ID: s},
			config.Source{Type: "greenhouse", ID: s, Region: "eu"},
			config.Source{Type: "bamboohr", ID: s},
			config.Source{Type: "breezy", ID: s},
			config.Source{Type: "lever", ID: s},
			config.Source{Type: "lever", ID: s, Region: "eu"},
			config.Source{Type: "ashby", ID: s},
			config.Source{Type: "workable", ID: s},
			config.Source{Type: "recruitee", ID: s},
			config.Source{Type: "personio", ID: s, Domain: "de"},
			config.Source{Type: "personio", ID: s, Domain: "com"},
			config.Source{Type: "smartrecruiters", ID: s, Pages: 1},
			config.Source{Type: "pinpoint", ID: s},
			config.Source{Type: "teamtailor", Domain: s + ".teamtailor.com"},
		)
		if c.Workday {
			for _, wd := range []string{"wd3", "wd103", "wd1", "wd5"} {
				for _, site := range []string{"Careers", "External", "careers", strings.ToUpper(s[:1]) + s[1:]} {
					out = append(out, config.Source{Type: "workday", Domain: s + "." + wd + ".myworkdayjobs.com", Tenant: s, Site: site, Pages: 1})
				}
			}
		}
	}
	domains := append([]string{}, c.Domains...)
	for _, s := range slugs(c) {
		if !strings.Contains(s, "-") {
			domains = append(domains, "careers."+s+".com", "jobs."+s+".com")
		}
	}
	for _, d := range domains {
		out = append(out,
			config.Source{Type: "teamtailor", Domain: d},
			config.Source{Type: "pinpoint", Domain: d},
			config.Source{Type: "successfactors", Domain: d},
		)
	}
	return out
}

type hit struct {
	Src    config.Source
	Count  int
	Sample []string
}

func main() {
	path := flag.String("candidates", "config/candidates.json", "elenco delle aziende da cercare")
	only := flag.String("only", "", "solo le aziende il cui nome contiene questo testo")
	parallel := flag.Int("parallel", 16, "richieste in parallelo")
	flag.Parse()

	b, err := os.ReadFile(*path)
	if err != nil {
		log.Fatal(err)
	}
	var file struct {
		Candidates []candidate `json:"candidates"`
	}
	if err := json.Unmarshal(b, &file); err != nil {
		log.Fatalf("%s: %v", *path, err)
	}

	client := fetch.New(15 * time.Second)
	client.Retries = 0
	sem := make(chan struct{}, *parallel)
	type result struct {
		c    candidate
		hits []hit
	}
	results := make([]result, len(file.Candidates))
	var wg sync.WaitGroup
	for i, c := range file.Candidates {
		if *only != "" && !strings.Contains(strings.ToLower(c.Company), strings.ToLower(*only)) {
			continue
		}
		results[i].c = c
		var mu sync.Mutex
		for _, src := range trials(c) {
			wg.Add(1)
			go func(i int, src config.Source) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()
				ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
				defer cancel()
				jobs, err := sources.ATS[src.Type](ctx, client, src, results[i].c.Company)
				if err != nil || len(jobs) == 0 {
					return
				}
				h := hit{Src: src, Count: len(jobs)}
				for k := 0; k < len(jobs) && k < 3; k++ {
					h.Sample = append(h.Sample, jobs[k].Title)
				}
				mu.Lock()
				results[i].hits = append(results[i].hits, h)
				mu.Unlock()
			}(i, src)
		}
	}
	wg.Wait()

	var found, missing []string
	fmt.Print("# Scoperta delle pagine carriere\n\n")
	fmt.Printf("_%s_\n\n", time.Now().UTC().Format("2006-01-02 15:04 UTC"))
	for _, r := range results {
		if r.c.Company == "" {
			continue
		}
		if len(r.hits) == 0 {
			missing = append(missing, r.c.Company)
			continue
		}
		found = append(found, r.c.Company)
		sort.Slice(r.hits, func(a, b int) bool { return r.hits[a].Count > r.hits[b].Count })
		fmt.Printf("## %s\n\n", r.c.Company)
		for _, h := range r.hits {
			snip, _ := json.Marshal(clean(h.Src))
			fmt.Printf("- **%d offerte** · `%s` · es.: %s\n", h.Count, snip, strings.Join(h.Sample, " · "))
		}
		fmt.Println()
	}
	fmt.Printf("**Trovate (%d):** %s\n\n", len(found), strings.Join(found, ", "))
	fmt.Printf("**Nessuna API trovata (%d):** %s\n", len(missing), strings.Join(missing, ", "))
}

// clean lascia solo i campi utili allo snippet.
func clean(s config.Source) map[string]any {
	m := map[string]any{"type": s.Type, "verified": true}
	for k, v := range map[string]string{"id": s.ID, "region": s.Region, "domain": s.Domain, "tenant": s.Tenant, "site": s.Site} {
		if v != "" {
			m[k] = v
		}
	}
	return m
}
