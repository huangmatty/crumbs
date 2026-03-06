package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/huangmatty/crumbs/internal/database"
)

const (
	maxUsernameLength = 50
	maxEmailLength    = 75
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding username and email: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode username and email")
		return
	}
	if params.Username == "" {
		respondWithError(w, http.StatusBadRequest, "Missing username")
		return
	}
	if len(params.Username) > maxUsernameLength {
		respondWithError(w, http.StatusBadRequest, "Username is too long")
		return
	}
	if params.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Missing email")
		return
	}
	if len(params.Email) > maxEmailLength {
		respondWithError(w, http.StatusBadRequest, "Email is too long")
		return
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Username: params.Username,
		Email:    params.Email,
	})
	if err != nil {
		log.Printf("Error creating user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}
	user := User{
		Username: dbUser.Username,
		Email:    dbUser.Email,
	}
	respondWithJSON(w, http.StatusCreated, user)
}
