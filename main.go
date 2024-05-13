package main

import (
	"log"
	"net/http"
)

const (
	ROOT_PATH string = "."
	PORT      string = "8080"
	HEALTHZ   string = "/healthz"
	METRICS   string = "/metrics"
	RESET     string = "/reset"
)

func main() {
	// create apiConfig state object
	apiCfg := apiConfig{
		fileServerHits: 0,
	}

	// create http server multiplexer
	serveMux := http.NewServeMux()

	// define file server
	serveMux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(ROOT_PATH)))))

	// let multiplexer handle specific endpoints for...
	// on HEALTHZ endpoint call, return readiness status
	serveMux.HandleFunc(HEALTHZ, healthzHandler)
	// on METRICS endpoint call, return visitor count
	serveMux.HandleFunc(METRICS, apiCfg.metricsHandler)
	// on RESET endpoint call, reset visitor count
	serveMux.HandleFunc(RESET, apiCfg.metricsResetHandler)

	// create http.Server object
	httpServer := &http.Server{
		Addr:    "localhost:" + PORT,
		Handler: serveMux,
	}

	// log info before server start
	log.Printf("Serving Yo Mama from %s on port: %s\n", ROOT_PATH, PORT)
	// log fatal errors and start server
	log.Fatal(httpServer.ListenAndServe())
}
