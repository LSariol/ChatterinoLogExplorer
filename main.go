package main

import (
	"ChatterinoLogExplorer/logs"
	"ChatterinoLogExplorer/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))

	mux.Handle("/", fs)
	mux.HandleFunc("/api/search", searchHandler)

	err := http.ListenAndServe(":9000", mux)
	if err != nil {
		panic(err)
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "use post", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	formData, err := processForm(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	foundLogs := logs.Search(formData)

	w.Header().Set("Content-Type", "text/plain; charset uft-8")
	fmt.Fprintf(w, foundLogs)

}

func processForm(r *http.Request) (models.FormData, error) {

	var processedData models.FormData

	//Process Terms
	rawTerms := strings.Split(r.FormValue("terms"), ",")
	var terms []string
	for _, t := range rawTerms {
		trimmed := strings.TrimSpace(t)
		terms = append(terms, trimmed)
	}

	//Process Duration
	duration, err := strconv.Atoi(r.FormValue("duration"))
	if err != nil {
		return processedData, err
	}

	//Process ExactMatch
	var exactMatch bool
	if r.FormValue("exactMatch") == "1" {
		exactMatch = true
	} else {
		exactMatch = false
	}

	processedData = models.FormData{
		Channel:    r.FormValue("channel"),
		Terms:      terms,
		Duration:   duration,
		User:       r.FormValue("user"),
		ExactMatch: exactMatch,
	}

	return processedData, nil
}
