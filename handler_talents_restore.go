package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerTalentsRestore(w http.ResponseWriter, r *http.Request) {
	talentIDStr := r.PathValue("talentID")
	talentID, err := uuid.Parse(talentIDStr)
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
	}

	userID := r.Context().Value(cfg.authUserContextKey)
	if userID != dbTalent.UserID {
		http.Error(w, "Cannot restore talent", http.StatusForbidden)
		return
	}

	dbTalent, err = cfg.db.RestoreTalent(r.Context(), talentID)
	if err != nil {
		log.Printf("Error restoring talent: %v", err)
		http.Error(w, "Failed to restore talent", http.StatusInternalServerError)
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
