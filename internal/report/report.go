// Package report produce la pagina HTML e il file Excel.
package report

import (
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"
	_ "time/tzdata" // fuso Europe/Rome anche dove manca tzdata di sistema

	"github.com/rirolli/work_search/internal/config"
	"github.com/rirolli/work_search/internal/model"
	"github.com/rirolli/work_search/internal/state"
	"github.com/rirolli/work_search/internal/xlsx"
)

// Data è tutto ciò che serve ai report.
type Data struct {
	Generated time.Time
	Jobs      []model.Job
	Gone      []state.Entry
	Sources   []model.SourceStatus
	Profile   config.Profile
	Fetched   int
}

var rome = func() *time.Location {
	l, err := time.LoadLocation("Europe/Rome")
	if err != nil {
		return time.UTC
	}
	return l
}()

func day(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(rome).Format("2006-01-02")
}

func published(j model.Job) string {
	d := day(j.Published)
	switch {
	case d == "" && j.PublishedNote != "":
		return j.PublishedNote
	case d == "":
		return "non indicata"
	case j.PublishedNote != "" && !strings.HasPrefix(strings.ToLower(j.PublishedNote), "posted"):
		return d + " (" + j.PublishedNote + ")"
	}
	return d
}

func yesNo(b bool) string {
	if b {
		return "Sì"
	}
	return ""
}

func (d Data) count(p string) int {
	n := 0
	for _, j := range d.Jobs {
		if j.Priority == p || (p == "new" && j.New) {
			n++
		}
	}
	return n
}

//go:embed template.html
var pageTemplate string

