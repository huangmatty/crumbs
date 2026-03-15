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
		http.Error(w, "Couldn't decode JSON", http.StatusBadRequest)
		return
	}
	if params.Username == "" {
		http.Error(w, "Missing username", http.StatusBadRequest)
		return
	}
	if len(params.Username) > maxUsernameLength {
		http.Error(w, "Username is too long", http.StatusBadRequest)
		return
	}
	if params.Email == "" {
		http.Error(w, "Missing email address", http.StatusBadRequest)
		return
	}
	if len(params.Email) > maxEmailLength {
		http.Error(w, "Email address is too long", http.StatusBadRequest)
		return
	}
	if len(params.Password) < minPasswordLength {
		http.Error(w, "Password must have at least 12 characters", http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Username:       params.Username,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
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
