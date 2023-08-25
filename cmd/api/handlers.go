package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/duffney/shorturl/internal/data"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]string{
		"status":  "available",
		"version": version,
	}

	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to generate JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (app *application) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"error": "Endpoint not found"}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to generate JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write(jsonResponse)
}

func (app *application) shortenMuxHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.listUrlHandler(w, r)
	case http.MethodPost:
		app.shortenUrlHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (app *application) shortenUrlHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Id        int64  `json:"id"`
		Long_url  string `json:"url"`
		Short_url string `json:"short_url"`
	}
	// convert the JSON request body to a struct
	err := app.readJSON(w, r, &input)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// validate the url
	if !app.isURLValid(input.Long_url) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	if app.models.Urls.LongUrlExists(input.Long_url) {
		url, err := app.models.Urls.GetByLongUrl(input.Long_url)
		if err != nil {
			http.Error(w, "Failed to get URL", http.StatusInternalServerError)
			return
		}
		err = app.writeJSON(w, http.StatusOK, envelope{"shortlink": url}, nil)
		if err != nil {
			http.Error(w, "The server encountered a problem and could not process your request.", http.StatusInternalServerError)
		}

		return
	}

	generator, err := NewIDGenerator(app.config.workerID)
	if err != nil {
		http.Error(w, "Failed to generate ID", http.StatusInternalServerError)
		return
	}

	input.Id = generator.GenerateID()
	app.logger.Println("Generated ID:", input.Id)
	hash := DecimalToBase62(input.Id)
	app.logger.Println("Hashed ID:", hash)

	s := &data.Url{
		Id:        input.Id,
		Long_url:  input.Long_url,
		Short_url: shortenerAddress + hash, // combine hash with shorturl address
	}

	app.logger.Printf("Add shorturl [%s] to database.", s.Short_url)
	err = app.models.Urls.Insert(s)
	if err != nil {
		app.logger.Print(err)
		http.Error(w, "The server encountered a problem and could not process your request.", http.StatusInternalServerError)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"shortlink": s}, nil)
	if err != nil {
		http.Error(w, "The server encountered a problem and could not process your request.", http.StatusInternalServerError)
		return
	}
}

func (app *application) redirectUrlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	hash := r.URL.Path[len("/v1/"):]
	app.logger.Println("Hash:", hash)

	// get the hash from the url
	fmt.Println("Hash:", hash)
	// reverse hash to id
	id := Base62ToDecimal(hash)
	fmt.Println("ID:", id)
	url, err := app.models.Urls.GetById(id)
	if err != nil {
		app.logger.Print(err)
		http.Error(w, "The server encountered a problem and could not process your request.", http.StatusInternalServerError)
		return
	}
	// 301 redirect to the url for tracking visits
	http.Redirect(w, r, url.Long_url, http.StatusFound)

	app.models.Urls.IncrementVisits(url.Id)
}

func (app *application) listUrlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Short_url string
		Long_url  string
		// Page         int // implement pagination
		// PageSize     int // implement pagination
		Pager        data.Pager
		Sort         string
		Direction    string
		SortSafeList []string
	}

	qs := r.URL.Query()

	// extract query string values and set defaults when empty
	input.Short_url = app.readString(qs, "short_url", "")
	input.Long_url = app.readString(qs, "long_url", "")

	input.Pager.Page = app.readInt(qs, "page", 1)
	input.Pager.PageSize = app.readInt(qs, "page_size", 20)

	input.Sort = app.readString(qs, "sort", "id") //#TODO: move logic into data/sort
	input.Direction = "ASC"
	input.SortSafeList = []string{"id", "long_url", "short_url", "created_at", "visits", "-id", "-long_url", "-short_url", "-created_at", "-visits"}

	safeSort := false
	// limit := input.PageSize
	// offset := (input.Page - 1) * input.PageSize

	for _, safeValue := range input.SortSafeList {
		if input.Sort == safeValue {

			if strings.HasPrefix(safeValue, "-") {
				input.Direction = "DESC"
				input.Sort = strings.TrimPrefix(input.Sort, "-")
			}

			safeSort = true

			break
		}
	}

	if !safeSort {
		panic("unsafe sort parameter" + input.Sort)
	}

	urls, metadata, err := app.models.Urls.GetAll(input.Long_url, input.Short_url, input.Sort, input.Direction, input.Pager)
	if err != nil {
		app.logger.Print(err)
		http.Error(w, "The server encountered a problem and could not process your request.", http.StatusInternalServerError)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"shortlinks": urls, "metadata": metadata}, nil)
	if err != nil {
		app.logger.Print(err)
		http.Error(w, "The server encountered a problem and could not process your request.", http.StatusInternalServerError)
	}
}
