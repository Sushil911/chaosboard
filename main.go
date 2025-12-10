package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	bolt "go.etcd.io/bbolt"
)

type Experiments struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Duration  int       `json:"duration"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	db      *bolt.DB
	store   = make(map[string]Experiments)
	storeMu sync.RWMutex
)

const (
	dbFile = "chaosboard.db"
	bucket = "experiments"
)

func initDB() error {
	var err error
	db, err = bolt.Open(dbFile, 0666, nil)
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		return err
	})
}

func saveToDB(exp Experiments) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		data, _ := json.Marshal(exp)
		return b.Put([]byte(exp.ID), data)
	})
}

func loadFromDB() error {
	return db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}
		return b.ForEach(func(k, v []byte) error {
			var exp Experiments
			if err := json.Unmarshal(v, &exp); err != nil {
				return err
			}
			storeMu.Lock()
			store[exp.ID] = exp
			storeMu.Unlock()
			return nil
		})
	})
}

func listExperiments(w http.ResponseWriter, r *http.Request) {
	storeMu.RLock()
	defer storeMu.RUnlock()

	list := make([]Experiments, 0, len(store))
	for _, e := range store {
		list = append(list, e)
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(list)
	if err != nil {
		log.Printf("Failed to encode experiments: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func createExperiments(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Type     string `json:"type"`
		Duration int    `json:"duration"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid values", http.StatusBadRequest)
		return
	}

	if req.Type == "" {
		http.Error(w, "Invalid experiment type. Please enter correct experiment", http.StatusBadRequest)
		return
	}

	if req.Duration <= 0 {
		req.Duration = 10
	}

	exp := Experiments{
		ID:        uuid.New().String(),
		Type:      req.Type,
		Duration:  req.Duration,
		Status:    "running",
		CreatedAt: time.Now(),
	}

	storeMu.Lock()
	store[exp.ID] = exp
	storeMu.Unlock()

	err = saveToDB(exp)
	if err != nil {
		log.Printf("failed to save to db: %v", err)
	}

	go runExperiments(exp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(exp)
	if err != nil {
		log.Printf("Failed to encode experiment response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func runExperiments(e Experiments) {
	log.Printf("[CHAOS START] id=%s type=%s duration=%ds", e.ID, e.Type, e.Duration)
	switch e.Type {
	case "cpu-hog":
		deadline := time.Now().Add(time.Duration(e.Duration) * time.Second)
		for time.Now().Before(deadline) {
			_ = 1<<63 - 1
		}
	}
	storeMu.Lock()
	exp, ok := store[e.ID]
	if ok {
		exp.Status = "completed"
		store[e.ID] = exp
		saveToDB(exp)
	}
	storeMu.Unlock()

	log.Printf("[CHAOS END] id=%s", e.ID)
}

func main() {

	if err := initDB(); err != nil {
		log.Fatalf("Failed to open db: %v", err)
	}

	if err := loadFromDB(); err != nil {
		log.Fatalf("Error loading from db:%v", err)
	}
	log.Printf("Loaded %d experiments from disk", len(store))

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
