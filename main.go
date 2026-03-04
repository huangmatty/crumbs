package main

import (
	"log"
	"net/http"
)

const port = "8080"
const filepathRoot = "."

func main() {
	mux := http.NewServeMux()
	srvr := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.HandleFunc("/ready", handlerReady)

	log.Printf("Starting Crumbs server on port %v...", port)
	log.Fatal(srvr.ListenAndServe())
}
