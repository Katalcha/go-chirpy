package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	"github.com/Katalcha/go-chirpy/internal/utils"
)

func (a *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	user := User{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	newUser, err := a.DB.CreateUser(user.Email, user.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Could not create user")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, User{
		ID:    newUser.ID,
		Email: newUser.Email,
	})
}

func (a *apiConfig) getUsersHandler(w http.ResponseWriter, _ *http.Request) {
	dbUsers, err := a.DB.GetUsers()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Could not retrieve users from database")
	}

	users := []User{}
	for _, dbUser := range dbUsers {
		users = append(users, User{
			ID:    dbUser.ID,
			Email: dbUser.Email,
		})
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})

	utils.RespondWithJSON(w, http.StatusOK, users)
}

func (a *apiConfig) getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
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

	utils.RespondWithJSON(w, http.StatusOK, User{
		ID:    dbUser.ID,
		Email: dbUser.Email,
	})
}
