package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerBuyersDelete(w http.ResponseWriter, r *http.Request) {
	buyerIDStr := r.PathValue("buyerID")
	buyerID, err := uuid.Parse(buyerIDStr)
	if err != nil {
		log.Printf("Error parsing buyer id: %v", err)
		http.Error(w, "Invalid buyer id", http.StatusBadRequest)
		return
	}

	dbBuyer, err := cfg.db.GetBuyerByID(r.Context(), buyerID)
	if err == sql.ErrNoRows {
		log.Printf("Error retrieving buyer: %v", err)
		http.Error(w, "Buyer doesn't exist", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Error retrieving buyer: %v", err)
		http.Error(w, "Couldn't retrieve buyer", http.StatusInternalServerError)
		return
	}

	userID := r.Context().Value(cfg.authUserContextKey).(uuid.UUID)
	if userID != dbBuyer.UserID {
		http.Error(w, "Cannot delete buyer", http.StatusForbidden)
		return
	}

	_, err = cfg.db.SoftDeleteBuyer(r.Context(), buyerID)
	if err != nil {
		log.Printf("Error soft deleting buyer: %v", err)
		http.Error(w, "Failed to delete buyer", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
