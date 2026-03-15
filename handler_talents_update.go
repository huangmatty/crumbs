package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/huangmatty/crumbs/internal/database"
)

func (cfg *apiConfig) handlerTalentsUpdate(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Name  *string `json:"name,omitempty"`
		Email *string `json:"email,omitempty"`
	}{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Couldn't decode JSON", http.StatusBadRequest)
		return
	}

	talentID, err := uuid.Parse(r.PathValue("talentID"))
	if err != nil {
		log.Printf("Error getting talent id: %v", err)
		http.Error(w, "Invalid talent id", http.StatusBadRequest)
		return
	}

	dbTalent, err := cfg.db.GetTalentByID(r.Context(), talentID)
	if err == sql.ErrNoRows {
		http.Error(w, "Talent doesn't exist", http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Printf("Error retrieving talent: %v", err)
		http.Error(w, "Couldn't retrieve talent", http.StatusInternalServerError)
		return
	}

	userID := r.Context().Value(cfg.authUserContextKey).(uuid.UUID)
	if userID != dbTalent.UserID {
		http.Error(w, "Cannot update talent", http.StatusForbidden)
		return
	}

	if params.Name != nil {
		if len(*params.Name) > maxNameLength {
			http.Error(w, "Name is too long", http.StatusBadRequest)
			return
		}
		dbTalent, err = cfg.db.UpdateTalentName(r.Context(), database.UpdateTalentNameParams{
			Name: *params.Name,
			ID:   talentID,
		})
		if err != nil {
			log.Printf("Error updating talent: %v", err)
			http.Error(w, "Failed to update talent", http.StatusInternalServerError)
			return
		}
	}
	if params.Email != nil {
		dbTalent, err = cfg.db.UpdateTalentEmail(r.Context(), database.UpdateTalentEmailParams{
			Email: sql.NullString{String: *params.Email, Valid: true},
			ID:    talentID,
		})
		if err != nil {
			log.Printf("Error updating talent: %v", err)
			http.Error(w, "Failed to update talent", http.StatusInternalServerError)
			return
		}
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
