package sources

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rirolli/work_search/internal/config"
	"github.com/rirolli/work_search/internal/fetch"
	"github.com/rirolli/work_search/internal/model"
)

// ErrMissingKey indica che la fonte richiede una chiave API non configurata.
var ErrMissingKey = errors.New("chiave API non configurata")

// AggFetcher esegue una ricerca su un aggregatore.
type AggFetcher func(ctx context.Context, c *fetch.Client, s config.Search) ([]model.Job, error)

// Aggregators è il registro degli aggregatori.
var Aggregators = map[string]AggFetcher{
	"adzuna":    adzuna,
	"jooble":    jooble,
	"arbeitnow": arbeitnow,
	"remotive":  remotive,
}

// Endpoint sovrascrivibili nei test.
var (
	AdzunaBase    = "https://api.adzuna.com"
	JoobleBase    = "https://jooble.org"
	ArbeitnowBase = "https://www.arbeitnow.com"
	RemotiveBase  = "https://remotive.com"
)

var adzunaCurrency = map[string]string{
	"it": "EUR", "de": "EUR", "fr": "EUR", "nl": "EUR", "es": "EUR", "at": "EUR", "be": "EUR",
	"gb": "GBP", "ch": "CHF", "pl": "PLN", "us": "USD",
}

func adzuna(ctx context.Context, c *fetch.Client, s config.Search) ([]model.Job, error) {
	id, k := os.Getenv("ADZUNA_APP_ID"), os.Getenv("ADZUNA_APP_KEY")
	if id == "" || k == "" {
		return nil, ErrMissingKey
	}
	country := s.Country
	if country == "" {
		country = "it"
	}
	pages := s.Pages
	if pages <= 0 {
		pages = 1
	}
	var out []model.Job
	for page := 1; page <= pages; page++ {
		q := url.Values{}
		q.Set("app_id", id)
		q.Set("app_key", k)
		q.Set("results_per_page", "50")
		if s.Query != "" {
			q.Set("what", s.Query)
		}
		if s.Where != "" {
			q.Set("where", s.Where)
		}
		if s.MaxDaysOld > 0 {
			q.Set("max_days_old", strconv.Itoa(s.MaxDaysOld))
		}
		var resp struct {
			Results []struct {
				ID          string  `json:"id"`
				Title       string  `json:"title"`
				Description string  `json:"description"`
				Created     string  `json:"created"`
				RedirectURL string  `json:"redirect_url"`
				SalaryMin   float64 `json:"salary_min"`
				SalaryMax   float64 `json:"salary_max"`
				Predicted   string  `json:"salary_is_predicted"`
				Company     struct {
					DisplayName string `json:"display_name"`
				} `json:"company"`
				Location struct {
					DisplayName string `json:"display_name"`
				} `json:"location"`
			} `json:"results"`
		}
		u := fmt.Sprintf("%s/v1/api/jobs/%s/search/%d?%s", AdzunaBase, country, page, q.Encode())
		if err := c.GetJSON(ctx, u, &resp); err != nil {
			return out, err
		}
		cur := adzunaCurrency[country]
		for _, r := range resp.Results {
			j := model.Job{
				Key:         "adzuna:" + r.ID,
				Company:     strings.TrimSpace(r.Company.DisplayName),
				Title:       fetch.PlainText(r.Title),
				Location:    r.Location.DisplayName,
				URL:         r.RedirectURL,
				Source:      "Adzuna",
				Published:   fetch.ParseTime(r.Created),
				Description: fetch.PlainText(r.Description),
			}
			if r.SalaryMin > 0 || r.SalaryMax > 0 {
				note := "/anno"
				if r.Predicted == "1" {
					note += " (stima Adzuna)"
				} else {
					j.SalaryMin, j.SalaryMax, j.SalaryCurrency, j.SalaryYearly = r.SalaryMin, r.SalaryMax, cur, true
				}
				j.Salary = fetch.FormatSalary(r.SalaryMin, r.SalaryMax, cur, note)
			}
			out = append(out, j)
		}
		if len(resp.Results) < 50 {
			break
		}
	}
	return out, nil
}

