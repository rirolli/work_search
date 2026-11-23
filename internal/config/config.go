// Package config carica i file JSON in config/.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Source è una fonte di offerte di una singola azienda (pagina carriere / ATS).
type Source struct {
	Type     string   `json:"type"`               // greenhouse, lever, ashby, personio, recruitee, workable, smartrecruiters, workday
	ID       string   `json:"id"`                 // board token / account / slug
	Region   string   `json:"region,omitempty"`   // lever: "eu" per api.eu.lever.co
	Domain   string   `json:"domain,omitempty"`   // personio: "de" o "com"; workday: host (es. ag.wd3.myworkdayjobs.com)
	Tenant   string   `json:"tenant,omitempty"`   // workday
	Site     string   `json:"site,omitempty"`     // workday
	Queries  []string `json:"queries,omitempty"`  // workday: testi di ricerca
	Pages    int      `json:"pages,omitempty"`    // workday/smartrecruiters: pagine massime
	BaseURL  string   `json:"base_url,omitempty"` // sovrascrive l'host dell'API (test o domini custom)
	Verified bool     `json:"verified"`           // true se l'endpoint è stato controllato a mano
	Note     string   `json:"note,omitempty"`
}

// Company descrive un'azienda da monitorare.
type Company struct {
	Name        string   `json:"name"`
	Aliases     []string `json:"aliases,omitempty"` // nomi con cui compare negli aggregatori
	Sector      string   `json:"sector"`
	MacroSector string   `json:"macro_sector"`
	Aerospace   bool     `json:"aerospace"`
	Rating      int      `json:"rating"` // 1-10, giudizio sull'azienda
	Country     string   `json:"country"`
	Exclude     bool     `json:"exclude,omitempty"` // es. datore di lavoro attuale
	Notes       string   `json:"notes,omitempty"`
	Sources     []Source `json:"sources,omitempty"`
}

// Search è una ricerca per parole chiave su un aggregatore.
type Search struct {
	Type       string `json:"type"` // adzuna, jooble, arbeitnow, remotive
	Query      string `json:"query,omitempty"`
	Where      string `json:"where,omitempty"`
	Country    string `json:"country,omitempty"` // adzuna: it, de, fr, nl, gb, es, at, ...
	MaxDaysOld int    `json:"max_days_old,omitempty"`
	Pages      int    `json:"pages,omitempty"`
	Disabled   bool   `json:"disabled,omitempty"`
	Note       string `json:"note,omitempty"`
}

// Weighted è un termine con peso.
type Weighted struct {
	Term   string `json:"term"`
	Weight int    `json:"weight"`
	// Weak (solo per le sedi): il termine aggiunge punti ma non salva un'offerta
	// la cui sede contiene anche un termine escluso (es. "Remote – USA").
	Weak bool `json:"weak,omitempty"`
}

// Profile contiene i criteri di filtro e di punteggio.
type Profile struct {
	TitleRequiredAny    []string   `json:"title_required_any"`
	TitleExclude        []string   `json:"title_exclude"`
	Keywords            []Weighted `json:"keywords"`
	Locations           []Weighted `json:"locations"`
	LocationExclude     []string   `json:"location_exclude"`
	MinScore            int        `json:"min_score"`
	PriorityHigh        int        `json:"priority_high"`
	PriorityMedium      int        `json:"priority_medium"`
	SalaryMinEUR        float64    `json:"salary_min_eur"`
	CompanyRatingFactor int        `json:"company_rating_factor"`
	RecencyDays         int        `json:"recency_days"`
	RecencyBonus        int        `json:"recency_bonus"`
	MaxAgeDays          int        `json:"max_age_days"` // scarta offerte con data nota più vecchia di così (0 = nessun limite)
	NewWithinHours      int        `json:"new_within_hours"`
	DisappearedDays     int        `json:"disappeared_days"`
	StateRetentionDays  int        `json:"state_retention_days"`
}

// Config raccoglie tutti i file di configurazione.
type Config struct {
	Companies []Company
	Searches  []Search
	Profile   Profile
}

func readJSON(path string, v any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, v); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	return nil
}

// Load legge companies.json, searches.json e profile.json da dir.
func Load(dir string) (*Config, error) {
	var c struct {
		Companies []Company `json:"companies"`
	}
	var s struct {
		Searches []Search `json:"searches"`
	}
	var p Profile
	if err := readJSON(filepath.Join(dir, "companies.json"), &c); err != nil {
		return nil, err
	}
	if err := readJSON(filepath.Join(dir, "searches.json"), &s); err != nil {
		return nil, err
	}
	if err := readJSON(filepath.Join(dir, "profile.json"), &p); err != nil {
		return nil, err
	}
	p.defaults()
	return &Config{Companies: c.Companies, Searches: s.Searches, Profile: p}, nil
}

func (p *Profile) defaults() {
	if p.PriorityHigh == 0 {
		p.PriorityHigh = 60
	}
	if p.PriorityMedium == 0 {
		p.PriorityMedium = 40
	}
	if p.CompanyRatingFactor == 0 {
		p.CompanyRatingFactor = 2
	}
	if p.NewWithinHours == 0 {
		p.NewWithinHours = 24
	}
	if p.DisappearedDays == 0 {
		p.DisappearedDays = 14
	}
	if p.StateRetentionDays == 0 {
		p.StateRetentionDays = 90
	}
}
