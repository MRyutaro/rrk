package history

import (
	"time"
)

// Entry 単一の履歴エントリを表す
type Entry struct {
	ID        int       `json:"id"`
	SessionID string    `json:"session_id"`
	CWD       string    `json:"cwd"`
	Command   string    `json:"command"`
	Timestamp time.Time `json:"timestamp"`
}

// EntryFilter 履歴エントリをフィルタリングするための条件を含む
type EntryFilter struct {
	SessionID *string
	CWD       *string
	Limit     int
}
