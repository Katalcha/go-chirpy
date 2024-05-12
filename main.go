package main

import (
	"net/http"
)

const (
	ROOT string = "."
	PORT string = "8080"
	LOGO string = "/assets/logo.png"
)

func main() {
	serveMux := http.NewServeMux()
	serveMux.Handle("/", http.FileServer(http.Dir(ROOT)))
	serveMux.Handle("/assets", http.FileServer(http.Dir(LOGO)))

	httpServer := &http.Server{
		Addr:    "localhost:" + PORT,
		Handler: serveMux,
	}

	httpServer.ListenAndServe()
}
