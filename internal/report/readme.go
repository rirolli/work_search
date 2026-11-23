package report

import (
	"fmt"
	"os"
	"strings"
)

// Marcatori della sezione del README riscritta a ogni esecuzione.
const (
	ReadmeStart = "<!-- RISULTATI:INIZIO -->"
	ReadmeEnd   = "<!-- RISULTATI:FINE -->"
)

func mdCell(s string) string {
	s = strings.NewReplacer("|", "/", "\n", " ", "\r", " ", "[", "(", "]", ")").Replace(s)
	if r := []rune(s); len(r) > 70 {
		s = string(r[:69]) + "…"
	}
	return strings.TrimSpace(s)
}

// ReadmeSection produce il markdown dei risultati (senza marcatori).
func ReadmeSection(d Data) string {
	var b strings.Builder
	fmt.Fprintf(&b, "_Aggiornato il %s (ora di Roma) · %d offerte in linea con il profilo su %d scaricate · "+
		"**%d alta**, %d media, %d bassa · **%d nuove**_\n\n",
		d.Generated.In(rome).Format("02/01/2006 15:04"), len(d.Jobs), d.Fetched,
		d.count("Alta"), d.count("Media"), d.count("Bassa"), d.count("new"))
	b.WriteString("Filtri e ordinamento: [`docs/index.html`](docs/index.html) · Excel: [`docs/offerte.xlsx`](docs/offerte.xlsx)\n\n")

	if len(d.Jobs) == 0 {
		b.WriteString("Nessuna offerta in linea con il profilo in questa esecuzione.\n")
	} else {
		fmt.Fprintf(&b, "| Priorità | Score | | Azienda | Ruolo | Sede | Pubblicata |\n|---|---:|---|---|---|---|---|\n")
		for _, j := range d.Jobs {
			nuova := ""
			if j.New {
				nuova = "🆕"
			}
			fmt.Fprintf(&b, "| %s | %d | %s | %s | [%s](%s) | %s | %s |\n",
				j.Priority, j.Score, nuova, mdCell(j.Company), mdCell(j.Title), j.URL, mdCell(j.Location), mdCell(published(j)))
		}
	}

	var bad []string
	for _, s := range d.Sources {
		if !s.OK && !s.Skipped {
			bad = append(bad, s.Name)
		}
	}
	if len(bad) > 0 {
		fmt.Fprintf(&b, "\n⚠️ Fonti in errore (%d): %s. Dettagli nel foglio *Fonti*.\n", len(bad), strings.Join(bad, ", "))
	}
	return b.String()
}

// UpdateReadme sostituisce la sezione tra i marcatori; se mancano, la aggiunge in fondo.
func UpdateReadme(path string, d Data) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	s := string(b)
	block := ReadmeStart + "\n" + ReadmeSection(d) + ReadmeEnd
	i, j := strings.Index(s, ReadmeStart), strings.Index(s, ReadmeEnd)
	if i >= 0 && j > i {
		s = s[:i] + block + s[j+len(ReadmeEnd):]
	} else {
		s = strings.TrimRight(s, "\n") + "\n\n## Risultati\n\n" + block + "\n"
	}
	return os.WriteFile(path, []byte(s), 0o644)
}
