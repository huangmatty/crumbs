package main

import (
	"log"
	"net/http"

	"github.com/huangmatty/crumbs/internal/auth"
)

func (cfg *apiConfig) handlerBuyersList(w http.ResponseWriter, r *http.Request) {
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

	dbBuyers, err := cfg.db.GetBuyers(r.Context(), userID)
	if err != nil {
		log.Printf("Error retrieving buyers: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't retieve buyers")
		return
	}

	buyers := []BuyerDTO{}
	for _, dbBuyer := range dbBuyers {
		buyer := BuyerDTO{
			ID:        dbBuyer.ID,
			CreatedAt: dbBuyer.CreatedAt,
			UpdatedAt: dbBuyer.UpdatedAt,
			Name:      dbBuyer.Name,
			Email:     dbBuyer.Email,
			UserID:    dbBuyer.UserID,
		}
		buyers = append(buyers, buyer)
	}
	respondWithJSON(w, http.StatusOK, buyers)
}
