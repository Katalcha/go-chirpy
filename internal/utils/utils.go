package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type errorResponse struct {
	Error string `json:"error"`
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
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

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5xx error: %s", msg)
	}
	response := errorResponse{Error: msg}
	RespondWithJSON(w, code, response)
}

func ReplaceBadWords(inputString string, badWords map[string]struct{}) string {
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

func ValidateChirp(body string) (string, error) {
	const maxChirp int = 140
	if len(body) > maxChirp {
		return "", errors.New("Chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	filtered := ReplaceBadWords(body, badWords)
	return filtered, nil
}

func DebugDeleteDatabase(databasePath string) error {
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

func HashThisPw(pw string) ([]byte, error) {
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(pw), 10)
	if err != nil {
		log.Printf("hashing password failed")
		return nil, errors.New("hasing password failed")
	}
	return hashedPw, nil
}
