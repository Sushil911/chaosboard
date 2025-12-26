package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"chaosboard/internal/api"
	"chaosboard/internal/db"
	"chaosboard/internal/metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	if err := db.Init(); err != nil {
		log.Fatalf("db init failed: %v", err)
	}
	defer db.Close()

	if err := db.LoadAll(); err != nil {
		log.Fatalf("db load failed: %v", err)
	}
	log.Printf("Loaded %d experiments from disk", len(db.GetStore()))

	mux := http.NewServeMux()
	mux.HandleFunc("/", api.RootHandler)
	mux.HandleFunc("/healthz", api.HealthHandler)
	mux.HandleFunc("POST /api/experiments", api.CreateExperiment)
	mux.HandleFunc("GET /api/experiments", api.ListExperiments)
	mux.Handle("/metrics", promhttp.Handler())

	var handler http.Handler = mux
	handler = metrics.TrackRequest(handler)

	server := &http.Server{Addr: ":8080", Handler: handler}

	go func() {
		log.Println("ChaosBoard listening on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("Shutdown signal received...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("forced shutdown: %v", err)
	}
	log.Println("Graceful shutdown complete")
}
