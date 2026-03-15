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

type contextKey string

type apiConfig struct {
	db                 *database.Queries
	platform           string
	jwtSecret          string
	authUserContextKey contextKey
	fsHits             atomic.Int32
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
		db:                 database.New(db),
		platform:           platform,
		jwtSecret:          jwtSecret,
		authUserContextKey: "user",
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

	mux.Handle("PUT /api/users", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerUsersUpdate)))
	mux.HandleFunc("POST /api/users", cfg.handlerUsersCreate)

	mux.Handle("GET /api/talents", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerTalentsList)))
	mux.Handle("GET /api/talents/{talentID}", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerTalentsGet)))
	mux.Handle("PUT /api/talents/{talentID}/name", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerTalentsUpdateName)))
	mux.Handle("PUT /api/talents/{talentID}/restore", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerTalentsRestore)))
	mux.Handle("POST /api/talents", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerTalentsCreate)))
	mux.Handle("DELETE /api/talents/{talentID}", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerTalentsDelete)))

	mux.Handle("GET /api/buyers", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerBuyersList)))
	mux.Handle("GET /api/buyers/{buyerID}", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerBuyersGet)))
	mux.Handle("PUT /api/buyers/{buyerID}/name", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerBuyersUpdateName)))
	mux.Handle("PUT /api/buyers/{buyerID}/restore", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerBuyersRestore)))
	mux.Handle("POST /api/buyers", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerBuyersCreate)))
	mux.Handle("DELETE /api/buyers/{buyerID}", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerBuyersDelete)))

	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)

	log.Printf("Starting Crumbs server on port %v...", port)
	log.Fatal(server.ListenAndServe())
}
