package db

import (
	"encoding/json"
	"sync"
	"time"

	"chaosboard/internal/metrics"
	"chaosboard/internal/models"

	"github.com/google/uuid"
	bolt "go.etcd.io/bbolt"
)

var (
	DB      *bolt.DB
	store   = make(map[string]models.Experiment)
	storeMu sync.RWMutex
)

const (
	dbFile = "chaosboard.db"
	bucket = "experiments"
)

func Init() error {
	var err error
	DB, err = bolt.Open(dbFile, 0666, nil)
	if err != nil {
		return err
	}

	return DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		return err
	})
}

func Close() {
	DB.Close()
}

func Save(exp models.Experiment) error {
	return DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		data, _ := json.Marshal(exp)
		return b.Put([]byte(exp.ID), data)
	})
}

func LoadAll() error {
	return DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}
		return b.ForEach(func(k, v []byte) error {
			var exp models.Experiment
			if err := json.Unmarshal(v, &exp); err != nil {
				return err
			}
			storeMu.Lock()
			store[string(k)] = exp
			storeMu.Unlock()
			return nil
		})
	})
}

func GetAll() []models.Experiment {
	storeMu.RLock()
	defer storeMu.RUnlock()

	list := make([]models.Experiment, 0, len(store))
	for _, e := range store {
		list = append(list, e)
	}
	return list
}

func GetStore() map[string]models.Experiment {
	storeMu.RLock()
	defer storeMu.RUnlock()
	copy := make(map[string]models.Experiment)
	for k, v := range store {
		copy[k] = v
	}
	return copy
}

func Update(exp models.Experiment) {
	storeMu.Lock()
	store[exp.ID] = exp
	storeMu.Unlock()
	Save(exp)
	metrics.ExperimentsActive.Dec()
}

func Create(reqType string, duration int) models.Experiment {
	if duration <= 0 {
		duration = 10
	}

	exp := models.Experiment{
		ID:        uuid.New().String(),
		Type:      reqType,
		Duration:  duration,
		Status:    "running",
		CreatedAt: time.Now(),
	}

	storeMu.Lock()
	store[exp.ID] = exp
	storeMu.Unlock()
	Save(exp)

	metrics.ExperimentsTotal.WithLabelValues(reqType).Inc()

	return exp
}
