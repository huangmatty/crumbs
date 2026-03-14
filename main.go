package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/huangmatty/crumbs/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const port = "8080"
const filepathRoot = "."

type apiConfig struct {
	db        *database.Queries
	platform  string
	jwtSecret string
	fsHits    atomic.Int32
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Couldn't establish connection to database: %v", err)
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM environment variable is not set")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	cfg := &apiConfig{
		db:        database.New(db),
		platform:  platform,
		jwtSecret: jwtSecret,
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: middlewareLog(mux),
	}

	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/ready", handlerReady)
	mux.HandleFunc("POST /api/login", cfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", cfg.handlerRevoke)

	mux.HandleFunc("PUT /api/users", cfg.handlerUsersUpdate)
	mux.HandleFunc("POST /api/users", cfg.handlerUsersCreate)

	mux.HandleFunc("GET /api/talents", cfg.handlerTalentsList)
	mux.HandleFunc("GET /api/talents/{talentID}", cfg.handlerTalentsGet)
	mux.HandleFunc("POST /api/talents", cfg.handlerTalentsCreate)
	mux.HandleFunc("DELETE /api/talent/{talentID}", cfg.handlerTalentsDelete)

	mux.HandleFunc("GET /api/buyers/{buyerID}", cfg.handlerBuyersGet)
	mux.HandleFunc("POST /api/buyers", cfg.handlerBuyersCreate)

	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)

	log.Printf("Starting Crumbs server on port %v...", port)
	log.Fatal(server.ListenAndServe())
}
