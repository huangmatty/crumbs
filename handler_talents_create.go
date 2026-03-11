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

const maxNameLength = 100

type TalentDTO struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerTalentsCreate(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting user token: %v", err)
		respondWithError(w, http.StatusBadRequest, "Couldn't get user token")
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error validating token: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	params := struct {
		Name string `json:"name"`
	}{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode JSON")
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

	dbTalent, err := cfg.db.CreateTalent(r.Context(), database.CreateTalentParams{
		Name:   params.Name,
		UserID: userID,
	})
	if err != nil {
		log.Printf("Error creating talent: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create talent")
		return
	}
	talent := TalentDTO{
		ID:        dbTalent.ID,
		CreatedAt: dbTalent.CreatedAt,
		UpdatedAt: dbTalent.UpdatedAt,
		Name:      dbTalent.Name,
		UserID:    dbTalent.UserID,
	}
	if dbTalent.Email.Valid {
		talent.Email = dbTalent.Email.String
	}
	respondWithJSON(w, http.StatusCreated, talent)
}
