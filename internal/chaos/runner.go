package chaos

import (
	"log"
	"os"
	"time"

	"chaosboard/internal/db"
	"chaosboard/internal/models"
)

func Run(e models.Experiment) {
	log.Printf("[CHAOS START] id=%s type=%s duration=%ds", e.ID, e.Type, e.Duration)

	switch e.Type {
	case "cpu-hog":
		cpuHog(e.Duration)

	case "memory-hog":
		memoryHog(e.Duration)

	case "disk-fill":
		diskFill(e.Duration)

	default:
		log.Printf("unknown chaos type %q", e.Type)
	}

	// Mark completed
	e.Status = "completed"
	db.Update(e)
	log.Printf("[CHAOS END] id=%s", e.ID)
}

func cpuHog(seconds int) {
	end := time.Now().Add(time.Duration(seconds) * time.Second)
	for time.Now().Before(end) {
		_ = 1<<63 - 1
	}
}

func memoryHog(seconds int) {
	var junk [][]byte
	end := time.Now().Add(time.Duration(seconds) * time.Second)
	for time.Now().Before(end) {
		junk = append(junk, make([]byte, 10*1024*1024)) // 10 MB
		time.Sleep(100 * time.Millisecond)
	}
}

func diskFill(seconds int) {
	file, _ := os.Create("chaos_junk.tmp")
	defer os.Remove("chaos_junk.tmp")
	end := time.Now().Add(time.Duration(seconds) * time.Second)
	for time.Now().Before(end) {
		file.Write(make([]byte, 1024*1024)) // 1 MB
		file.Sync()
		time.Sleep(100 * time.Millisecond)
	}
	file.Close()
}
