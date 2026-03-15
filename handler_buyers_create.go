package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/huangmatty/crumbs/internal/database"
)

type BuyerDTO struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerBuyersCreate(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Name  string `json:"name"`
		Email string `json:"email"`
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
	if params.Email == "" {
		http.Error(w, "Missing email address", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(cfg.authUserContextKey).(uuid.UUID)
	dbBuyer, err := cfg.db.CreateBuyer(r.Context(), database.CreateBuyerParams{
		Name:   params.Name,
		Email:  params.Email,
		UserID: userID,
	})
	if err != nil {
		log.Printf("Error creating buyer: %v", err)
		http.Error(w, "Failed to create buyer", http.StatusInternalServerError)
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
	respondWithJSON(w, http.StatusCreated, buyer)
}
