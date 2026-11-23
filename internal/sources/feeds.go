package sources

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/rirolli/work_search/internal/config"
	"github.com/rirolli/work_search/internal/fetch"
	"github.com/rirolli/work_search/internal/model"
)

func init() {
	ATS["teamtailor"] = teamtailor
	ATS["pinpoint"] = pinpoint
	ATS["successfactors"] = successFactors
	ATS["bamboohr"] = bambooHR
	ATS["breezy"] = breezy
}

// --- BambooHR (JSON /careers/list) --------------------------------------------

func bambooHR(ctx context.Context, c *fetch.Client, s config.Source, company string) ([]model.Job, error) {
	h := base(s, "https://"+s.ID+".bamboohr.com")
	var resp struct {
		Result []struct {
			ID              json64 `json:"id"`
			JobOpeningName  string `json:"jobOpeningName"`
			DepartmentLabel string `json:"departmentLabel"`
			IsRemote        *bool  `json:"isRemote"`
			LocationType    string `json:"locationType"`
			Location        struct {
				City  string `json:"city"`
				State string `json:"state"`
			} `json:"location"`
		} `json:"result"`
	}
	if err := c.GetJSON(ctx, h+"/careers/list", &resp); err != nil {
		return nil, err
	}
	out := make([]model.Job, 0, len(resp.Result))
	for _, r := range resp.Result {
		out = append(out, model.Job{
			Key:        key("bamboohr", s.ID, string(r.ID)),
			Company:    company,
			Title:      strings.TrimSpace(r.JobOpeningName),
			Location:   strings.Trim(r.Location.City+", "+r.Location.State, ", "),
			Remote:     (r.IsRemote != nil && *r.IsRemote) || r.LocationType == "1",
			Department: r.DepartmentLabel,
			URL:        h + "/careers/" + string(r.ID),
			Source:     "BambooHR",
		})
	}
	return out, nil
}

// --- Breezy HR (JSON /json) -----------------------------------------------------

func breezy(ctx context.Context, c *fetch.Client, s config.Source, company string) ([]model.Job, error) {
	h := base(s, "https://"+s.ID+".breezy.hr")
	var resp []struct {
		ID            string `json:"id"`
		Name          string `json:"name"`
		URL           string `json:"url"`
		PublishedDate string `json:"published_date"`
		Department    string `json:"department"`
		Location      struct {
			Name     string `json:"name"`
			IsRemote bool   `json:"is_remote"`
		} `json:"location"`
	}
	if err := c.GetJSON(ctx, h+"/json", &resp); err != nil {
		return nil, err
	}
	out := make([]model.Job, 0, len(resp))
	for _, r := range resp {
		out = append(out, model.Job{
			Key:        key("breezy", s.ID, r.ID),
			Company:    company,
			Title:      strings.TrimSpace(r.Name),
			Location:   r.Location.Name,
			Remote:     r.Location.IsRemote,
			Department: r.Department,
			URL:        r.URL,
			Source:     "Breezy",
			Published:  fetch.ParseTime(r.PublishedDate),
		})
	}
	return out, nil
}

// host restituisce l'URL base di una fonte che vive su un dominio proprio.
func host(s config.Source) (string, error) {
	switch {
	case s.BaseURL != "":
		return strings.TrimRight(s.BaseURL, "/"), nil
	case s.Domain != "":
		return "https://" + strings.TrimRight(s.Domain, "/"), nil
	}
	return "", fmt.Errorf("%s: serve domain (es. careers.azienda.com)", s.Type)
}

// --- Teamtailor (feed RSS /jobs.rss) -----------------------------------------

type rssFeed struct {
	Items []struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		GUID        string `xml:"guid"`
		PubDate     string `xml:"pubDate"`
		Description string `xml:"description"`
		Remote      string `xml:"remoteStatus"`
		Department  string `xml:"department"`
		Locations   []struct {
			Name    string `xml:"name"`
			City    string `xml:"city"`
			Country string `xml:"country"`
		} `xml:"locations>location"`
	} `xml:"channel>item"`
}

func parseRSS(data []byte) (*rssFeed, error) {
	var f rssFeed
	if err := xml.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("RSS non valido: %w", err)
	}
	return &f, nil
}

func teamtailor(ctx context.Context, c *fetch.Client, s config.Source, company string) ([]model.Job, error) {
	h, err := host(s)
	if err != nil {
		return nil, err
	}
	data, err := c.Do(ctx, http.MethodGet, h+"/jobs.rss", nil, nil)
	if err != nil {
		return nil, err
	}
	f, err := parseRSS(data)
	if err != nil {
		return nil, err
	}
	out := make([]model.Job, 0, len(f.Items))
	for _, it := range f.Items {
		var locs []string
		for _, l := range it.Locations {
			parts := []string{}
			for _, p := range []string{l.City, l.Country} {
				if p != "" {
					parts = append(parts, p)
				}
			}
			if len(parts) == 0 && l.Name != "" {
				parts = append(parts, l.Name)
			}
			if len(parts) > 0 {
				locs = append(locs, strings.Join(parts, ", "))
			}
		}
		id := it.GUID
		if id == "" {
			id = it.Link
		}
		r := strings.ToLower(it.Remote)
		out = append(out, model.Job{
			Key:         key("teamtailor", h, id),
			Company:     company,
			Title:       strings.TrimSpace(it.Title),
			Location:    strings.Join(locs, " / "),
			Remote:      r == "fully" || r == "remote",
			Department:  it.Department,
			URL:         it.Link,
			Source:      "Teamtailor",
			Published:   parseRFC1123(it.PubDate),
			Description: fetch.PlainText(it.Description),
		})
	}
	return out, nil
}

