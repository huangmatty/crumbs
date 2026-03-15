package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerBuyersList(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(cfg.authUserContextKey).(uuid.UUID)
	dbBuyers, err := cfg.db.GetBuyers(r.Context(), userID)
	if err != nil {
		log.Printf("Error retrieving buyers: %v", err)
		http.Error(w, "Couldn't retrieve buyers", http.StatusInternalServerError)
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
