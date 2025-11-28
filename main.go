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
		fmt.Fprintf(w, "Chaosboard v0.1 - Chaos Engineering Platform \n")
		fmt.Fprintf(w, "GET /healthz -> 200 OK \n")
		fmt.Fprintf(w, "GET /chaos/start -> trigger chaos(coming soon) \n")
	})

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	server := &http.Server{Addr: ":8080", Handler: mux}
	go func() {
		log.Println("Chaosboard listening on 8080")
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed %v", err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	log.Println("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}
