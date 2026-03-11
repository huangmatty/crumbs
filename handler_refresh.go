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
		respondWithError(w, http.StatusBadRequest, "Couldn't get refresh token")
		return
	}
	dbRefreshToken, err := cfg.db.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}
	if dbRefreshToken.RevokedAt.Valid || dbRefreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token has expired or has been revoked")
		return
	}

	dbUser, err := cfg.db.GetUserFromRefreshToken(r.Context(), dbRefreshToken.Token)
	if err != nil {
		log.Printf("Error getting user from refresh token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user from refresh token")
		return
	}
	token, err := auth.CreateJWT(dbUser.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		log.Printf("Error creating JWT: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create JWT")
		return
	}
	respondWithJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{Token: token})
}
