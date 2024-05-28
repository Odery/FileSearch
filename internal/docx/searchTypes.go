package docx

import "sync"

type SearchEntry struct {
	SearchedString string
	Path           string
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
func (s *SearchResult) AddEntry(path, searchedString string) {
	s.Lock()
	s.results[s.currentID] = SearchEntry{SearchedString: searchedString, Path: path}
	s.currentID += 1
	s.Unlock()
}

func (s *SearchResult) GetEntry(id int) (SearchEntry, bool) {
	v, ok := s.results[id]
	return v, ok
}
