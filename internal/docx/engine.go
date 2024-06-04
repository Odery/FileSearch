package docx

import (
	"fmt"
	"github.com/fumiama/go-docx"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ProcessSearchRequest is the core function in the whole searching logic.
func ProcessSearchRequest(path, query1, query2 string) (*SearchResult, *SearchProgress, error) {
	result := NewSearchResult()
	progress := new(SearchProgress)
	regex1, regex2, err := compileRegex(query1, query2)

	// Traverse path and subfolders to look for doc, docx files
	err = filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && (strings.HasSuffix(d.Name(), ".docx") || strings.HasSuffix(d.Name(), ".doc")) {
			if !strings.HasPrefix(d.Name(), "~$") {
				progress.Add()
				go processFile(path, regex1, regex2, result, progress)
			}
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return result, progress, err
}

// processFile processes a DOC(X) file, searching for text that matches the provided regular expressions.
// If a match is found, it adds an entry to the search result. This function handles errors in opening,
// reading, and parsing the file, and uses a wait group to manage overall app progress.
func processFile(filepath string, regex1, regex2 *regexp.Regexp, result *SearchResult, progress *SearchProgress) {
	defer progress.Done()

	readFile, err := os.Open(filepath)
	if err != nil {
		log.Println("[ERROR] Error opening a file ", filepath, ":", err)
		return
	}
	defer readFile.Close()

	fileInfo, err := readFile.Stat()
	if err != nil {
		log.Println("[ERROR] Error reading a file ", filepath, ":", err)
		return
	}

	doc, err := docx.Parse(readFile, fileInfo.Size())
	if err != nil {
		log.Println("[ERROR] Error parsing a file as DOC(X) ", filepath, ":", err)
		return
	}

	regex1Bool, regex2Bool := false, false
	for _, element := range doc.Document.Body.Items {
		switch elem := element.(type) {
		case *docx.Paragraph, *docx.Table:
			text := fmt.Sprint(elem)
			if regex1Bool && (regex2 == nil || regex2Bool) {
				result.AddEntry(fileInfo.Name(), filepath, "text", fileInfo.ModTime())
				return
			}

			if regex1Bool == false {
				regex1Bool = regex1.MatchString(text)
			}
			if regex2 != nil {
				regex2Bool = regex2.MatchString(text)
			}
		}
	}
}

// compileRegex constructs regular expressions from input patterns
// and returns them along with an error if the first pattern is empty.
// If a second pattern is provided, it compiles a regular expression for it too.
// Makes a case-insensitive regex
func compileRegex(pattern1, pattern2 string) (*regexp.Regexp, *regexp.Regexp, error) {
	// Input sanitization
	pattern1 = regexp.QuoteMeta(pattern1)
	pattern2 = regexp.QuoteMeta(pattern2)

	if pattern1 == "" {
		return nil, nil, fmt.Errorf("pattern is empty")
	}
	wordsInPattern1 := strings.Split(strings.TrimSpace(pattern1), " ")

	tempRegex := "(?i)"
	for index, value := range wordsInPattern1 {
		if index == len(wordsInPattern1)-1 {
			tempRegex += value + `\p{L}*`
		} else {
			tempRegex += value + `\p{L}*[[:blank:]]+`
		}
	}
	r1, err := regexp.Compile(tempRegex)
	if err != nil {
		return nil, nil, err
	}

	if pattern2 != "" {
		wordsInPattern2 := strings.Split(strings.TrimSpace(pattern2), " ")
		tempRegex2 := "(?i)"
		for index, value := range wordsInPattern2 {
			if index == len(wordsInPattern2)-1 {
				tempRegex2 += value + `\w*`
			} else {
				tempRegex2 += value + `\w*[[:blank:]]+`
			}
		}

		r2, err := regexp.Compile(tempRegex2)
		if err != nil {
			return nil, nil, err
		}
		return r1, r2, nil
	}
	return r1, nil, nil
}
