package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/huangmatty/crumbs/internal/auth"
	"github.com/huangmatty/crumbs/internal/database"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Username *string `json:"username,omitempty"`
		Email    *string `json:"email,omitempty"`
		Password *string `json:"password,omitempty"`
	}{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Couldn't decode JSON", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(cfg.authUserContextKey).(uuid.UUID)
	dbUser, err := cfg.db.GetUserByID(r.Context(), userID)
	if err == sql.ErrNoRows {
		http.Error(w, "User doesn't exist", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		http.Error(w, "Couldn't retrieve user", http.StatusInternalServerError)
		return
	}

	if params.Username != nil {
		dbUser, err = cfg.db.UpdateUsername(r.Context(), database.UpdateUsernameParams{
			Username: *params.Username,
			ID:       userID,
		})
		if err != nil {
			log.Printf("Error updating user: %v", err)
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			return
		}
	}
	if params.Email != nil {
		dbUser, err = cfg.db.UpdateUserEmail(r.Context(), database.UpdateUserEmailParams{
			Email: *params.Email,
			ID:    userID,
		})
		if err != nil {
			log.Printf("Error updating user: %v", err)
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			return
		}
	}
	if params.Password != nil {
		if len(*params.Password) < minPasswordLength {
			http.Error(w, "Password must have at least 12 characters", http.StatusBadRequest)
			return
		}
		hashedPassword, err := auth.HashPassword(*params.Password)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		dbUser, err = cfg.db.UpdateUserPassword(r.Context(), database.UpdateUserPasswordParams{
			HashedPassword: hashedPassword,
			ID:             userID,
		})
		if err != nil {
			log.Printf("Error updating user: %v", err)
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			return
		}
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
