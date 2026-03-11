package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/huangmatty/crumbs/internal/auth"
	"github.com/huangmatty/crumbs/internal/database"
)

const (
	maxUsernameLength = 50
	maxEmailLength    = 75
	minPasswordLength = 12
)

type UserDTO struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	AccessToken  string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode JSON")
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
	if len(params.Password) < minPasswordLength {
		respondWithError(w, http.StatusBadRequest, "Password must have at least 12 characters")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Username:       params.Username,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("Error creating user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}
	user := UserDTO{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Username:  dbUser.Username,
		Email:     dbUser.Email,
	}
	respondWithJSON(w, http.StatusCreated, user)
}
