package api

import (
	"encoding/json"
	"net/http"

	"chaosboard/internal/chaos"
	"chaosboard/internal/db"
)

func CreateExperiment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Type     string `json:"type"`
		Duration int    `json:"duration"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	exp := db.Create(req.Type, req.Duration)

	go chaos.Run(exp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(exp)
}

func ListExperiments(w http.ResponseWriter, r *http.Request) {
	list := db.GetAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}
