package core

import (
	"fmt"
	"net/http"
)

type SearchResult struct {
	Title string
	Url string
	Snippet string
}

func (s *SearchResult) WriteToHtml(w *http.ResponseWriter) {
	fmt.Fprintf(*w,"<div>")
	fmt.Fprintf(*w,"<a href=\"%s\">%s</a>", s.Url, s.Title)
	fmt.Fprintf(*w,"<br>%s", s.Snippet)
	fmt.Fprintf(*w,"</div>")
}
