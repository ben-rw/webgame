package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type config struct {
	Port         string
	FilepathRoot string
	RootURL      string
	mu           *sync.RWMutex
	activeRooms  map[string]*Room
}

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}

	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	type errorVals struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorVals{Error: msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
