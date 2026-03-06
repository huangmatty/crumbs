package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerTalentsList(w http.ResponseWriter, r *http.Request) {
	dbTalents, err := cfg.db.GetTalents(r.Context())
	if err != nil {
		log.Printf("Error retrieving talents: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve talents")
		return
	}

	talents := []TalentDTO{}
	for _, dbTalent := range dbTalents {
		talent := TalentDTO{
			ID:     dbTalent.ID,
			UserID: dbTalent.UserID,
			Name:   dbTalent.Name,
		}
		if dbTalent.Email.Valid {
			talent.Email = dbTalent.Email.String
		}
		talents = append(talents, talent)
	}
	respondWithJSON(w, http.StatusOK, talents)
}
