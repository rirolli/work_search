// Package fetch fornisce un client HTTP con retry e utilità per testo e date.
package fetch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const userAgent = "work_search-bot/1.0 (+https://github.com/rirolli/work_search)"

// Client esegue richieste HTTP con un paio di tentativi su 429/5xx.
type Client struct {
	HTTP    *http.Client
	Retries int
}

// New crea un client con timeout per singola richiesta.
func New(timeout time.Duration) *Client {
	return &Client{HTTP: &http.Client{Timeout: timeout}, Retries: 2}
}

// Do esegue la richiesta e restituisce il corpo. body può essere nil.
func (c *Client) Do(ctx context.Context, method, url string, body []byte, headers map[string]string) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt <= c.Retries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(attempt*2) * time.Second):
			}
		}
		var rd io.Reader
		if body != nil {
			rd = bytes.NewReader(body)
		}
		req, err := http.NewRequestWithContext(ctx, method, url, rd)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("Accept", "application/json, application/xml;q=0.9, */*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9,it;q=0.8")
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		resp, err := c.HTTP.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		data, err := io.ReadAll(io.LimitReader(resp.Body, 64<<20))
		resp.Body.Close()
		if err != nil {
			lastErr = err
			continue
		}
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("HTTP %d", resp.StatusCode)
			continue
		}
		if resp.StatusCode >= 400 {
			return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, snippet(data))
		}
		return data, nil
	}
	return nil, lastErr
}

// GetJSON scarica url e decodifica il JSON in v.
func (c *Client) GetJSON(ctx context.Context, url string, v any) error {
	b, err := c.Do(ctx, http.MethodGet, url, nil, nil)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, v); err != nil {
		return fmt.Errorf("JSON non valido: %w (%s)", err, snippet(b))
	}
	return nil
}

// PostJSON invia body come JSON e decodifica la risposta in v.
func (c *Client) PostJSON(ctx context.Context, url string, body, v any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}
	b, err := c.Do(ctx, http.MethodPost, url, payload, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, v); err != nil {
		return fmt.Errorf("JSON non valido: %w (%s)", err, snippet(b))
	}
	return nil
}

func snippet(b []byte) string {
	s := strings.TrimSpace(string(b))
	if len(s) > 160 {
		s = s[:160] + "…"
	}
	return s
}

var (
	tagRe   = regexp.MustCompile(`(?s)<[^>]*>`)
	spaceRe = regexp.MustCompile(`\s+`)
)

// PlainText converte HTML (anche con entità già codificate) in testo semplice.
func PlainText(s string) string {
	s = html.UnescapeString(s) // Greenhouse codifica l'HTML due volte
	s = tagRe.ReplaceAllString(s, " ")
	s = html.UnescapeString(s)
	return strings.TrimSpace(spaceRe.ReplaceAllString(s, " "))
}

var timeLayouts = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05.0000000",
	"2006-01-02T15:04:05.000Z",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05 MST",
	"2006-01-02 15:04:05 -0700",
	"2006-01-02 15:04:05",
	"2006-01-02",
}

// ParseTime prova i formati più comuni delle API; restituisce zero se non riesce.
func ParseTime(s string) time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}
	}
	for _, l := range timeLayouts {
		if t, err := time.Parse(l, s); err == nil {
			return t.UTC()
		}
	}
	return time.Time{}
}

// FromMillis converte un timestamp Unix in millisecondi.
func FromMillis(ms int64) time.Time {
	if ms <= 0 {
		return time.Time{}
	}
	return time.UnixMilli(ms).UTC()
}

// FormatSalary produce un testo leggibile per un intervallo retributivo.
func FormatSalary(min, max float64, currency, period string) string {
	if min <= 0 && max <= 0 {
		return ""
	}
	f := func(v float64) string {
		if v >= 1000 {
			return fmt.Sprintf("%.0fk", v/1000)
		}
		return fmt.Sprintf("%.0f", v)
	}
	var s string
	switch {
	case min > 0 && max > 0 && min != max:
		s = f(min) + "–" + f(max)
	case max > 0:
		s = f(max)
	default:
		s = f(min)
	}
	s = strings.TrimSpace(s + " " + currency)
	if period != "" {
		s += " " + period
	}
	return s
}
