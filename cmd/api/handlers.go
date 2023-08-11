package main

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	// err := json.NewDecoder(r.Body).Decode(&input)
	// if err != nil {
	// 	http.Error(w, "Bad request", http.StatusBadRequest)
	// 	return
	// }

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

	// for _, v := range app.dbMap { // improve with go routines
	// 	if v.Long_url == input.Url {
	// 		// fmt.Fprintf(w, "%+v\n", v.Short_url)
	// 		err := app.writeJSON(w, http.StatusOK, envelope{"shortlinks": v}, nil)
	// 		if err != nil {
	// 			app.logger.Print(err)
	// 			http.Error(w, "The server encountered a problem and could not process your request.", http.StatusInternalServerError)
	// 		}
	// 		return
	// 	}
	// }

	generator, err := NewIDGenerator(app.config.workerID)
	if err != nil {
		http.Error(w, "Failed to generate ID", http.StatusInternalServerError)
		return
	}

	input.Id = generator.GenerateID()
	app.logger.Println("Generated ID:", input.Id)
	hash := DecimalToBase62(input.Id)
	app.logger.Println("Hashed ID:", hash)

	// save the url and the id in a map
	s := &data.Url{
		Id:        input.Id,
		Long_url:  input.Long_url,
		Short_url: shortenerAddress + hash, // combine hash with shorturl address
	}
	// s := Shorten{
	// 	Long_url:  input.Url,
	// 	Short_url: shortenerAddress + hash, // combine hash with shorturl address
	// 	CreatedAt: time.Now(),
	// }

	// store Shorten in a map
	app.logger.Printf("Add shorturl [%s] to database.", s.Short_url)
	// app.dbMap[id] = s
	err = app.models.Urls.Insert(s)
	if err != nil {
		app.logger.Print(err)
		http.Error(w, "The server encountered a problem and could not process your request.", http.StatusInternalServerError)
		return
	}

	// Return the shortened url as a string
	// fmt.Fprintf(w, "%+v\n", app.db[id].Short_url)
	// app.writeJSON(w, http.StatusOK, envelope{"shortlink": s}, nil)
	// if err != nil {
	// 	app.logger.Print(err)
	// 	http.Error(w, "The server encountered a problem and could not process your request.", http.StatusInternalServerError)
	// }
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

	// the the has from the shorturl example: http://localhost:4000/v1/7VQTFRr8s
	hash := r.URL.Path[len("/v1/"):]
	app.logger.Println("Hash:", hash)

	// get the hash from the url
	// hash := r.URL.Path[len("/v1/"):]
	fmt.Println("Hash:", hash)
	// reverse hash to id
	id := Base62ToDecimal(hash)
	fmt.Println("ID:", id)
	// get the url from the map
	// url := app.dbMap[id].Long_url
	url, err := app.models.Urls.GetById(id)
	if err != nil {
		app.logger.Print(err)
		http.Error(w, "The server encountered a problem and could not process your request.", http.StatusInternalServerError)
		return
	}
	// 301 redirect to the url
	http.Redirect(w, r, url.Long_url, http.StatusMovedPermanently)
}

func (app *application) listUrlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	urls, err := app.models.Urls.GetAll()
	if err != nil {
		app.logger.Print(err)
		http.Error(w, "The server encountered a problem and could not process your request.", http.StatusInternalServerError)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"shortlinks": urls}, nil)
	if err != nil {
		app.logger.Print(err)
		http.Error(w, "The server encountered a problem and could not process your request.", http.StatusInternalServerError)
	}
}
