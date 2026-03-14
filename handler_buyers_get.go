package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/huangmatty/crumbs/internal/auth"
)

func (cfg *apiConfig) handlerBuyersGet(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting JWT: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Couldn't get access token")
		return
	}
	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error validating JWT")
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
		respondWithError(w, http.StatusNotFound, "Buyer doesn't exist")
		return
	}
	if err != nil {
		log.Printf("Error retrieving buyer: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve buyer")
		return
	}
	if userID != dbBuyer.UserID {
		respondWithError(w, http.StatusForbidden, "Cannot access buyer")
		return
	}

	buyer := BuyerDTO{
		ID:        dbBuyer.ID,
		CreatedAt: dbBuyer.CreatedAt,
		UpdatedAt: dbBuyer.UpdatedAt,
		Name:      dbBuyer.Name,
		Email:     dbBuyer.Email,
		UserID:    dbBuyer.UserID,
	}
	respondWithJSON(w, http.StatusOK, buyer)
}
