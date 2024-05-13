package main

import (
	"log"
	"net/http"
)

const (
	ROOT_PATH string = "."
	PORT      string = "8080"
	HEALTHZ   string = "/healthz"
)

func main() {
	serveMux := http.NewServeMux()
	serveMux.Handle("/app/*", http.StripPrefix("/app", http.FileServer(http.Dir(ROOT_PATH))))
	serveMux.HandleFunc(HEALTHZ, healthzHandler)

	httpServer := &http.Server{
		Addr:    "localhost:" + PORT,
		Handler: serveMux,
	}

	log.Printf("Serving Yo Mama from %s on port: %s\n", ROOT_PATH, PORT)
	log.Fatal(httpServer.ListenAndServe())
}

func healthzHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(http.StatusText(http.StatusOK)))
}
