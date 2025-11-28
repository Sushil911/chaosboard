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

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to Chaosboard - v0.1 \n")
		fmt.Fprintf(w, "GET /healthz - check health of your container \n")
		fmt.Fprintf(w, "GET /chaos - check the resiliency of your deployed containers \n")
	})

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok \n"))
	})

	mux.HandleFunc("/chaos", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Coming soon \n")
	})

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