// WriteHTML scrive la pagina autonoma (nessuna risorsa esterna).
func WriteHTML(path string, d Data) error {
	funcs := template.FuncMap{
		"day":       day,
		"published": published,
		"local":     func(t time.Time) string { return t.In(rome).Format("02/01/2006 15:04") },
		"join":      strings.Join,
		"count":     d.count,
		"failed": func() int {
			n := 0
			for _, s := range d.Sources {
				if !s.OK && !s.Skipped {
					n++
				}
			}
			return n
		},
	}
	t, err := template.New("page").Funcs(funcs).Parse(pageTemplate)
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := t.Execute(f, d); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func prioStyle(p string) int {
	switch p {
	case "Alta":
		return xlsx.High
	case "Media":
		return xlsx.Medium
	}
	return xlsx.Low
}

// WriteXLSX scrive il file Excel con i fogli Offerte, Non più visibili, Fonti, Info.
func WriteXLSX(path string, d Data) error {
	h := func(cols ...string) []xlsx.Cell {
		row := make([]xlsx.Cell, len(cols))
		for i, c := range cols {
			row[i] = xlsx.T(c, xlsx.Header)
		}
		return row
	}

	jobs := xlsx.Sheet{Name: "Offerte", FreezeRow: true, AutoFilter: true,
		Widths: []float64{9, 7, 7, 24, 44, 26, 9, 22, 24, 11, 16, 16, 14, 20, 9, 16, 40, 10, 30}}
	jobs.Rows = append(jobs.Rows, h("Priorità", "Score", "Nuova", "Azienda", "Ruolo", "Sede", "Remoto",
		"Macro-settore", "Settore", "Aerospaziale?", "Data pubblicazione", "Prima vista", "Data scadenza",
		"Retribuzione", "≥ soglia?", "Fonte", "Perché (termini trovati)", "Rating azienda", "Link"))
	for _, j := range d.Jobs {
		newCell := xlsx.T("", xlsx.Default)
		if j.New {
			newCell = xlsx.T("NUOVA", xlsx.New)
		}
		deadline := j.Deadline
		if deadline == "" {
			deadline = "non indicata"
		}
		rating := xlsx.T("", xlsx.Default)
		if j.CompanyRating > 0 {
			rating = xlsx.N(float64(j.CompanyRating), xlsx.Default)
		}
		jobs.Rows = append(jobs.Rows, []xlsx.Cell{
			xlsx.T(j.Priority, prioStyle(j.Priority)),
			xlsx.N(float64(j.Score), xlsx.Default),
			newCell,
			xlsx.T(j.Company, xlsx.Bold),
			xlsx.T(j.Title, xlsx.Wrap),
			xlsx.T(j.Location, xlsx.Wrap),
			xlsx.T(yesNo(j.Remote), xlsx.Default),
			xlsx.T(j.MacroSector, xlsx.Default),
			xlsx.T(j.Sector, xlsx.Wrap),
			xlsx.T(j.Aerospace, xlsx.Default),
			xlsx.T(published(j), xlsx.Wrap),
			xlsx.T(day(j.FirstSeen), xlsx.Default),
			xlsx.T(deadline, xlsx.Default),
			xlsx.T(j.Salary, xlsx.Wrap),
			xlsx.T(j.SalaryOK, xlsx.Default),
			xlsx.T(j.Source, xlsx.Default),
			xlsx.T(strings.Join(j.Matched, ", "), xlsx.Wrap),
			rating,
			xlsx.L("Apri annuncio", j.URL),
		})
	}

	gone := xlsx.Sheet{Name: "Non più visibili", FreezeRow: true, AutoFilter: true,
		Widths: []float64{24, 44, 26, 9, 7, 14, 14, 16}}
	gone.Rows = append(gone.Rows, h("Azienda", "Ruolo", "Sede", "Priorità", "Score", "Prima vista", "Ultima vista", "Link"))
	for _, e := range d.Gone {
		gone.Rows = append(gone.Rows, []xlsx.Cell{
			xlsx.T(e.Company, xlsx.Bold), xlsx.T(e.Title, xlsx.Wrap), xlsx.T(e.Location, xlsx.Wrap),
			xlsx.T(e.Priority, prioStyle(e.Priority)), xlsx.N(float64(e.Score), xlsx.Default),
			xlsx.T(day(e.FirstSeen), xlsx.Default), xlsx.T(day(e.LastSeen), xlsx.Default),
			xlsx.L("Apri (potrebbe essere chiuso)", e.URL),
		})
	}

	src := xlsx.Sheet{Name: "Fonti", FreezeRow: true, AutoFilter: true, Widths: []float64{50, 16, 12, 10, 10, 12, 60}}
	src.Rows = append(src.Rows, h("Fonte", "Tipo", "Esito", "Scaricate", "Tenute", "Verificata a mano", "Errore / nota"))
	for _, s := range d.Sources {
		esito, st := "OK", xlsx.High
		switch {
		case s.Skipped:
			esito, st = "Saltata", xlsx.Low
		case !s.OK:
			esito, st = "Errore", xlsx.Medium
		}
		src.Rows = append(src.Rows, []xlsx.Cell{
			xlsx.T(s.Name, xlsx.Default), xlsx.T(s.Type, xlsx.Default), xlsx.T(esito, st),
			xlsx.N(float64(s.Fetched), xlsx.Default), xlsx.N(float64(s.Kept), xlsx.Default),
			xlsx.T(yesNo(s.Verified), xlsx.Default), xlsx.T(s.Error, xlsx.Wrap),
		})
	}

	p := d.Profile
	info := xlsx.Sheet{Name: "Info", Widths: []float64{34, 90}}
	kv := func(k, v string) []xlsx.Cell { return []xlsx.Cell{xlsx.T(k, xlsx.Bold), xlsx.T(v, xlsx.Wrap)} }
	info.Rows = [][]xlsx.Cell{
		h("Voce", "Valore"),
		kv("Generato", d.Generated.In(rome).Format("02/01/2006 15:04")+" (Europe/Rome)"),
		kv("Offerte scaricate / tenute", fmt.Sprintf("%d / %d", d.Fetched, len(d.Jobs))),
		kv("Priorità", fmt.Sprintf("Alta ≥ %d, Media ≥ %d, Bassa sotto; soglia minima per comparire: %d", p.PriorityHigh, p.PriorityMedium, p.MinScore)),
		kv("Punteggio", fmt.Sprintf("parola chiave nel titolo = peso×2, nella descrizione = peso; + peso della sede; + rating azienda×%d; + %d se pubblicata negli ultimi %d giorni", p.CompanyRatingFactor, p.RecencyBonus, p.RecencyDays)),
		kv("Nuova", fmt.Sprintf("vista per la prima volta nelle ultime %d ore", p.NewWithinHours)),
		kv("Soglia retribuzione", fmt.Sprintf("%.0f EUR/anno; \"≥ soglia?\" vale solo quando la fonte dà un importo annuo in EUR, altrimenti n.d.", p.SalaryMinEUR)),
		kv("Data pubblicazione", "presa dall'API della fonte; Workday dà solo date relative (\"30+ Days Ago\"); Jooble dà la data di aggiornamento. \"Prima vista\" è la data in cui questa routine l'ha trovata."),
		kv("Data scadenza", "quasi nessuna API la espone: va controllata sulla pagina dell'annuncio."),
		kv("Configurazione", "config/companies.json (aziende e fonti), config/searches.json (ricerche sugli aggregatori), config/profile.json (parole chiave, sedi, soglie)."),
	}
	return xlsx.Write(path, []xlsx.Sheet{jobs, gone, src, info})
}
