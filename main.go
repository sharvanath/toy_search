package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/sharvanath/toy_search/core"
)

var searchIndex core.SearchIndex

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<title>ToySearch</title>")
	queryStr := r.URL.Query().Get("q")

	fmt.Fprintf(w, "<h2>Results for %s: </h2>", queryStr)
	for _, r := range searchIndex.Search(queryStr) {
		r.WriteToHtml(&w)
		fmt.Fprintf(w, "<br>")
	}
}

func main() {
	searchIndex.IndexDir("./datasets/space")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

