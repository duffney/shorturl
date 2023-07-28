package main

import (
	"encoding/json"
	"net/http"
	"net/url"
)

/*
DONE: Create `writeJson` and `readJson` helpers
	DONE: Update `fmt.Fprintf` statements to use json helpers
DONE: isUrlValid
*/

type envelope map[string]any

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) isURLValid(urlStr string) bool {
	_, err := url.ParseRequestURI(urlStr)
	return err == nil
}
