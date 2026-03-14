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

type BuyerDTO struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerBuyersCreate(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting JWT: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Couldn't get access token")
		return
	}
	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error validating JWT")
		respondWithError(w, http.StatusUnauthorized, "Invalid access token")
		return
	}

	params := struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		respondWithError(w, http.StatusBadRequest, "Couldn't decode JSON")
		return
	}
	if params.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Missing name")
		return
	}
	if len(params.Name) > maxNameLength {
		respondWithError(w, http.StatusBadRequest, "Name is too long")
		return
	}
	if params.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Missing email")
		return
	}

	dbBuyer, err := cfg.db.CreateBuyer(r.Context(), database.CreateBuyerParams{
		Name:   params.Name,
		Email:  params.Email,
		UserID: userID,
	})
	if err != nil {
		log.Printf("Error creating buyer: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create buyer")
		return
	}

	buyer := BuyerDTO{
		ID:        dbBuyer.ID,
		CreatedAt: dbBuyer.CreatedAt,
		UpdatedAt: dbBuyer.UpdatedAt,
		Name:      dbBuyer.Name,
		Email:     dbBuyer.Email,
		UserID:    dbBuyer.UserID,
	}
	respondWithJSON(w, http.StatusCreated, buyer)
}
