package main

import (
	"archive/zip"
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/rirolli/work_search/internal/config"
	"github.com/rirolli/work_search/internal/fetch"
	"github.com/rirolli/work_search/internal/report"
	"github.com/rirolli/work_search/internal/state"
)

// Esecuzione completa con una fonte finta e la configurazione reale del repository.
// Con WS_SAMPLE_OUT=<cartella> i report vengono lasciati lì per ispezionarli.
func TestEndToEnd(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"jobs":[
		 {"id":1,"title":"Machine Learning Engineer – SAR","absolute_url":"https://example.com/1","first_published":"`+time.Now().UTC().Format(time.RFC3339)+`",
		  "content":"Earth observation, PyTorch, onboard AI","location":{"name":"Rome, Italy"}},
		 {"id":2,"title":"Engineering Manager, Computer Vision","absolute_url":"https://example.com/2","first_published":"2026-10-01T00:00:00Z",
		  "content":"geospatial","location":{"name":"Berlin, Germany"}},
		 {"id":3,"title":"Office Manager","absolute_url":"https://example.com/3","location":{"name":"Rome"}},
		 {"id":4,"title":"Computer Vision Engineer","absolute_url":"https://example.com/4","location":{"name":"San Francisco, USA"}},
		 {"id":5,"title":"ML Intern","absolute_url":"https://example.com/5","location":{"name":"Rome"}}]}`)
	}))
	defer srv.Close()

	cfg, err := config.Load(filepath.Join("..", "..", "config"))
	if err != nil {
		t.Fatal(err)
	}
	cfg.Companies = []config.Company{{Name: "Acme Space", Sector: "Spazio – EO", MacroSector: "Aerospazio", Aerospace: true, Rating: 8,
		Sources: []config.Source{{Type: "greenhouse", ID: "acme", BaseURL: srv.URL, Verified: true}}}}
	cfg.Searches = nil

	raw, statuses := collect(context.Background(), fetch.New(5*time.Second), cfg, "", 2)
	if len(raw) != 5 || len(statuses) != 1 || !statuses[0].OK {
		t.Fatalf("raw=%d statuses=%+v", len(raw), statuses)
	}
	now := time.Now().UTC()
	kept, _ := process(raw, cfg, now)
	if len(kept) != 2 {
		for _, k := range kept {
			t.Logf("tenuta: %s (%d)", k.Title, k.Score)
		}
		t.Fatalf("attese 2 offerte tenute, ottenute %d", len(kept))
	}
	if kept[0].Location != "Rome, Italy" || kept[0].Priority != "Alta" {
		t.Errorf("prima offerta inattesa: %+v", kept[0])
	}

	st := &state.State{Jobs: map[string]*state.Entry{
		"sparita":  {Company: "Vecchia", Title: "AI Engineer", Origin: "Acme Space · Greenhouse", LastSeen: now.Add(-12 * time.Hour)},
		"in-pausa": {Company: "Altra", Title: "ML Engineer", Origin: "Altra · Lever", LastSeen: now.Add(-12 * time.Hour)},
	}}
	st.Merge(kept, now, 24*time.Hour, map[string]bool{"Altra · Lever": true}, 7*24*time.Hour)
	if _, ok := st.Jobs["sparita"]; ok {
		t.Error("un'offerta non più presente deve essere rimossa dallo stato")
	}
	if _, ok := st.Jobs["in-pausa"]; !ok {
		t.Error("le offerte di una fonte in errore devono restare in memoria")
	}
	if !kept[0].New || len(st.Jobs) != 3 {
		t.Fatalf("merge: new=%v voci=%d", kept[0].New, len(st.Jobs))
	}

	out := os.Getenv("WS_SAMPLE_OUT")
	if out == "" {
		out = t.TempDir()
	}
	d := report.Data{Generated: now, Jobs: kept, Sources: statuses, Profile: cfg.Profile, Fetched: len(raw)}
	if err := report.WriteHTML(filepath.Join(out, "index.html"), d); err != nil {
		t.Fatal(err)
	}
	if err := report.WriteXLSX(filepath.Join(out, "offerte.xlsx"), d); err != nil {
		t.Fatal(err)
	}
	if err := writeJSON(filepath.Join(out, "offerte.json"), d); err != nil {
		t.Fatal(err)
	}

	// Ogni parte XML dell'xlsx deve essere ben formata.
	zr, err := zip.OpenReader(filepath.Join(out, "offerte.xlsx"))
	if err != nil {
		t.Fatal(err)
	}
	defer zr.Close()
	for _, f := range zr.File {
		rc, _ := f.Open()
		dec := xml.NewDecoder(rc)
		for {
			if _, err := dec.Token(); err == io.EOF {
				break
			} else if err != nil {
				t.Fatalf("%s: XML non valido: %v", f.Name, err)
			}
		}
		rc.Close()
	}
	readme := filepath.Join(out, "README.md")
	os.WriteFile(readme, []byte("# Titolo\n\nprima\n"+report.ReadmeStart+"\nvecchio\n"+report.ReadmeEnd+"\n\ndopo\n"), 0o644)
	for i := 0; i < 2; i++ { // due volte: la sezione va sostituita, non duplicata
		if err := report.UpdateReadme(readme, d); err != nil {
			t.Fatal(err)
		}
	}
	md, _ := os.ReadFile(readme)
	if s := string(md); strings.Contains(s, "vecchio") || strings.Count(s, report.ReadmeStart) != 1 ||
		!strings.Contains(s, "](https://example.com/1)") || !strings.HasSuffix(s, "dopo\n") {
		t.Errorf("README non aggiornato correttamente:\n%s", s)
	}

	html, _ := os.ReadFile(filepath.Join(out, "index.html"))
	if !strings.Contains(string(html), "Machine Learning Engineer – SAR") || !strings.Contains(string(html), "NUOVA") {
		t.Error("l'HTML non contiene l'offerta attesa")
	}
}
