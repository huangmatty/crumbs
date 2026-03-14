package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/huangmatty/crumbs/internal/auth"
)

func (cfg *apiConfig) handlerBuyersDelete(w http.ResponseWriter, r *http.Request) {
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

	buyerIDStr := r.PathValue("buyerID")
	buyerID, err := uuid.Parse(buyerIDStr)
	if err != nil {
		log.Printf("Error parsing buyer id: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid buyer id")
		return
	}

	dbBuyer, err := cfg.db.GetBuyerByID(r.Context(), buyerID)
	if err == sql.ErrNoRows {
		log.Printf("Error retrieving buyer: %v", err)
		respondWithError(w, http.StatusNotFound, "Buyer doesn't exist")
		return
	}
	if err != nil {
		log.Printf("Error retrieving buyer: %v", err)
		respondWithError(w, http.StatusNotFound, "Couldn't retrieve buyer")
		return
	}
	if userID != dbBuyer.UserID {
		respondWithError(w, http.StatusForbidden, "Cannot delete buyer")
		return
	}

	_, err = cfg.db.SoftDeleteBuyer(r.Context(), buyerID)
	if err != nil {
		log.Printf("Error soft deleting buyer: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to delete buyer")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