func jooble(ctx context.Context, c *fetch.Client, s config.Search) ([]model.Job, error) {
	k := os.Getenv("JOOBLE_API_KEY")
	if k == "" {
		return nil, ErrMissingKey
	}
	pages := s.Pages
	if pages <= 0 {
		pages = 1
	}
	var out []model.Job
	for page := 1; page <= pages; page++ {
		body := map[string]string{"keywords": s.Query, "location": s.Where, "page": strconv.Itoa(page)}
		var resp struct {
			Jobs []struct {
				ID       json64 `json:"id"`
				Title    string `json:"title"`
				Location string `json:"location"`
				Snippet  string `json:"snippet"`
				Salary   string `json:"salary"`
				Link     string `json:"link"`
				Company  string `json:"company"`
				Updated  string `json:"updated"`
				Source   string `json:"source"`
			} `json:"jobs"`
		}
		if err := c.PostJSON(ctx, JoobleBase+"/api/"+url.PathEscape(k), body, &resp); err != nil {
			return out, err
		}
		for _, r := range resp.Jobs {
			out = append(out, model.Job{
				Key:           "jooble:" + string(r.ID),
				Company:       strings.TrimSpace(r.Company),
				Title:         fetch.PlainText(r.Title),
				Location:      r.Location,
				URL:           r.Link,
				Source:        "Jooble" + suffix(r.Source),
				Published:     fetch.ParseTime(r.Updated),
				PublishedNote: "data di aggiornamento (Jooble)",
				Salary:        r.Salary,
				Description:   fetch.PlainText(r.Snippet),
			})
		}
		if len(resp.Jobs) == 0 {
			break
		}
	}
	return out, nil
}

func suffix(s string) string {
	if s == "" {
		return ""
	}
	return " (" + s + ")"
}

// json64 accetta id numerici o stringa.
type json64 string

func (j *json64) UnmarshalJSON(b []byte) error {
	*j = json64(strings.Trim(string(b), `"`))
	return nil
}

func arbeitnow(ctx context.Context, c *fetch.Client, s config.Search) ([]model.Job, error) {
	pages := s.Pages
	if pages <= 0 {
		pages = 3
	}
	var out []model.Job
	for page := 1; page <= pages; page++ {
		var resp struct {
			Data []struct {
				Slug        string   `json:"slug"`
				CompanyName string   `json:"company_name"`
				Title       string   `json:"title"`
				Description string   `json:"description"`
				Remote      bool     `json:"remote"`
				URL         string   `json:"url"`
				Tags        []string `json:"tags"`
				Location    string   `json:"location"`
				CreatedAt   int64    `json:"created_at"`
			} `json:"data"`
		}
		if err := c.GetJSON(ctx, fmt.Sprintf("%s/api/job-board-api?page=%d", ArbeitnowBase, page), &resp); err != nil {
			return out, err
		}
		for _, r := range resp.Data {
			out = append(out, model.Job{
				Key:         "arbeitnow:" + r.Slug,
				Company:     r.CompanyName,
				Title:       r.Title,
				Location:    r.Location,
				Remote:      r.Remote,
				URL:         r.URL,
				Source:      "Arbeitnow",
				Published:   time.Unix(r.CreatedAt, 0).UTC(),
				Description: fetch.PlainText(r.Description) + " " + strings.Join(r.Tags, " "),
			})
		}
		if len(resp.Data) == 0 {
			break
		}
	}
	return out, nil
}

func remotive(ctx context.Context, c *fetch.Client, s config.Search) ([]model.Job, error) {
	q := url.Values{}
	if s.Query != "" {
		q.Set("search", s.Query)
	}
	q.Set("limit", "100")
	var resp struct {
		Jobs []struct {
			ID                int64  `json:"id"`
			URL               string `json:"url"`
			Title             string `json:"title"`
			CompanyName       string `json:"company_name"`
			PublicationDate   string `json:"publication_date"`
			CandidateLocation string `json:"candidate_required_location"`
			Salary            string `json:"salary"`
			Description       string `json:"description"`
		} `json:"jobs"`
	}
	if err := c.GetJSON(ctx, RemotiveBase+"/api/remote-jobs?"+q.Encode(), &resp); err != nil {
		return nil, err
	}
	out := make([]model.Job, 0, len(resp.Jobs))
	for _, r := range resp.Jobs {
		out = append(out, model.Job{
			Key:         "remotive:" + strconv.FormatInt(r.ID, 10),
			Company:     r.CompanyName,
			Title:       r.Title,
			Location:    "Remote – " + r.CandidateLocation,
			Remote:      true,
			URL:         r.URL,
			Source:      "Remotive",
			Published:   fetch.ParseTime(r.PublicationDate),
			Salary:      r.Salary,
			Description: fetch.PlainText(r.Description),
		})
	}
	return out, nil
}
