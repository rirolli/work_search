// Package sources contiene i connettori verso le pagine carriere (ATS) e gli aggregatori.
package sources

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rirolli/work_search/internal/config"
	"github.com/rirolli/work_search/internal/fetch"
	"github.com/rirolli/work_search/internal/model"
)

// ATSFetcher scarica le offerte di una singola azienda da una fonte.
type ATSFetcher func(ctx context.Context, c *fetch.Client, s config.Source, company string) ([]model.Job, error)

// ATS è il registro dei connettori per tipo di fonte.
var ATS = map[string]ATSFetcher{
	"greenhouse":      greenhouse,
	"lever":           lever,
	"ashby":           ashby,
	"personio":        personio,
	"recruitee":       recruitee,
	"workable":        workable,
	"smartrecruiters": smartRecruiters,
	"workday":         workday,
}

func base(s config.Source, def string) string {
	if s.BaseURL != "" {
		return strings.TrimRight(s.BaseURL, "/")
	}
	return def
}

func key(src, company, id string) string {
	return src + ":" + strings.ToLower(company) + ":" + id
}

// --- Greenhouse -------------------------------------------------------------

func greenhouse(ctx context.Context, c *fetch.Client, s config.Source, company string) ([]model.Job, error) {
	var resp struct {
		Jobs []struct {
			ID             int64  `json:"id"`
			Title          string `json:"title"`
			AbsoluteURL    string `json:"absolute_url"`
			UpdatedAt      string `json:"updated_at"`
			FirstPublished string `json:"first_published"`
			Content        string `json:"content"`
			Location       struct {
				Name string `json:"name"`
			} `json:"location"`
			Departments []struct {
				Name string `json:"name"`
			} `json:"departments"`
		} `json:"jobs"`
	}
	u := base(s, "https://boards-api.greenhouse.io") + "/v1/boards/" + url.PathEscape(s.ID) + "/jobs?content=true"
	if err := c.GetJSON(ctx, u, &resp); err != nil {
		return nil, err
	}
	out := make([]model.Job, 0, len(resp.Jobs))
	for _, j := range resp.Jobs {
		job := model.Job{
			Key:         key("greenhouse", s.ID, strconv.FormatInt(j.ID, 10)),
			Company:     company,
			Title:       strings.TrimSpace(j.Title),
			Location:    j.Location.Name,
			URL:         j.AbsoluteURL,
			Source:      "Greenhouse",
			Description: fetch.PlainText(j.Content),
		}
		if t := fetch.ParseTime(j.FirstPublished); !t.IsZero() {
			job.Published = t
		} else if t := fetch.ParseTime(j.UpdatedAt); !t.IsZero() {
			job.Published = t
			job.PublishedNote = "data di aggiornamento"
		}
		if len(j.Departments) > 0 {
			job.Department = j.Departments[0].Name
		}
		job.Remote = strings.Contains(strings.ToLower(job.Location), "remote")
		out = append(out, job)
	}
	return out, nil
}

// --- Lever ------------------------------------------------------------------

func lever(ctx context.Context, c *fetch.Client, s config.Source, company string) ([]model.Job, error) {
	host := "https://api.lever.co"
	if strings.EqualFold(s.Region, "eu") {
		host = "https://api.eu.lever.co"
	}
	var resp []struct {
		ID               string `json:"id"`
		Text             string `json:"text"`
		HostedURL        string `json:"hostedUrl"`
		CreatedAt        int64  `json:"createdAt"`
		DescriptionPlain string `json:"descriptionPlain"`
		AdditionalPlain  string `json:"additionalPlain"`
		WorkplaceType    string `json:"workplaceType"`
		Categories       struct {
			Location     string   `json:"location"`
			Team         string   `json:"team"`
			AllLocations []string `json:"allLocations"`
		} `json:"categories"`
		Lists []struct {
			Content string `json:"content"`
		} `json:"lists"`
		SalaryRange *struct {
			Currency string  `json:"currency"`
			Interval string  `json:"interval"`
			Min      float64 `json:"min"`
			Max      float64 `json:"max"`
		} `json:"salaryRange"`
	}
	u := base(s, host) + "/v0/postings/" + url.PathEscape(s.ID) + "?mode=json"
	if err := c.GetJSON(ctx, u, &resp); err != nil {
		return nil, err
	}
	out := make([]model.Job, 0, len(resp))
	for _, j := range resp {
		desc := j.DescriptionPlain + " " + j.AdditionalPlain
		for _, l := range j.Lists {
			desc += " " + fetch.PlainText(l.Content)
		}
		loc := j.Categories.Location
		if len(j.Categories.AllLocations) > 1 {
			loc = strings.Join(j.Categories.AllLocations, " / ")
		}
		job := model.Job{
			Key:         key("lever", s.ID, j.ID),
			Company:     company,
			Title:       strings.TrimSpace(j.Text),
			Location:    loc,
			Remote:      strings.EqualFold(j.WorkplaceType, "remote"),
			Department:  j.Categories.Team,
			URL:         j.HostedURL,
			Source:      "Lever",
			Published:   fetch.FromMillis(j.CreatedAt),
			Description: desc,
		}
		if r := j.SalaryRange; r != nil {
			yearly := strings.Contains(r.Interval, "year")
			job.SalaryMin, job.SalaryMax, job.SalaryCurrency, job.SalaryYearly = r.Min, r.Max, r.Currency, yearly
			period := ""
			if yearly {
				period = "/anno"
			}
			job.Salary = fetch.FormatSalary(r.Min, r.Max, r.Currency, period)
		}
		out = append(out, job)
	}
	return out, nil
}

