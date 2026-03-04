package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

const port = "8080"
const filepathRoot = "."

type apiConfig struct {
	fsHits atomic.Int32
}

func main() {
	cfg := &apiConfig{}
	mux := http.NewServeMux()
	srvr := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/ready", handlerReady)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)

	log.Printf("Starting Crumbs server on port %v...", port)
	log.Fatal(srvr.ListenAndServe())
}
