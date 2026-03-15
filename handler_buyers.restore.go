package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerBuyersRestore(w http.ResponseWriter, r *http.Request) {
	buyerIDStr := r.PathValue("buyerID")
	buyerID, err := uuid.Parse(buyerIDStr)
	if err != nil {
		log.Printf("Error parsing buyer id: %v", err)
		http.Error(w, "Invalid buyer id", http.StatusBadRequest)
		return
	}

	dbBuyer, err := cfg.db.GetBuyerByID(r.Context(), buyerID)
	if err == sql.ErrNoRows {
		http.Error(w, "Buyer doesn't exist", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error retrieving buyer: %v", err)
		http.Error(w, "Couldn't retrieve buyer", http.StatusInternalServerError)
		return
	}

	userID := r.Context().Value(cfg.authUserContextKey)
	if userID != dbBuyer.UserID {
		http.Error(w, "Cannot restore buyer", http.StatusForbidden)
		return
	}

	dbBuyer, err = cfg.db.RestoreBuyer(r.Context(), buyerID)
	if err != nil {
		log.Printf("Error restoring buyer: %v", err)
		http.Error(w, "Failed to restore buyer", http.StatusInternalServerError)
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