// --- Ashby ------------------------------------------------------------------

func ashby(ctx context.Context, c *fetch.Client, s config.Source, company string) ([]model.Job, error) {
	var resp struct {
		Jobs []struct {
			ID                 string `json:"id"`
			Title              string `json:"title"`
			Department         string `json:"department"`
			Location           string `json:"location"`
			SecondaryLocations []struct {
				Location string `json:"location"`
			} `json:"secondaryLocations"`
			PublishedAt      string `json:"publishedAt"`
			IsListed         *bool  `json:"isListed"`
			IsRemote         bool   `json:"isRemote"`
			WorkplaceType    string `json:"workplaceType"`
			JobURL           string `json:"jobUrl"`
			DescriptionPlain string `json:"descriptionPlain"`
			Compensation     *struct {
				Summary string `json:"compensationTierSummary"`
			} `json:"compensation"`
		} `json:"jobs"`
	}
	u := base(s, "https://api.ashbyhq.com") + "/posting-api/job-board/" + url.PathEscape(s.ID) + "?includeCompensation=true"
	if err := c.GetJSON(ctx, u, &resp); err != nil {
		return nil, err
	}
	out := make([]model.Job, 0, len(resp.Jobs))
	for _, j := range resp.Jobs {
		if j.IsListed != nil && !*j.IsListed {
			continue
		}
		locs := []string{j.Location}
		for _, l := range j.SecondaryLocations {
			if l.Location != "" {
				locs = append(locs, l.Location)
			}
		}
		job := model.Job{
			Key:         key("ashby", s.ID, j.ID),
			Company:     company,
			Title:       strings.TrimSpace(j.Title),
			Location:    strings.Join(locs, " / "),
			Remote:      j.IsRemote || strings.EqualFold(j.WorkplaceType, "remote"),
			Department:  j.Department,
			URL:         j.JobURL,
			Source:      "Ashby",
			Published:   fetch.ParseTime(j.PublishedAt),
			Description: j.DescriptionPlain,
		}
		if j.Compensation != nil {
			job.Salary = j.Compensation.Summary
		}
		out = append(out, job)
	}
	return out, nil
}

// --- Personio (feed XML pubblico) --------------------------------------------

type personioFeed struct {
	Positions []struct {
		ID         string `xml:"id"`
		Office     string `xml:"office"`
		Additional struct {
			Offices []string `xml:"office"`
		} `xml:"additionalOffices"`
		Department string `xml:"department"`
		Name       string `xml:"name"`
		CreatedAt  string `xml:"createdAt"`
		Schedule   string `xml:"schedule"`
		Descs      []struct {
			Value string `xml:"value"`
		} `xml:"jobDescriptions>jobDescription"`
	} `xml:"position"`
}

