package sources

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rirolli/work_search/internal/config"
	"github.com/rirolli/work_search/internal/fetch"
)

func server(t *testing.T, routes map[string]string) *httptest.Server {
	t.Helper()
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for prefix, body := range routes {
			if strings.HasPrefix(r.URL.Path, prefix) {
				if r.Method == http.MethodPost {
					b, _ := io.ReadAll(r.Body)
					if !strings.Contains(string(b), `"limit":20`) {
						http.Error(w, "bad body", 400)
						return
					}
				}
				io.WriteString(w, body)
				return
			}
		}
		http.NotFound(w, r)
	}))
	t.Cleanup(s.Close)
	return s
}

func TestParsers(t *testing.T) {
	s := server(t, map[string]string{
		"/v1/boards/acme/jobs": `{"jobs":[{"id":1,"title":"ML Engineer","absolute_url":"https://x/1","first_published":"2026-11-20T10:00:00Z",
			"content":"&lt;p&gt;Earth observation &amp;amp; SAR&lt;/p&gt;","location":{"name":"Rome, Italy"},"departments":[{"name":"AI"}]}]}`,
		"/v0/postings/acme": `[{"id":"a","text":"CV Engineer","hostedUrl":"https://x/a","createdAt":1763600000000,"descriptionPlain":"pytorch",
			"workplaceType":"remote","categories":{"location":"Remote EU","team":"ML"},"salaryRange":{"currency":"EUR","interval":"per-year-salary","min":60000,"max":80000}}]`,
		"/posting-api/job-board/acme": `{"jobs":[{"id":"b","title":"AI Lead","location":"Berlin","publishedAt":"2026-11-01T00:00:00.000+00:00","isListed":true,
			"jobUrl":"https://x/b","descriptionPlain":"geospatial","compensation":{"compensationTierSummary":"€70K – €90K"}},
			{"id":"c","title":"Hidden","isListed":false,"jobUrl":"https://x/c"}]}`,
		"/xml": `<?xml version="1.0"?><workzag-jobs><position><id>42</id><office>Berlin</office><name>Senior ML Engineer</name>
			<createdAt>2026-10-01T09:00:00+00:00</createdAt><jobDescriptions><jobDescription><name>Tasks</name><value><![CDATA[<p>Remote sensing</p>]]></value></jobDescription></jobDescriptions></position></workzag-jobs>`,
		"/wday/cxs/acme/Careers/jobs": `{"total":1,"jobPostings":[{"title":"Data Scientist","externalPath":"/job/Rome/Data-Scientist_R1","locationsText":"Rome","postedOn":"Posted 3 Days Ago"}]}`,
	})
	c := fetch.New(5 * time.Second)
	ctx := context.Background()

	gh, err := ATS["greenhouse"](ctx, c, config.Source{ID: "acme", BaseURL: s.URL}, "Acme")
	if err != nil || len(gh) != 1 || gh[0].Published.Day() != 20 || !strings.Contains(gh[0].Description, "Earth observation & SAR") {
		t.Fatalf("greenhouse: %v %+v", err, gh)
	}
	lv, err := ATS["lever"](ctx, c, config.Source{ID: "acme", BaseURL: s.URL}, "Acme")
	if err != nil || len(lv) != 1 || !lv[0].Remote || lv[0].SalaryMax != 80000 || lv[0].Published.IsZero() {
		t.Fatalf("lever: %v %+v", err, lv)
	}
	as, err := ATS["ashby"](ctx, c, config.Source{ID: "acme", BaseURL: s.URL}, "Acme")
	if err != nil || len(as) != 1 || as[0].Salary == "" {
		t.Fatalf("ashby: %v %+v", err, as)
	}
	ps, err := ATS["personio"](ctx, c, config.Source{ID: "acme", BaseURL: s.URL}, "Acme")
	if err != nil || len(ps) != 1 || ps[0].URL != s.URL+"/job/42" || !strings.Contains(ps[0].Description, "Remote sensing") {
		t.Fatalf("personio: %v %+v", err, ps)
	}
	wd, err := ATS["workday"](ctx, c, config.Source{Domain: "x", Tenant: "acme", Site: "Careers", BaseURL: s.URL, Queries: []string{"data"}}, "Acme")
	if err != nil || len(wd) != 1 || wd[0].Published.IsZero() || wd[0].URL != s.URL+"/Careers/job/Rome/Data-Scientist_R1" {
		t.Fatalf("workday: %v %+v", err, wd)
	}
}

func TestWorkdayPosted(t *testing.T) {
	now := time.Date(2026, 11, 23, 9, 0, 0, 0, time.UTC)
	if d := WorkdayPosted("Posted 30+ Days Ago", now); !d.IsZero() {
		t.Errorf("30+ deve restare ignota, ottenuto %v", d)
	}
	if d := WorkdayPosted("Posted Yesterday", now); d.Day() != 22 {
		t.Errorf("yesterday: %v", d)
	}
	if d := WorkdayPosted("Posted 5 Days Ago", now); d.Day() != 18 {
		t.Errorf("5 days: %v", d)
	}
}

func TestMissingKey(t *testing.T) {
	t.Setenv("ADZUNA_APP_ID", "")
	if _, err := Aggregators["adzuna"](context.Background(), fetch.New(time.Second), config.Search{Query: "x"}); err != ErrMissingKey {
		t.Fatalf("atteso ErrMissingKey, ottenuto %v", err)
	}
}
