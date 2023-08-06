package main

import (
	"encoding/json"
	"net/http"
	"time"
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
	/*
		psudo code:
		1. get the url from the request body
		2. validate the url (later)
		3. generate a unique id
		4. save the url and the id in a map
		5. pass the id to the shortener function
		6. return the shortened url
	*/

	var input struct {
		Url string `json:"url"`
	}
	// convert the JSON request body to a struct
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// validate the url
	if !app.isURLValid(input.Url) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// check if url is already in the database
	for _, v := range app.db { // improve with go routines
		if v.Long_url == input.Url {
			// fmt.Fprintf(w, "%+v\n", v.Short_url)
			err := app.writeJSON(w, http.StatusOK, envelope{"shortlinks": v}, nil)
			if err != nil {
				app.logger.Print(err)
				http.Error(w, "The server encountered a problem and could not process your request.", http.StatusInternalServerError)
			}
			return
		}
	}

	generator, err := NewIDGenerator(app.config.workerID) // DONE: add env var for workerID
	if err != nil {
		http.Error(w, "Failed to generate ID", http.StatusInternalServerError)
		return
	}

	id := generator.GenerateID()
	app.logger.Println("Generated ID:", id)
	hash := DecimalToBase62(id)
	app.logger.Println("Hashed ID:", hash)

	// save the url and the id in a map
	s := Shorten{
		Long_url:  input.Url,
		Short_url: shortenerAddress + hash, // combine hash with shorturl address
		CreatedAt: time.Now(),
	}

	// store Shorten in a map
	app.logger.Printf("Add shorturl [%s] to database.", s.Short_url)
	app.db[id] = s

	// Return the shortened url as a string
	// fmt.Fprintf(w, "%+v\n", app.db[id].Short_url)
	app.writeJSON(w, http.StatusOK, envelope{"shortlink": s}, nil)
	if err != nil {
		app.logger.Print(err)
		http.Error(w, "The server encountered a problem and could not process your request.", http.StatusInternalServerError)
	}
}

func (app *application) listUrlHandler(w http.ResponseWriter, r *http.Request) {
	// DONE: add an envelope to the response
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := app.writeJSON(w, http.StatusOK, envelope{"shortlinks": app.db}, nil)
	if err != nil {
		app.logger.Print(err)
		http.Error(w, "The server encountered a problem and could not process your request.", http.StatusInternalServerError)
	}

	// js, err := json.MarshalIndent(app.db, "", "\t")
	// if err != nil {
	// 	http.Error(w, "Failed to generate JSON response", http.StatusInternalServerError)
	// 	return
	// }

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusFound)
	// w.Write(js)
}

func (app *application) redirectUrlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// get the hash from the url
	hash := r.URL.Path[len("/v1/"):]
	// reverse hash to id
	id := Base62ToDecimal(hash)
	// get the url from the map
	url := app.db[id].Long_url
	// 301 redirect to the url
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}
