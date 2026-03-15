package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerTalentsGet(w http.ResponseWriter, r *http.Request) {
	talentID, err := uuid.Parse(r.PathValue("talentID"))
	if err != nil {
		log.Printf("Error parsing talent id: %v", err)
		http.Error(w, "Invalid talent id", http.StatusBadRequest)
		return
	}

	dbTalent, err := cfg.db.GetTalentByID(r.Context(), talentID)
	if err == sql.ErrNoRows {
		http.Error(w, "Talent doesn't exist", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error retrieving talent: %v", err)
		http.Error(w, "Couldn't retrieve talent", http.StatusInternalServerError)
		return
	}

	userID := r.Context().Value(cfg.authUserContextKey).(uuid.UUID)
	if userID != dbTalent.UserID {
		http.Error(w, "Cannot access talent", http.StatusForbidden)
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
	respondWithJSON(w, http.StatusOK, talent)
}
