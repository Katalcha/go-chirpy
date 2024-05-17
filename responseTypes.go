package main

type chirpType struct {
	Body string `json:"body"`
}

type returnVals struct {
	CleanedBody string `json:"cleaned_body"`
}

type errorResponse struct {
	Error string `json:"error"`
}
