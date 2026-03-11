package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/huangmatty/crumbs/internal/auth"
	"github.com/huangmatty/crumbs/internal/database"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting JWT: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Couldn't get access token")
		return
	}
	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error validating JWT: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Invalid access token")
		return
	}

	params := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		respondWithError(w, http.StatusBadRequest, "Couldn't decode JSON")
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
	_, err = cfg.db.UpdateUserPassword(r.Context(), database.UpdateUserPasswordParams{
		HashedPassword: hashedPassword,
		ID:             userID,
	})
	if err != nil {
		log.Printf("Error updating user password: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to update password")
		return
	}
	dbUser, err := cfg.db.UpdateUserEmail(r.Context(), database.UpdateUserEmailParams{
		Email: params.Email,
		ID:    userID,
	})
	if err != nil {
		log.Printf("Error updating user email: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to update email")
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
