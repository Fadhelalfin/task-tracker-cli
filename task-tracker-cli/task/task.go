package task

import (
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // e.g., "pending", "in-progress", "completed"
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

const (
	StatusPending    = "pending"
	StatusInProgress = "in-progress"
	StatusCompleted  = "completed"
)