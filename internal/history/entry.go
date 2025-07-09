package history

import (
	"time"
)

// Entry represents a single history entry
type Entry struct {
	ID        int       `json:"id"`
	SessionID string    `json:"session_id"`
	CWD       string    `json:"cwd"`
	Command   string    `json:"command"`
	Timestamp time.Time `json:"timestamp"`
}

// EntryFilter contains criteria for filtering history entries
type EntryFilter struct {
	SessionID *string
	CWD       *string
	Limit     int
}
