package main

import (
	"log"
	"net/http"
	"time"

	"github.com/huangmatty/crumbs/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting refresh token: %v", err)
		http.Error(w, "Didn't get refresh token", http.StatusUnauthorized)
		return
	}
	dbRefreshToken, err := cfg.db.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}
	if dbRefreshToken.RevokedAt.Valid || dbRefreshToken.ExpiresAt.Before(time.Now()) {
		http.Error(w, "Refresh token has expired or has been revoked", http.StatusUnauthorized)
		return
	}

	dbUser, err := cfg.db.GetUserFromRefreshToken(r.Context(), dbRefreshToken.Token)
	if err != nil {
		log.Printf("Error getting user from refresh token: %v", err)
		http.Error(w, "Couldn't get user from refresh token", http.StatusInternalServerError)
		return
	}
	token, err := auth.CreateJWT(dbUser.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		log.Printf("Error creating JWT: %v", err)
		http.Error(w, "Failed to create access token", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{Token: token})
}
