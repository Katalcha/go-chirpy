package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Katalcha/go-chirpy/internal/auth"
	"github.com/Katalcha/go-chirpy/internal/database"
	"github.com/Katalcha/go-chirpy/internal/utils"
)

func (a *apiConfig) webhookhandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		}
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "could not find api key")
		return
	}

	if apiKey != a.polkaKey {
		utils.RespondWithError(w, http.StatusUnauthorized, "api key is invalid")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not decode parameters")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = a.DB.UpgradeChirpyRed(params.Data.UserID)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			utils.RespondWithError(w, http.StatusNotFound, "could not find user")
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "could not update user")
		return
	}

	w.WriteHeader(http.StatusOK)
}
