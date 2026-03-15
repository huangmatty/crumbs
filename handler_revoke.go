package main

import (
	"log"
	"net/http"

	"github.com/huangmatty/crumbs/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting refresh token: %v", err)
		http.Error(w, "Couldn't get refresh token", http.StatusBadRequest)
		return
	}
	_, err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		log.Printf("Error revoking refresh token: %v", err)
		http.Error(w, "Failed to revoke refresh token", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
