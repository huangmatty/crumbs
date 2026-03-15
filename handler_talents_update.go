package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/huangmatty/crumbs/internal/database"
)

func (cfg *apiConfig) handlerTalentsUpdateName(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Name string `json:"name"`
	}{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Couldn't decode JSON", http.StatusBadRequest)
		return
	}
	if params.Name == "" {
		http.Error(w, "Missing name", http.StatusBadRequest)
		return
	}
	if len(params.Name) > maxNameLength {
		http.Error(w, "Name is too long", http.StatusBadRequest)
		return
	}

	talentID, err := uuid.Parse(r.PathValue("talentID"))
	if err != nil {
		log.Printf("Error getting talent id: %v", err)
		http.Error(w, "Invalid talent id", http.StatusBadRequest)
		return
	}

	dbTalentUserID, err := cfg.db.GetUserIDForTalent(r.Context(), talentID)
	if err == sql.ErrNoRows {
		http.Error(w, "Talent doesn't exist", http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Printf("Error retrieving talent's user id: %v", err)
		http.Error(w, "Couldn't retrieve talent's user id", http.StatusInternalServerError)
		return
	}

	userID := r.Context().Value(cfg.authUserContextKey).(uuid.UUID)
	if userID != dbTalentUserID {
		http.Error(w, "Cannot update talent", http.StatusForbidden)
		return
	}

	dbTalent, err := cfg.db.UpdateTalentName(r.Context(), database.UpdateTalentNameParams{
		Name: params.Name,
		ID:   talentID,
	})
	if err != nil {
		log.Printf("Error updating talent: %v", err)
		http.Error(w, "Failed to update talent", http.StatusInternalServerError)
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
