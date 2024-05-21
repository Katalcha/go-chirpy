package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	"github.com/Katalcha/go-chirpy/internal/auth"
	"github.com/Katalcha/go-chirpy/internal/utils"
)

type Chirp struct {
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
	Body     string `json:"body"`
}

func (a *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "could not find jwt")
		return
	}

	subject, err := auth.ValidateJWT(token, a.jwtSecret)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "could not validate jwt")
		return
	}

	userID, err := strconv.Atoi(subject)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "coul not parse user id")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not decode parameters")
		return
	}

	cleaned, err := utils.ValidateChirp(params.Body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := a.DB.CreateChirp(cleaned, userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not create chirp")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, Chirp{
		ID:       chirp.ID,
		AuthorID: chirp.AuthorID,
		Body:     chirp.Body,
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
			ID:       dbChirp.ID,
			AuthorID: dbChirp.AuthorID,
			Body:     dbChirp.Body,
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
		ID:       dbChirp.ID,
		AuthorID: dbChirp.AuthorID,
		Body:     dbChirp.Body,
	})
}

func (a *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	const matchingPattern string = "chirpID"
	chirpIDString := r.PathValue(matchingPattern)
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid chirp id")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "could not find jwt")
		return
	}

	subject, err := auth.ValidateJWT(token, a.jwtSecret)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "could not validate jwt")
		return
	}

	userID, err := strconv.Atoi(subject)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "could not parse user id")
		return
	}

	dbChirp, err := a.DB.GetChirpByID(chirpID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "could not find chirp")
		return
	}

	if dbChirp.AuthorID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "you cannot delete this chirp")
		return
	}

	err = a.DB.DeleteChirp(chirpID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not delete chirp")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
