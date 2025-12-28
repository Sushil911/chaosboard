package chaos

import (
	"log"
	"os"
	"time"

	"chaosboard/internal/db"
	"chaosboard/internal/metrics"
	"chaosboard/internal/models"
)

func Run(e models.Experiment) {
	log.Printf("[CHAOS START] id=%s type=%s duration=%ds", e.ID, e.Type, e.Duration)

	metrics.ExperimentsActive.Inc()
	defer metrics.ExperimentsActive.Dec()

	var runErr error
	switch e.Type {
	case "cpu-hog":
		runErr = cpuHog(e.Duration)

	case "memory-hog":
		runErr = memoryHog(e.Duration)

	case "disk-fill":
		runErr = diskFill(e.Duration)

	default:
		log.Printf("unknown chaos type %q", e.Type)
	}

	if runErr != nil {
		e.Status = "failed"
		metrics.ExperimentsCompleted.WithLabelValues(e.Type, e.Status).Inc()
	} else {
		e.Status = "completed"
		metrics.ExperimentsCompleted.WithLabelValues(e.Type, e.Status).Inc()
	}
	db.Update(e)
	log.Printf("[CHAOS END] id=%s status=%s", e.ID, e.Status)
}

func cpuHog(seconds int) error {
	end := time.Now().Add(time.Duration(seconds) * time.Second)
	for time.Now().Before(end) {
		_ = 1<<63 - 1
	}
	return nil
}

func memoryHog(seconds int) error {
	var junk [][]byte
	end := time.Now().Add(time.Duration(seconds) * time.Second)
	for time.Now().Before(end) {
		junk = append(junk, make([]byte, 10*1024*1024)) // 10 MB
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

func diskFill(seconds int) error {
	file, _ := os.Create("chaos_junk.tmp")
	defer os.Remove("chaos_junk.tmp")
	end := time.Now().Add(time.Duration(seconds) * time.Second)
	for time.Now().Before(end) {
		file.Write(make([]byte, 1024*1024)) // 1 MB
		file.Sync()
		time.Sleep(100 * time.Millisecond)
	}
	file.Close()
	return nil
}
