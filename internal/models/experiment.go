package models
import "time"

type Experiment struct {
    ID        string    `json:"id"`
    Type      string    `json:"type"`
    Duration  int       `json:"duration"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
}
