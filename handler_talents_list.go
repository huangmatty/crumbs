package main

import (
	"log"
	"net/http"

	"github.com/huangmatty/crumbs/internal/auth"
)

func (cfg *apiConfig) handlerTalentsList(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting access token: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Couldn't get access token")
		return
	}
	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error validating access token: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Invalid accesss token")
		return
	}

	dbTalents, err := cfg.db.GetTalents(r.Context(), userID)
	if err != nil {
		log.Printf("Error retrieving talents: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve talents")
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