// --- Pinpoint (JSON /postings.json) ------------------------------------------

func pinpoint(ctx context.Context, c *fetch.Client, s config.Source, company string) ([]model.Job, error) {
	h, err := host(s)
	if err != nil {
		if s.ID == "" {
			return nil, err
		}
		h = "https://" + s.ID + ".pinpointhq.com"
	}
	var resp struct {
		Data []struct {
			ID            json64 `json:"id"`
			Title         string `json:"title"`
			URL           string `json:"url"`
			WorkplaceType string `json:"workplace_type"`
			Description   string `json:"description"`
			Requirements  string `json:"key_responsibilities"`
			Location      struct {
				Name     string `json:"name"`
				City     string `json:"city"`
				Province string `json:"province"`
			} `json:"location"`
			Department struct {
				Name string `json:"name"`
			} `json:"department"`
			Compensation        string `json:"compensation"`
			CompensationVisible bool   `json:"compensation_visible"`
		} `json:"data"`
	}
	b, err := c.Do(ctx, http.MethodGet, h+"/postings.json", nil, nil)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, fmt.Errorf("JSON non valido: %w", err)
	}
	// Per un account inesistente Pinpoint risponde con una board dimostrativa:
	// teniamo solo gli annunci il cui link è sul dominio richiesto.
	reqHost := strings.TrimPrefix(strings.TrimPrefix(h, "https://"), "http://")
	out := make([]model.Job, 0, len(resp.Data))
	for _, p := range resp.Data {
		if u, err := url.Parse(p.URL); err == nil && u.Host != "" && !strings.EqualFold(u.Host, reqHost) {
			continue
		}
		loc := p.Location.Name
		if loc == "" {
			loc = strings.Trim(p.Location.City+", "+p.Location.Province, ", ")
		}
		j := model.Job{
			Key:         key("pinpoint", h, string(p.ID)),
			Company:     company,
			Title:       strings.TrimSpace(p.Title),
			Location:    loc,
			Remote:      strings.Contains(strings.ToLower(p.WorkplaceType), "remote"),
			Department:  p.Department.Name,
			URL:         p.URL,
			Source:      "Pinpoint",
			Description: fetch.PlainText(p.Description + " " + p.Requirements),
		}
		if p.CompensationVisible {
			j.Salary = p.Compensation
		}
		out = append(out, j)
	}
	return out, nil
}

// --- SAP SuccessFactors Recruiting Marketing (feed RSS /services/rss/job/) -----

var sfLocRe = regexp.MustCompile(`\(([^()]*)\)\s*$`)

func successFactors(ctx context.Context, c *fetch.Client, s config.Source, company string) ([]model.Job, error) {
	h, err := host(s)
	if err != nil {
		return nil, err
	}
	queries := s.Queries
	if len(queries) == 0 {
		queries = []string{""}
	}
	seen := map[string]bool{}
	var out []model.Job
	for _, q := range queries {
		u := h + "/services/rss/job/?locale=en_GB&keywords=" + url.QueryEscape(q)
		data, err := c.Do(ctx, http.MethodGet, u, nil, nil)
		if err != nil {
			if len(out) > 0 {
				return out, nil
			}
			return nil, err
		}
		f, err := parseRSS(data)
		if err != nil {
			return nil, err
		}
		for _, it := range f.Items {
			if seen[it.Link] {
				continue
			}
			seen[it.Link] = true
			title, loc := strings.TrimSpace(it.Title), ""
			if m := sfLocRe.FindStringSubmatch(title); m != nil { // "Titolo (Città, Paese)"
				loc = m[1]
				title = strings.TrimSpace(strings.TrimSuffix(title, m[0]))
			}
			out = append(out, model.Job{
				Key:         key("successfactors", h, it.Link),
				Company:     company,
				Title:       title,
				Location:    loc,
				URL:         it.Link,
				Source:      "SuccessFactors",
				Published:   parseRFC1123(it.PubDate),
				Description: fetch.PlainText(it.Description),
			})
		}
	}
	return out, nil
}

func parseRFC1123(s string) time.Time {
	s = strings.TrimSpace(s)
	for _, l := range []string{time.RFC1123Z, time.RFC1123, "Mon, 2 Jan 2006 15:04:05 -0700", "Mon, 2 Jan 2006 15:04:05 MST"} {
		if v, err := time.Parse(l, s); err == nil {
			return v.UTC()
		}
	}
	return fetch.ParseTime(s)
}