func personio(ctx context.Context, c *fetch.Client, s config.Source, company string) ([]model.Job, error) {
	tld := s.Domain
	if tld == "" {
		tld = "de"
	}
	host := fmt.Sprintf("https://%s.jobs.personio.%s", s.ID, tld)
	b := base(s, host)
	data, err := c.Do(ctx, http.MethodGet, b+"/xml?language=en", nil, nil)
	if err != nil {
		return nil, err
	}
	var feed personioFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, fmt.Errorf("XML non valido: %w", err)
	}
	out := make([]model.Job, 0, len(feed.Positions))
	for _, p := range feed.Positions {
		locs := append([]string{p.Office}, p.Additional.Offices...)
		var desc strings.Builder
		for _, d := range p.Descs {
			desc.WriteString(fetch.PlainText(d.Value))
			desc.WriteString(" ")
		}
		loc := strings.Trim(strings.Join(locs, " / "), " /")
		out = append(out, model.Job{
			Key:         key("personio", s.ID, p.ID),
			Company:     company,
			Title:       strings.TrimSpace(p.Name),
			Location:    loc,
			Remote:      strings.Contains(strings.ToLower(loc), "remote"),
			Department:  p.Department,
			URL:         b + "/job/" + p.ID,
			Source:      "Personio",
			Published:   fetch.ParseTime(p.CreatedAt),
			Description: desc.String(),
		})
	}
	return out, nil
}

// --- Recruitee --------------------------------------------------------------

func recruitee(ctx context.Context, c *fetch.Client, s config.Source, company string) ([]model.Job, error) {
	var resp struct {
		Offers []struct {
			ID           int64  `json:"id"`
			Title        string `json:"title"`
			CareersURL   string `json:"careers_url"`
			Location     string `json:"location"`
			Remote       bool   `json:"remote"`
			Department   string `json:"department"`
			PublishedAt  string `json:"published_at"`
			CreatedAt    string `json:"created_at"`
			Description  string `json:"description"`
			Requirements string `json:"requirements"`
		} `json:"offers"`
	}
	u := base(s, "https://"+s.ID+".recruitee.com") + "/api/offers/"
	if err := c.GetJSON(ctx, u, &resp); err != nil {
		return nil, err
	}
	out := make([]model.Job, 0, len(resp.Offers))
	for _, o := range resp.Offers {
		pub := fetch.ParseTime(o.PublishedAt)
		if pub.IsZero() {
			pub = fetch.ParseTime(o.CreatedAt)
		}
		out = append(out, model.Job{
			Key:         key("recruitee", s.ID, strconv.FormatInt(o.ID, 10)),
			Company:     company,
			Title:       strings.TrimSpace(o.Title),
			Location:    o.Location,
			Remote:      o.Remote,
			Department:  o.Department,
			URL:         o.CareersURL,
			Source:      "Recruitee",
			Published:   pub,
			Description: fetch.PlainText(o.Description + " " + o.Requirements),
		})
	}
	return out, nil
}

// --- Workable ---------------------------------------------------------------

func workable(ctx context.Context, c *fetch.Client, s config.Source, company string) ([]model.Job, error) {
	var resp struct {
		Jobs []struct {
			Title         string `json:"title"`
			Shortcode     string `json:"shortcode"`
			URL           string `json:"url"`
			Telecommuting bool   `json:"telecommuting"`
			Department    string `json:"department"`
			PublishedOn   string `json:"published_on"`
			CreatedAt     string `json:"created_at"`
			City          string `json:"city"`
			Country       string `json:"country"`
			Description   string `json:"description"`
		} `json:"jobs"`
	}
	u := base(s, "https://apply.workable.com") + "/api/v1/widget/accounts/" + url.PathEscape(s.ID) + "?details=true"
	if err := c.GetJSON(ctx, u, &resp); err != nil {
		return nil, err
	}
	out := make([]model.Job, 0, len(resp.Jobs))
	for _, j := range resp.Jobs {
		pub := fetch.ParseTime(j.PublishedOn)
		if pub.IsZero() {
			pub = fetch.ParseTime(j.CreatedAt)
		}
		out = append(out, model.Job{
			Key:         key("workable", s.ID, j.Shortcode),
			Company:     company,
			Title:       strings.TrimSpace(j.Title),
			Location:    strings.Trim(j.City+", "+j.Country, ", "),
			Remote:      j.Telecommuting,
			Department:  j.Department,
			URL:         j.URL,
			Source:      "Workable",
			Published:   pub,
			Description: fetch.PlainText(j.Description),
		})
	}
	return out, nil
}

// --- SmartRecruiters --------------------------------------------------------

