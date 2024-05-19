package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	"github.com/Katalcha/go-chirpy/internal/utils"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func (a *apiConfig) createChirpHandler(writer http.ResponseWriter, request *http.Request) {
	chirp := Chirp{}
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&chirp)
	if err != nil {
		utils.RespondWithError(writer, http.StatusBadRequest, err.Error())
		return
	}

	filtered, err := utils.ValidateChirp(chirp.Body)
	if err != nil {
		utils.RespondWithError(writer, http.StatusBadRequest, err.Error())
		return
	}

	newChirp, err := a.DB.CreateChirp(filtered)
	if err != nil {
		utils.RespondWithError(writer, http.StatusInternalServerError, "Could not create chirp")
		return
	}

	utils.RespondWithJSON(writer, http.StatusCreated, Chirp{
		ID:   newChirp.ID,
		Body: newChirp.Body,
	})
}

// func validateChirp(body string) (string, error) {
// 	const maxChirpLength = 140
// 	if len(body) > maxChirpLength {
// 		return "", errors.New("Chirp is too long")
// 	}

// 	badWords := map[string]struct{}{
// 		"kerfuffle": {},
// 		"sharbert":  {},
// 		"fornax":    {},
// 	}
// 	cleaned := utils.ReplaceBadWords(body, badWords)
// 	return cleaned, nil
// }

func (a *apiConfig) getChirpsHandler(writer http.ResponseWriter, _ *http.Request) {
	dbChirps, err := a.DB.GetChirps()
	if err != nil {
		utils.RespondWithError(writer, http.StatusInternalServerError, "Could not retrieve chirps from database")
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:   dbChirp.ID,
			Body: dbChirp.Body,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	utils.RespondWithJSON(writer, http.StatusOK, chirps)
}

func (a *apiConfig) getChirpByIdHandler(w http.ResponseWriter, r *http.Request) {
	const matchingPattern string = "chirpID"
	chirpIDString := r.PathValue(matchingPattern)
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Chirp ID")
		return
	}

	dbChirp, err := a.DB.GetChirpByID(chirpID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "No chirp found in database")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, Chirp{
		ID:   dbChirp.ID,
		Body: dbChirp.Body,
	})
}
