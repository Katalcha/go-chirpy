package main

// type returnVals struct {
// 	CleanedBody string `json:"cleaned_body"`
// }

type errorResponse struct {
	Error string `json:"error"`
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}
