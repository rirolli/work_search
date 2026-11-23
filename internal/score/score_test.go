package score

import (
	"testing"
	"time"

	"github.com/rirolli/work_search/internal/config"
	"github.com/rirolli/work_search/internal/model"
)

func profile() config.Profile {
	return config.Profile{
		TitleRequiredAny: []string{"machine learning", "computer vision", "ai"},
		TitleExclude:     []string{"intern", "internship"},
		Keywords: []config.Weighted{
			{Term: "earth observation", Weight: 15},
			{Term: "machine learning", Weight: 8},
			{Term: "sar", Weight: 12},
			{Term: "junior", Weight: -15},
		},
		Locations: []config.Weighted{
			{Term: "roma", Weight: 30},
			{Term: "germany", Weight: 6},
			{Term: "remote", Weight: 10, Weak: true},
		},
		LocationExclude:     []string{"usa", "us"},
		MinScore:            10,
		PriorityHigh:        60,
		PriorityMedium:      40,
		SalaryMinEUR:        45000,
		CompanyRatingFactor: 2,
	}
}

func TestNormWholeWords(t *testing.T) {
	n := Norm("Senior ML/AI Engineer (SAR) – Roma")
	for _, term := range []string{"ml", "ai", "sar", "roma", "senior ml"} {
		if !Has(n, term) {
			t.Errorf("atteso %q in %q", term, n)
		}
	}
	for _, term := range []string{"sa", "rom", "aI engineer x"} {
		if Has(n, term) {
			t.Errorf("non atteso %q in %q", term, n)
		}
	}
	if Has(Norm("Maintenance planner"), "ai") {
		t.Error("'ai' non deve corrispondere dentro 'maintenance'")
	}
}

func TestEvaluate(t *testing.T) {
	now := time.Date(2026, 11, 23, 8, 0, 0, 0, time.UTC)
	p := profile()
	cases := []struct {
		name     string
		job      model.Job
		keep     bool
		priority string
	}{
		{"roma EO", model.Job{Title: "Machine Learning Engineer", Location: "Roma, Italia", Description: "SAR earth observation", CompanyRating: 8}, true, "Alta"},
		{"titolo fuori profilo", model.Job{Title: "Accountant", Location: "Roma"}, false, ""},
		{"stage", model.Job{Title: "Machine Learning Intern", Location: "Roma"}, false, ""},
		{"solo USA", model.Job{Title: "Computer Vision Engineer", Location: "Seattle, USA"}, false, ""},
		{"remoto USA", model.Job{Title: "Computer Vision Engineer", Location: "Remote - US", Remote: true}, false, ""},
		{"Germania + USA", model.Job{Title: "Computer Vision Engineer", Location: "Berlin, Germany / New York, USA", Description: "earth observation"}, true, "Bassa"},
	}
	for _, c := range cases {
		j := c.job
		keep, why := Evaluate(&j, p, now)
		if keep != c.keep {
			t.Errorf("%s: keep=%v (%s), atteso %v; score %d", c.name, keep, why, c.keep, j.Score)
			continue
		}
		if keep && j.Priority != c.priority {
			t.Errorf("%s: priorità %s (score %d), attesa %s", c.name, j.Priority, j.Score, c.priority)
		}
	}
}

func TestSalary(t *testing.T) {
	p := profile()
	j := model.Job{Title: "AI Engineer", Location: "Roma", SalaryMax: 50000, SalaryCurrency: "EUR", SalaryYearly: true}
	Evaluate(&j, p, time.Now())
	if j.SalaryOK != "Sì" {
		t.Errorf("SalaryOK=%s", j.SalaryOK)
	}
	j = model.Job{Title: "AI Engineer", Location: "Roma", SalaryMax: 60000, SalaryCurrency: "GBP", SalaryYearly: true}
	Evaluate(&j, p, time.Now())
	if j.SalaryOK != "n.d." {
		t.Errorf("valuta diversa da EUR: SalaryOK=%s", j.SalaryOK)
	}
}
