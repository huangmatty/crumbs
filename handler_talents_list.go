package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerTalentsList(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cfg.authUserContextKey).(uuid.UUID)
	dbTalents, err := cfg.db.GetTalents(r.Context(), userID)
	if err != nil {
		log.Printf("Error retrieving talents: %v", err)
		http.Error(w, "Couldn't retrieve talents", http.StatusInternalServerError)
		return
	}

	talents := []TalentDTO{}
	for _, dbTalent := range dbTalents {
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
		talents = append(talents, talent)
	}
	respondWithJSON(w, http.StatusOK, talents)
}
