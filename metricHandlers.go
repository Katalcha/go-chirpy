package main

import (
	"fmt"
	"net/http"
)

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

// handler to be used with serveMux.HandleFunc()
// this handler returns a response on with current readiness of the server
func healthzHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(http.StatusText(http.StatusOK)))
}
