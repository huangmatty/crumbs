package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/huangmatty/crumbs/internal/auth"
	"github.com/huangmatty/crumbs/internal/database"
)

const defaultTokenDuration = time.Hour

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode JSON")
		return
	}

	dbUser, err := cfg.db.GetUserByUsername(r.Context(), params.Username)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect username or password")
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect username or password")
		return
	}

	accessToken, err := auth.CreateJWT(dbUser.ID, cfg.jwtSecret, defaultTokenDuration)
	if err != nil {
		log.Printf("Error creating JWT: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create JWT")
		return
	}

	refreshToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     auth.MakeRefreshToken(),
		ExpiresAt: time.Now().AddDate(0, 0, 60),
		UserID:    dbUser.ID,
	})
	if err != nil {
		log.Printf("Error creating refresh token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create refresh token")
		return
	}

	user := UserDTO{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Username:     dbUser.Username,
		Email:        dbUser.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
	}
	respondWithJSON(w, http.StatusOK, user)
}
