package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error while marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5xx error: %s", msg)
	}
	response := errorResponse{Error: msg}
	respondWithJSON(w, code, response)
}

func replaceBadWords(inputString string, badWords map[string]struct{}) string {
	splittedInput := strings.Split(inputString, " ")

	for i, word := range splittedInput {
		loweredWord := strings.ToLower(word)
		_, ok := badWords[loweredWord]
		if ok {
			splittedInput[i] = "****"
		}
	}

	output := strings.Join(splittedInput, " ")
	return output
}

func validateChirp(body string) (string, error) {
	const maxChirp int = 140
	if len(body) > maxChirp {
		return "", errors.New("Chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	filtered := replaceBadWords(body, badWords)
	return filtered, nil
}

func debugDeleteDatabase(databasePath string) error {
	err := os.Remove(databasePath)
	if err != nil {
		return errors.New("no database to delete")
	}
	_, err = os.ReadFile(databasePath)
	if err != nil {
		return err
	}
	return nil
}
