package main

import (
	"log"
	"net/http"
)

const (
	PORT               string = "8080"
	FILE_ROOT_PATH     string = "."
	FILE_SERVER_PATH   string = "/app/*"
	API_VALIDATE_CHIRP string = "/api/validate_chirp"
	API_HEALTHZ        string = "/api/healthz"
	API_RESET          string = "/api/reset"
	ADMIN_METRICS      string = "/admin/metrics"
)

const (
	GET  string = "GET "
	POST string = "POST "
)

func main() {
	// create apiConfig state object
	apiCfg := apiConfig{
		fileServerHits: 0,
	}

	// create http server multiplexer
	serveMux := http.NewServeMux()

	// define file server
	serveMux.Handle(FILE_SERVER_PATH, apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(FILE_ROOT_PATH)))))

	// let multiplexer handle specific endpoints for...
	// on HEALTHZ endpoint call, return readiness status
	serveMux.HandleFunc(GET+API_HEALTHZ, healthzHandler)
	// on METRICS endpoint call, return visitor count
	serveMux.HandleFunc(GET+ADMIN_METRICS, apiCfg.metricsHandler)
	// on RESET endpoint call, reset visitor count
	serveMux.HandleFunc(GET+API_RESET, apiCfg.metricsResetHandler)
	// on VALIDATE_CHIRP call, response with json
	serveMux.HandleFunc(POST+API_VALIDATE_CHIRP, validateChirpHandler)

	// create http.Server object
	httpServer := &http.Server{
		Addr:    "localhost:" + PORT,
		Handler: serveMux,
	}

	// log info before server start
	log.Printf("Serving Yo Mama from %s on port: %s\n", FILE_ROOT_PATH, PORT)
	// log fatal errors and start server
	log.Fatal(httpServer.ListenAndServe())
}
