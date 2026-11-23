// Package state ricorda le offerte già viste tra un'esecuzione e l'altra.
package state

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/rirolli/work_search/internal/model"
)

// Entry è il ricordo di un'offerta.
type Entry struct {
	Company   string    `json:"company"`
	Title     string    `json:"title"`
	Location  string    `json:"location"`
	URL       string    `json:"url"`
	Origin    string    `json:"origin"`
	Priority  string    `json:"priority"`
	Score     int       `json:"score"`
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
}

// State è il contenuto di data/state.json.
type State struct {
	LastRun time.Time         `json:"last_run"`
	Jobs    map[string]*Entry `json:"jobs"`
}

// Load legge lo stato; un file assente produce uno stato vuoto.
func Load(path string) (*State, error) {
	s := &State{Jobs: map[string]*Entry{}}
	b, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, s); err != nil {
		return nil, err
	}
	if s.Jobs == nil {
		s.Jobs = map[string]*Entry{}
	}
	return s, nil
}

// Merge aggiorna lo stato con le offerte di questa esecuzione e imposta
// FirstSeen/New su ciascuna. Le offerte non più presenti vengono rimosse.
// Fanno eccezione quelle di una fonte che in questa esecuzione è andata in
// errore (failed): restano in memoria, così quando la fonte torna a rispondere
// non vengono segnate di nuovo come nuove. Anche loro vengono rimosse se la
// fonte non le restituisce per più di retention.
func (s *State) Merge(jobs []model.Job, now time.Time, newWithin time.Duration, failed map[string]bool, retention time.Duration) {
	current := make(map[string]bool, len(jobs))
	for i := range jobs {
		j := &jobs[i]
		current[j.Key] = true
		e, ok := s.Jobs[j.Key]
		if !ok {
			e = &Entry{FirstSeen: now}
			s.Jobs[j.Key] = e
		}
		e.Company, e.Title, e.Location, e.URL, e.Origin = j.Company, j.Title, j.Location, j.URL, j.Origin
		e.Priority, e.Score, e.LastSeen = j.Priority, j.Score, now
		j.FirstSeen = e.FirstSeen
		j.New = now.Sub(e.FirstSeen) <= newWithin
	}
	for k, e := range s.Jobs {
		if current[k] {
			continue
		}
		if failed[e.Origin] && now.Sub(e.LastSeen) <= retention {
			continue
		}
		delete(s.Jobs, k)
	}
	s.LastRun = now
}

// Save scrive lo stato in modo atomico.
func (s *State) Save(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
