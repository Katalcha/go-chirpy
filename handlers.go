package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// handler to be used with serveMux.HandleFunc()
// this handler returns a response on with current readiness of the server
func healthzHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(http.StatusText(http.StatusOK)))
}

// handler to be used with serveMux.HandleFunc()
// this handler returns a json response, if incoming POST request is valid
// sends json error msg otherwise
func validateChirpHandler(writer http.ResponseWriter, request *http.Request) {
	const maxChirp = 140

	type chirpType struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		Valid bool `json:"valid"`
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	// set up decoder
	decoder := json.NewDecoder(request.Body)

	// set up struct to fit decoded info into
	chirp := chirpType{}

	// decode info into struct, save status in err
	err := decoder.Decode(&chirp)

	// if the request-body could not be decoded -> prepare error
	if err != nil {
		data, err := json.Marshal(errorResponse{Error: "Could not decode Message"})
		// if the error itself is broken
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write(data)
		return
	}

	// if the request-body is too long -> prepare error
	if len(chirp.Body) > maxChirp {
		data, err := json.Marshal(errorResponse{Error: "Chirp is too long"})
		// if the error itself is broken
		if err != nil {
			log.Printf("Error marshalling JSON %s", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(data)
		return
	}

	// if the request-body lacks a body field
	if chirp.Body == "" {
		data, err := json.Marshal(errorResponse{Error: "Missing 'body field in JSON"})
		if err != nil {
			log.Printf("Error marshalling JSON %s", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(data)
		return
	}

	// if the request-body is valid
	// encode json answer for response
	data, err := json.Marshal(returnVals{Valid: true})
	// if json answer itself is broken
	if err != nil {
		log.Printf("Error marshalling JSON %s", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(data)
}

// handler to be used with serveMux.HandleFunc()
// this handler returns a response with state of visitors
func (a *apiConfig) metricsHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", a.fileServerHits)))
}

// handler to be used with serveMux.HandleFunc()
// this handler returns a response which sets visitor count back to 0
// and shows new state of visitors
func (a *apiConfig) metricsResetHandler(writer http.ResponseWriter, request *http.Request) {
	a.fileServerHits = 0
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(fmt.Sprintf("Hits reset to %d", a.fileServerHits)))
}
