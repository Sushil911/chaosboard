package api

import (
	"fmt"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ChaosBoard v1.0 – Clean architecture\n\n")
	fmt.Fprintf(w, "POST /api/experiments → start chaos\n")
	fmt.Fprintf(w, "GET  /api/experiments → list all\n")
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok\n")); err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}
