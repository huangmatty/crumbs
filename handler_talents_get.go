package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerTalentsGet(w http.ResponseWriter, r *http.Request) {
	talentIDStr := r.PathValue("talentID")
	talentID, err := uuid.Parse(talentIDStr)
	if err != nil {
		log.Printf("Error parsing talent id string: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid talent id")
		return
	}

	dbTalent, err := cfg.db.GetTalentByID(r.Context(), talentID)
	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusNotFound, "Talent doesn't exist")
		return
	}
	if err != nil {
		log.Printf("Error retrieving talent: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve talent")
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
	respondWithJSON(w, http.StatusOK, talent)
}
