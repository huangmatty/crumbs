package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/huangmatty/crumbs/internal/auth"
	"github.com/huangmatty/crumbs/internal/database"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Couldn't decode JSON", http.StatusBadRequest)
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

	userID := r.Context().Value(cfg.authUserContextKey).(uuid.UUID)
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	_, err = cfg.db.UpdateUserPassword(r.Context(), database.UpdateUserPasswordParams{
		HashedPassword: hashedPassword,
		ID:             userID,
	})
	if err != nil {
		log.Printf("Error updating user password: %v", err)
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}
	dbUser, err := cfg.db.UpdateUserEmail(r.Context(), database.UpdateUserEmailParams{
		Email: params.Email,
		ID:    userID,
	})
	if err != nil {
		log.Printf("Error updating user email: %v", err)
		http.Error(w, "Failed to update email address", http.StatusInternalServerError)
		return
	}

	user := UserDTO{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Username:  dbUser.Username,
		Email:     dbUser.Email,
	}
	respondWithJSON(w, http.StatusOK, user)
}
