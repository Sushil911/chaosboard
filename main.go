package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func listExperiments(w http.ResponseWriter, r *http.Request) {

}

func createExperiments(w http.ResponseWriter, r *http.Request) {

}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to Chaosboard - v0.1 \n")
		fmt.Fprintf(w, "GET /healthz - check health of your container \n")
		fmt.Fprintf(w, "GET /api/experiments - list all chaos experiments \n")
		fmt.Fprintf(w, "POST /api/experiments - start chaos experiments \n")

	})

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok \n"))
	})

	mux.HandleFunc("GET /api/experiments", listExperiments)
	mux.HandleFunc("POST /api/experiments", createExperiments)

	server := &http.Server{Addr: ":8080", Handler: mux}

	go func() {
		log.Println("Chaosboard listening on 8080")
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server start failed %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("Shutdown signal received")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		log.Printf("Shutting down forcefully %v", err)
	} else {
		log.Println("Shutting down gracefully")
	}
}
