package docx

import (
	"sort"
	"sync"
	"time"
)

type SearchProgress struct {
	sync.Mutex
	done  int
	total int
}

func (s *SearchProgress) IsDone() bool {
	s.Lock()
	defer s.Unlock()
	if s.done == s.total {
		return true
	} else {
		return false
	}
}

func (s *SearchProgress) Done() {
	s.Lock()
	s.done += 1
	s.Unlock()
}

func (s *SearchProgress) Add() {
	s.Lock()
	s.total += 1
	s.Unlock()
}

func (s *SearchProgress) GetDone() int {
	s.Lock()
	defer s.Unlock()
	return s.done
}

// GetProgress will return a value from .0 to 1
func (s *SearchProgress) GetProgress() float64 {
	return float64(s.done) / float64(s.total)
}

type SearchEntry struct {
	SearchedString string
	Path           string
	Name           string
	LastModified   time.Time
}

func (s *SearchEntry) GetFormattedDate() string {
	return s.LastModified.Format("2006-01-02 15:04")
}

// Types for sorting SearchEntry in ascending and descending order:

type byPathAscending []SearchEntry

func (a byPathAscending) Len() int           { return len(a) }
func (a byPathAscending) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPathAscending) Less(i, j int) bool { return a[i].Path < a[j].Path }

type byPathDescending []SearchEntry

func (a byPathDescending) Len() int           { return len(a) }
func (a byPathDescending) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPathDescending) Less(i, j int) bool { return a[i].Path > a[j].Path }

type byNameAscending []SearchEntry

func (a byNameAscending) Len() int           { return len(a) }
func (a byNameAscending) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byNameAscending) Less(i, j int) bool { return a[i].Name < a[j].Name }

type byNameDescending []SearchEntry

func (a byNameDescending) Len() int           { return len(a) }
func (a byNameDescending) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byNameDescending) Less(i, j int) bool { return a[i].Name > a[j].Name }

type byLastModifiedAscending []SearchEntry

func (a byLastModifiedAscending) Len() int      { return len(a) }
func (a byLastModifiedAscending) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byLastModifiedAscending) Less(i, j int) bool {
	return a[i].LastModified.Before(a[j].LastModified)
}

type byLastModifiedDescending []SearchEntry

func (a byLastModifiedDescending) Len() int      { return len(a) }
func (a byLastModifiedDescending) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byLastModifiedDescending) Less(i, j int) bool {
	return a[i].LastModified.After(a[j].LastModified)
}

type SearchResult struct {
	sync.Mutex
	Results []SearchEntry
}

func NewSearchResult() *SearchResult {
	return &SearchResult{
		Results: make([]SearchEntry, 0),
	}
}

// AddEntry updates Search Result with the latest given entry
// concurrency safe
func (s *SearchResult) AddEntry(name, path, searchedString string, lastModified time.Time) {
	s.Lock()
	s.Results = append(s.Results, SearchEntry{
		Name:           name,
		Path:           path,
		SearchedString: searchedString,
		LastModified:   lastModified,
	})
	s.Unlock()
}

func (s *SearchResult) GetEntry(id int) *SearchEntry {
	if id > len(s.Results) {
		return &s.Results[id]
	} else {
		return nil
	}
}

func (s *SearchResult) SortByNameAscending() {
	sort.Sort(byNameAscending(s.Results))
}

func (s *SearchResult) SortByNameDescending() {
	sort.Sort(byNameDescending(s.Results))
}

func (s *SearchResult) SortByLastModifiedAscending() {
	sort.Sort(byLastModifiedAscending(s.Results))
}

func (s *SearchResult) SortByLastModifiedDescending() {
	sort.Sort(byLastModifiedDescending(s.Results))
}

func (s *SearchResult) SortByPathAscending() {
	sort.Sort(byPathAscending(s.Results))
}

func (s *SearchResult) SortByPathDescending() {
	sort.Sort(byPathDescending(s.Results))
}
