package main

import (
	"context"
	"log"
	"net/http"

	"github.com/huangmatty/crumbs/internal/auth"
)

func (cfg *apiConfig) middlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			log.Printf("Error getting JWT: %v", err)
			http.Error(w, "Didn't get access token", http.StatusUnauthorized)
			return
		}
		userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
		if err != nil {
			log.Printf("Error validating JWT: %v", err)
			http.Error(w, "Invalid access token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), cfg.authUserContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
