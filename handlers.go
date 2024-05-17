package main

import (
	"encoding/json"
	"fmt"
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
	// maximum characters per chirp
	const maxChirp int = 140
	// bad words filter list
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	// set up struct to fit decoded info into
	chirp := chirpType{}

	// set up decoder
	decoder := json.NewDecoder(request.Body)
	// decode info into struct, save status in err
	err := decoder.Decode(&chirp)

	// if the request-body could not be decoded -> prepare error
	if err != nil {
		// respond with error
		respondWithError(writer, http.StatusInternalServerError, "Could not decode Message")
		return
	}

	// if the request-body is too long -> prepare error
	if len(chirp.Body) > maxChirp {
		// respond with error
		respondWithError(writer, http.StatusBadRequest, "Chirp is too long")
		return
	}

	// if the request-body lacks a body field -> prepare error
	if chirp.Body == "" {
		// respond with error
		respondWithError(writer, http.StatusBadRequest, "Empty body field or missing body field in JSON")
		return
	}

	// if the request-body is valid
	// filter for bad words und replace them
	filteredBody := replaceBadWords(chirp.Body, badWords)
	// respond with JSON
	respondWithJSON(writer, http.StatusOK, returnVals{CleanedBody: filteredBody})
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
