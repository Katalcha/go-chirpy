package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strconv"

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

func (a *apiConfig) loginUserHandler(w http.ResponseWriter, r *http.Request) {
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

	utils.RespondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}