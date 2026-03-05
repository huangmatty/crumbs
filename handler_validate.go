package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const maxNameLength = 100

func handlerValidateName(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Name string `json:"name"`
	}{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding name: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode name")
		return
	}
	if len(params.Name) > maxNameLength {
		respondWithError(w, http.StatusBadRequest, "Name is too long")
		return
	}
	respondWithJSON(w, http.StatusOK, params)
}
