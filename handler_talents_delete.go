package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerTalentsDelete(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Cannot delete talent", http.StatusForbidden)
		return
	}

	_, err = cfg.db.SoftDeleteTalent(r.Context(), talentID)
	if err != nil {
		log.Printf("Error soft-deleting talent: %v", err)
		http.Error(w, "Failed to delete talent", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
