package core

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"unicode"
)

/* Auxiliary methods and classes for the Index of documents */
func tokenizer(c rune) bool {
	return !unicode.IsLetter(c) && !unicode.IsNumber(c)
}

type processedDocument struct {
	title string
	url string
	tokenizedText []string
}

type rawDocument struct {
	title string
	body string
	url string
	mimeType string
}

func (r *rawDocument) readFromFile(path string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	r.body = string(content)
}

func (r *rawDocument) toProcessedDocument() processedDocument {
	var p processedDocument
	if r.mimeType == "text/plain" {
		p.title = r.title
		p.url = r.url
		p.tokenizedText = strings.FieldsFunc(r.body, tokenizer)
	} else {
		log.Fatal("Not implemented")
	}
	return p
}

type indexEntry struct {
	// idx for doc
	idx int

	// pos of the term in doc, any one
	// TODO: make it better by storing all pos
	pos int
}

/* The actual Index of documents */
type SearchIndex struct {
	rawDocuments []rawDocument
	processedDocuments []processedDocument
	reverseIndex map[string][]indexEntry
}

func (s *SearchIndex) IndexDir(path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal("Could not read index dir " + path, err)
	}

	for _, f := range files {
		r := rawDocument{}
		r.readFromFile(path + "/" + f.Name())
		r.title = f.Name()
		r.mimeType = "text/plain"

		if absPath, err := filepath.Abs(path + "/" + f.Name()); err == nil {
			r.url = "file://" + absPath
		}

		s.rawDocuments = append(s.rawDocuments, r)
	}

	s.reverseIndex = make(map[string][]indexEntry)
	for idx, r := range s.rawDocuments {
		p := r.toProcessedDocument()
		s.processedDocuments = append(s.processedDocuments, p)
		seenTerms := make(map[string]bool)

		for term_idx, term := range p.tokenizedText {
			if seenTerms[term] {
				continue
			}

			var i indexEntry
			i.idx = idx
			i.pos = term_idx
			s.reverseIndex[term] = append(s.reverseIndex[term], i)
			seenTerms[term] = true
		}
	}
}

func minInt(x,y int) int {
	if x < y {
		return x
	}
	return y
}

func maxInt(x,y int) int {
	if x > y {
		return x
	}
	return y
}

// Let's start with assumption of single term query.
func (s *SearchIndex) Search(query string) []SearchResult{
	tokenizedQuery := strings.FieldsFunc(query, tokenizer)
  if query == "" {
    return nil
  }

	if len(tokenizedQuery) != 1 {
		log.Print("Only supports single term query for now")
    return nil
	}

	allEntries := s.reverseIndex[tokenizedQuery[0]]
	searchResults := []SearchResult{}
	for _, entry := range allEntries {
		var res SearchResult
		res.Title = s.processedDocuments[entry.idx].title
		res.Url = s.processedDocuments[entry.idx].url

		text := s.processedDocuments[entry.idx].tokenizedText
		var snippetSlice []string
		snippetSlice = append(snippetSlice, text[maxInt(0, entry.pos - 10):maxInt(0, entry.pos - 1)]...)
		snippetSlice = append(snippetSlice, "<b>" + text[entry.pos] + "</b>")
		snippetSlice = append(snippetSlice, text[minInt(len(text) - 1, entry.pos + 1):minInt(len(text) - 1, entry.pos + 10)]...)
		// TODO: this should be actual snippet from body.
		res.Snippet = strings.Join(snippetSlice, " ")
		searchResults = append(searchResults, res)
	}

	return searchResults
}
