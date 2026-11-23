// Package score filtra le offerte e assegna un punteggio rispetto al profilo.
package score

import (
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/rirolli/work_search/internal/config"
	"github.com/rirolli/work_search/internal/model"
)

// Norm porta il testo in minuscolo, sostituisce la punteggiatura con spazi e
// aggiunge uno spazio ai bordi, così i termini si confrontano come parole intere.
func Norm(s string) string {
	var b strings.Builder
	b.Grow(len(s) + 2)
	b.WriteByte(' ')
	space := true
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			space = false
		} else if !space {
			b.WriteByte(' ')
			space = true
		}
	}
	if !space {
		b.WriteByte(' ')
	}
	return b.String()
}

// Has dice se il testo normalizzato contiene il termine come parola/e intere.
func Has(normText, term string) bool {
	t := Norm(term)
	if strings.TrimSpace(t) == "" {
		return false
	}
	return strings.Contains(normText, t)
}

// Reason spiega perché un'offerta è stata scartata (vuota se tenuta).
type Reason string

// Evaluate calcola punteggio e priorità. Restituisce false se l'offerta va scartata.
func Evaluate(j *model.Job, p config.Profile, now time.Time) (bool, Reason) {
	title := Norm(j.Title)
	desc := Norm(j.Description + " " + j.Department)
	loc := Norm(j.Location)
	if j.Remote {
		loc += "remote "
	}

	for _, t := range p.TitleExclude {
		if Has(title, t) {
			return false, Reason("titolo escluso: " + t)
		}
	}
	if len(p.TitleRequiredAny) > 0 && !hasAny(title, p.TitleRequiredAny) {
		// Ruolo gestionale: ammesso solo se il contesto è ML/AI/EO.
		if !hasAny(title, p.TitleManagerAny) {
			return false, "titolo fuori profilo"
		}
		if len(p.ManagerRequiresAny) > 0 && !hasAny(title+desc, p.ManagerRequiresAny) {
			return false, "ruolo gestionale fuori ambito ML/EO"
		}
	}

	locScore, locMatched, strong := 0, "", false
	for _, l := range p.Locations {
		if !Has(loc, l.Term) {
			continue
		}
		if !l.Weak {
			strong = true
		}
		if locMatched == "" || l.Weight > locScore {
			locScore, locMatched = l.Weight, l.Term
		}
	}
	if !strong {
		for _, x := range p.LocationExclude {
			if Has(loc, x) {
				return false, Reason("sede esclusa: " + x)
			}
		}
	}

	if p.MaxAgeDays > 0 && !j.Published.IsZero() && now.Sub(j.Published) > time.Duration(p.MaxAgeDays)*24*time.Hour {
		return false, "offerta troppo vecchia"
	}

	score, descScore := 0, 0
	var matched []string
	for _, k := range p.Keywords {
		switch {
		case Has(title, k.Term):
			score += 2 * k.Weight
			matched = append(matched, k.Term)
		case Has(desc, k.Term):
			descScore += k.Weight
			if k.Weight > 0 {
				matched = append(matched, k.Term)
			}
		}
	}
	if p.DescriptionMax > 0 && descScore > p.DescriptionMax {
		descScore = p.DescriptionMax
	}
	score += descScore + locScore
	if locMatched != "" {
		matched = append(matched, "sede: "+locMatched)
	}
	score += j.CompanyRating * p.CompanyRatingFactor
	if p.RecencyDays > 0 && !j.Published.IsZero() && now.Sub(j.Published) <= time.Duration(p.RecencyDays)*24*time.Hour {
		score += p.RecencyBonus
	}
	if score < 0 {
		score = 0
	}
	j.Score = score
	j.Matched = matched

	switch {
	case score >= p.PriorityHigh:
		j.Priority = "Alta"
	case score >= p.PriorityMedium:
		j.Priority = "Media"
	default:
		j.Priority = "Bassa"
	}

	j.SalaryOK = "n.d."
	if j.SalaryYearly && strings.EqualFold(j.SalaryCurrency, "EUR") && (j.SalaryMax > 0 || j.SalaryMin > 0) {
		top := j.SalaryMax
		if top == 0 {
			top = j.SalaryMin
		}
		if top >= p.SalaryMinEUR {
			j.SalaryOK = "Sì"
		} else {
			j.SalaryOK = "No"
		}
	}

	if score < p.MinScore {
		return false, "punteggio sotto soglia"
	}
	return true, ""
}

func hasAny(normText string, terms []string) bool {
	for _, t := range terms {
		if Has(normText, t) {
			return true
		}
	}
	return false
}

// Sort ordina per priorità, poi punteggio, poi data.
func Sort(jobs []model.Job) {
	rank := map[string]int{"Alta": 0, "Media": 1, "Bassa": 2}
	sort.SliceStable(jobs, func(a, b int) bool {
		ja, jb := jobs[a], jobs[b]
		if rank[ja.Priority] != rank[jb.Priority] {
			return rank[ja.Priority] < rank[jb.Priority]
		}
		if ja.Score != jb.Score {
			return ja.Score > jb.Score
		}
		return ja.Published.After(jb.Published)
	})
}
