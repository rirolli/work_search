// Package model contiene i tipi condivisi tra fonti, punteggio e report.
package model

import "time"

// Job è un'offerta normalizzata, indipendente dalla fonte.
type Job struct {
	Key        string `json:"key"` // chiave univoca: fonte:azienda:id
	Company    string `json:"company"`
	Title      string `json:"title"`
	Location   string `json:"location"`
	Remote     bool   `json:"remote"`
	Department string `json:"department,omitempty"`
	URL        string `json:"url"`
	Source     string `json:"source"`
	Origin     string `json:"-"` // nome della fonte/ricerca che l'ha prodotta (per i conteggi)

	Published     time.Time `json:"published"`                // zero se non nota
	PublishedNote string    `json:"published_note,omitempty"` // es. "Posted 30+ Days Ago"
	Deadline      string    `json:"deadline,omitempty"`

	Salary         string  `json:"salary,omitempty"` // testo da mostrare
	SalaryMin      float64 `json:"-"`
	SalaryMax      float64 `json:"-"`
	SalaryCurrency string  `json:"-"`
	SalaryYearly   bool    `json:"-"`

	Description string `json:"-"` // testo semplice, usato solo per il punteggio

	// Metadati dell'azienda (da config/companies.json, se riconosciuta).
	Sector        string `json:"sector,omitempty"`
	MacroSector   string `json:"macro_sector,omitempty"`
	Aerospace     string `json:"aerospace,omitempty"` // "Sì" / "No" / ""
	CompanyRating int    `json:"company_rating,omitempty"`
	Country       string `json:"country,omitempty"`

	// Calcolati.
	Score     int       `json:"score"`
	Priority  string    `json:"priority"`
	Matched   []string  `json:"matched,omitempty"`
	SalaryOK  string    `json:"salary_ok"` // "Sì" / "No" / "n.d."
	FirstSeen time.Time `json:"first_seen"`
	New       bool      `json:"new"`
}

// SourceStatus descrive l'esito di una singola fonte in un'esecuzione.
type SourceStatus struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	OK       bool   `json:"ok"`
	Skipped  bool   `json:"skipped"`
	Fetched  int    `json:"fetched"`
	Kept     int    `json:"kept"`
	Error    string `json:"error,omitempty"`
	Verified bool   `json:"verified"`
	Millis   int64  `json:"millis"`
}
