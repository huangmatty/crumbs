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
		respondWithError(w, http.StatusBadRequest, "Couldn't get refresh token")
		return
	}
	_, err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		log.Printf("Error revoking refresh token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to revoke refresh token")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
