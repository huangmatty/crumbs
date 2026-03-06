package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/huangmatty/crumbs/internal/database"
)

const maxNameLength = 100

func (cfg *apiConfig) handlerTalentsCreate(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Name   string    `json:"name"`
		UserID uuid.UUID `json:"user_id"`
	}{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding name: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode name")
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
		UserID: params.UserID,
	})
	if err != nil {
		log.Printf("Error creating talent: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create talent")
		return
	}
	talent := TalentDTO{
		ID:     dbTalent.ID,
		UserID: dbTalent.UserID,
		Name:   dbTalent.Name,
	}
	if dbTalent.Email.Valid {
		talent.Email = dbTalent.Email.String
	}
	respondWithJSON(w, http.StatusCreated, talent)
}
