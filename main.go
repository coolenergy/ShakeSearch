package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

var cache = struct {
	sync.RWMutex
	m map[string][]string
}{m: make(map[string][]string)}

func main() {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", handleSearch(searcher))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("Listening on port %s...", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

type Searcher struct {
	CompleteWorks string
	SuffixArray   *suffixarray.Index
}

type PaginatedResults struct {
	Results     []string `json:"results"`
	CurrentPage int      `json:"current_page"`
	TotalPages  int      `json:"total_pages"`
	PerPage     int      `json:"per_page"`
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}

		pageStr := r.URL.Query().Get("page")
		perPageStr := r.URL.Query().Get("perPage")

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		perPage, err := strconv.Atoi(perPageStr)
		if err != nil || perPage < 1 {
			perPage = 10
		}

		results := searcher.Search(strings.ToLower(query[0]))
		totalPages := (len(results) + perPage - 1) / perPage

		start := (page - 1) * perPage
		end := start + perPage
		if end > len(results) {
			end = len(results)
		}

		paginatedResults := PaginatedResults{
			Results:     results[start:end],
			CurrentPage: page,
			TotalPages:  totalPages,
			PerPage:     perPage,
		}

		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err = enc.Encode(paginatedResults)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("encoding failure"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	}
}

func (s *Searcher) Load(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	s.CompleteWorks = string(content)
	s.SuffixArray = suffixarray.New(content)

	return nil
}

func (s *Searcher) Search(query string) []string {
	reg := regexp.MustCompile(fmt.Sprintf(`(?i)(\b%s[\w]*\b)`, query))
	matches := s.SuffixArray.FindAllIndex(reg, -1)

	results := []string{}
	for _, match := range matches {
		start := match[0]
		end := match[1]

		excerptStart := strings.LastIndexFunc(s.CompleteWorks[:start], func(c rune) bool {
			return unicode.IsUpper(c) || c == '.'
		}) + 1
		if excerptStart < 0 {
			excerptStart = 0
		}

		excerptEnd := strings.Index(s.CompleteWorks[end:], ".") + end
		if excerptEnd < 0 {
			excerptEnd = len(s.CompleteWorks)
		}

		results = append(results, s.CompleteWorks[excerptStart:excerptEnd])
	}

	return results
}
