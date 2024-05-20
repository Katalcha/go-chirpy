package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/Katalcha/go-chirpy/internal/auth"
	"github.com/Katalcha/go-chirpy/internal/database"
	"github.com/Katalcha/go-chirpy/internal/utils"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func (a *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not decode parameters")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not hash password")
		return
	}

	user, err := a.DB.CreateUser(params.Email, hashedPassword)
	if err != nil {
		if errors.Is(err, database.ErrAlreadyExists) {
			utils.RespondWithError(w, http.StatusConflict, "user already exists")
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "could not create user")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}

func (a *apiConfig) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	dbUsers, err := a.DB.GetUsers()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Could not get users from database")
		return
	}

	users := []database.User{}
	for _, dbUser := range dbUsers {
		users = append(users, database.User{
			ID:             dbUser.ID,
			Email:          dbUser.Email,
			HashedPassword: dbUser.HashedPassword,
		})
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})

	utils.RespondWithJSON(w, http.StatusOK, users)
}

func (a *apiConfig) getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		User
	}

	const matchingPattern string = "userID"
	userIDString := r.PathValue(matchingPattern)
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	dbUser, err := a.DB.GetUserByID(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "No user found in database")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:    dbUser.ID,
			Email: dbUser.Email,
		},
	})
}

func (a *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
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

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not decode parameters")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not hash password")
		return
	}

	userIDInt, err := strconv.Atoi(subject)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not parse user id")
		return
	}

	user, err := a.DB.UpdateUser(userIDInt, params.Email, hashedPassword)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not update user")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}

func (a *apiConfig) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not decode parameters")
		return
	}

	user, err := a.DB.GetUserByEmail(params.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not get user")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "invalid password")
		return
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		a.jwtSecret,
		time.Hour,
	)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not create access jwt")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not create refresh token")
		return
	}

	err = a.DB.SaveRefreshToken(user.ID, refreshToken)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not save refresh token")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}

func (a *apiConfig) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "could not find token")
		return
	}

	user, err := a.DB.UserForRefreshToken(refreshToken)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "could not get user for refresh token")
		return
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		a.jwtSecret,
		time.Hour,
	)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "could not validate token")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}

func (a *apiConfig) revokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "could not find token")
		return
	}

	err = a.DB.RevokeRefreshToken(refreshToken)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not revoke session")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
