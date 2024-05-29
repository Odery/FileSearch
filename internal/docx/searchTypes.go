package docx

import (
	"sync"
	"time"
)

type SearchEntry struct {
	SearchedString string
	Path           string
	Name           string
	LastModified   time.Time
}

type SearchResult struct {
	sync.Mutex
	results   map[int]SearchEntry
	currentID int
}

func NewSearchResult() *SearchResult {
	return &SearchResult{
		results:   make(map[int]SearchEntry),
		currentID: 0,
	}
}

// AddEntry updates Search Result with the latest given entry
// concurrency safe
func (s *SearchResult) AddEntry(name, path, searchedString string, lastModified time.Time) {
	s.Lock()
	s.results[s.currentID] = SearchEntry{
		Name:           name,
		Path:           path,
		SearchedString: searchedString,
		LastModified:   lastModified,
	}
	s.currentID += 1
	s.Unlock()
}

func (s *SearchResult) GetEntry(id int) (SearchEntry, bool) {
	v, ok := s.results[id]
	return v, ok
}
