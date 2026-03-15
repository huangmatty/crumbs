package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/huangmatty/crumbs/internal/database"
)

const maxNameLength = 100

type TalentDTO struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerTalentsCreate(w http.ResponseWriter, r *http.Request) {
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
	email := sql.NullString{}
	if params.Email != "" {
		email.String = params.Email
		email.Valid = true
	}

	userID := r.Context().Value(cfg.authUserContextKey).(uuid.UUID)
	dbTalent, err := cfg.db.CreateTalent(r.Context(), database.CreateTalentParams{
		Name:   params.Name,
		Email:  email,
		UserID: userID,
	})
	if err != nil {
		log.Printf("Error creating talent: %v", err)
		http.Error(w, "Failed to create talent", http.StatusInternalServerError)
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
	respondWithJSON(w, http.StatusCreated, talent)
}
