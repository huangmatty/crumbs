package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/huangmatty/crumbs/internal/database"
)

func (cfg *apiConfig) handlerBuyersUpdateName(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Name string `json:"name"`
	}{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Couldn't decode JSON", http.StatusBadRequest)
		return
	}
	if params.Name == "" {
		http.Error(w, "Missing name", http.StatusBadRequest)
		return
	}
	if len(params.Name) > maxNameLength {
		http.Error(w, "Name is too long", http.StatusBadRequest)
		return
	}

	buyerID, err := uuid.Parse(r.PathValue("buyerID"))
	if err != nil {
		log.Printf("Error getting buyer id: %v", err)
		http.Error(w, "Invalid buyer id", http.StatusBadRequest)
		return
	}

	dbBuyerUserID, err := cfg.db.GetUserIDForBuyer(r.Context(), buyerID)
	if err == sql.ErrNoRows {
		http.Error(w, "Buyer doesn't exist", http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Printf("Error retrieving buyer's user id: %v", err)
		http.Error(w, "Couldn't retrieve buyer's user id", http.StatusInternalServerError)
		return
	}

	userID := r.Context().Value(cfg.authUserContextKey).(uuid.UUID)
	if userID != dbBuyerUserID {
		http.Error(w, "Cannot update buyer", http.StatusForbidden)
		return
	}

	dbBuyer, err := cfg.db.UpdateBuyerName(r.Context(), database.UpdateBuyerNameParams{
		Name: params.Name,
		ID:   buyerID,
	})
	if err != nil {
		log.Printf("Error updating buyer: %v", err)
		http.Error(w, "Failed to update buyer", http.StatusInternalServerError)
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
