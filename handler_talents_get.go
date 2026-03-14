package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/huangmatty/crumbs/internal/auth"
)

func (cfg *apiConfig) handlerTalentsGet(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting JWT: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Couldn't get access token")
		return
	}
	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error validating JWT: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Invalid access token")
		return
	}

	talentIDStr := r.PathValue("talentID")
	talentID, err := uuid.Parse(talentIDStr)
	if err != nil {
		log.Printf("Error parsing talent id: %v", err)
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
	if userID != dbTalent.UserID {
		respondWithError(w, http.StatusForbidden, "Cannot access talent")
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
