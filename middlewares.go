package main

import "net/http"

// returns a http.Handler in by use of http.HandlerFunc()
// this Handler is called and increases apiConfig.fileServerHits by 1
// this function should be used as middleware or wrapper for http.FileServer()
func (a *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		a.fileServerHits++
		next.ServeHTTP(writer, request)
	})
}