func smartRecruiters(ctx context.Context, c *fetch.Client, s config.Source, company string) ([]model.Job, error) {
	pages := s.Pages
	if pages <= 0 {
		pages = 5
	}
	var out []model.Job
	for page := 0; page < pages; page++ {
		var resp struct {
			TotalFound int `json:"totalFound"`
			Content    []struct {
				ID           string `json:"id"`
				Name         string `json:"name"`
				ReleasedDate string `json:"releasedDate"`
				Location     struct {
					FullLocation string `json:"fullLocation"`
					City         string `json:"city"`
					Country      string `json:"country"`
					Remote       bool   `json:"remote"`
				} `json:"location"`
				Department struct {
					Label string `json:"label"`
				} `json:"department"`
			} `json:"content"`
		}
		u := fmt.Sprintf("%s/v1/companies/%s/postings?limit=100&offset=%d",
			base(s, "https://api.smartrecruiters.com"), url.PathEscape(s.ID), page*100)
		if err := c.GetJSON(ctx, u, &resp); err != nil {
			return out, err
		}
		for _, p := range resp.Content {
			loc := p.Location.FullLocation
			if loc == "" {
				loc = strings.Trim(p.Location.City+", "+p.Location.Country, ", ")
			}
			out = append(out, model.Job{
				Key:        key("smartrecruiters", s.ID, p.ID),
				Company:    company,
				Title:      strings.TrimSpace(p.Name),
				Location:   loc,
				Remote:     p.Location.Remote,
				Department: p.Department.Label,
				URL:        "https://jobs.smartrecruiters.com/" + s.ID + "/" + p.ID,
				Source:     "SmartRecruiters",
				Published:  fetch.ParseTime(p.ReleasedDate),
			})
		}
		if (page+1)*100 >= resp.TotalFound || len(resp.Content) == 0 {
			break
		}
	}
	return out, nil
}

// --- Workday (endpoint JSON usato dalle pagine carriere *.myworkdayjobs.com) --

var workdayDaysRe = regexp.MustCompile(`(?i)(\d+)\+?\s*days?`)

// WorkdayPosted interpreta "Posted Today", "Posted Yesterday", "Posted 3 Days Ago", "Posted 30+ Days Ago".
func WorkdayPosted(s string, now time.Time) time.Time {
	l := strings.ToLower(s)
	switch {
	case strings.Contains(l, "today"):
		return now.Truncate(24 * time.Hour)
	case strings.Contains(l, "yesterday"):
		return now.Add(-24 * time.Hour).Truncate(24 * time.Hour)
	case strings.Contains(l, "+"):
		return time.Time{} // "30+ days": data esatta ignota
	}
	if m := workdayDaysRe.FindStringSubmatch(l); m != nil {
		n, _ := strconv.Atoi(m[1])
		return now.Add(-time.Duration(n) * 24 * time.Hour).Truncate(24 * time.Hour)
	}
	return time.Time{}
}

func workday(ctx context.Context, c *fetch.Client, s config.Source, company string) ([]model.Job, error) {
	if s.Domain == "" || s.Tenant == "" || s.Site == "" {
		return nil, fmt.Errorf("workday: servono domain, tenant e site")
	}
	host := base(s, "https://"+s.Domain)
	endpoint := fmt.Sprintf("%s/wday/cxs/%s/%s/jobs", host, s.Tenant, s.Site)
	queries := s.Queries
	if len(queries) == 0 {
		queries = []string{""}
	}
	pages := s.Pages
	if pages <= 0 {
		pages = 3
	}
	now := time.Now().UTC()
	seen := map[string]bool{}
	var out []model.Job
	for _, q := range queries {
		for page := 0; page < pages; page++ {
			body := map[string]any{"appliedFacets": map[string]any{}, "limit": 20, "offset": page * 20, "searchText": q}
			var resp struct {
				Total       int `json:"total"`
				JobPostings []struct {
					Title         string   `json:"title"`
					ExternalPath  string   `json:"externalPath"`
					LocationsText string   `json:"locationsText"`
					PostedOn      string   `json:"postedOn"`
					BulletFields  []string `json:"bulletFields"`
				} `json:"jobPostings"`
			}
			if err := c.PostJSON(ctx, endpoint, body, &resp); err != nil {
				if len(out) > 0 {
					return out, nil // risultati parziali
				}
				return nil, err
			}
			for _, p := range resp.JobPostings {
				if p.ExternalPath == "" || seen[p.ExternalPath] {
					continue
				}
				seen[p.ExternalPath] = true
				out = append(out, model.Job{
					Key:           key("workday", s.Tenant, p.ExternalPath),
					Company:       company,
					Title:         strings.TrimSpace(p.Title),
					Location:      p.LocationsText,
					URL:           host + "/" + s.Site + p.ExternalPath,
					Source:        "Workday",
					Published:     WorkdayPosted(p.PostedOn, now),
					PublishedNote: p.PostedOn,
				})
			}
			if (page+1)*20 >= resp.Total || len(resp.JobPostings) == 0 {
				break
			}
		}
	}
	return out, nil
}
