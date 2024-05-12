package main

import (
	"net/http"
)

const (
	ROOT string = "."
	PORT string = "8080"
)

func main() {
	serveMux := http.NewServeMux()
	serveMux.Handle("/", http.FileServer(http.Dir(ROOT)))

	httpServer := &http.Server{
		Addr:    "localhost:" + PORT,
		Handler: serveMux,
	}

	httpServer.ListenAndServe()
}
