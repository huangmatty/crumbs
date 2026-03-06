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
	db     *database.Queries
	fsHits atomic.Int32
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Couldn't establish connection to database: %v", err)
	}

	cfg := &apiConfig{
		db: database.New(db),
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: middlewareLog(mux),
	}

	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/ready", handlerReady)
	mux.HandleFunc("GET /api/talents", cfg.handlerTalentsList)
	mux.HandleFunc("GET /api/talents/{talentID}", cfg.handlerTalentsGet)
	mux.HandleFunc("POST /api/users", cfg.handlerUsersCreate)
	mux.HandleFunc("POST /api/talents", cfg.handlerTalentsCreate)

	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)

	log.Printf("Starting Crumbs server on port %v...", port)
	log.Fatal(server.ListenAndServe())
}
