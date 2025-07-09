package session

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetCurrentSessionID returns the current shell session ID
func GetCurrentSessionID() (string, error) {
	// Try to get from environment variable first
	if sessionID := os.Getenv("RRK_SESSION_ID"); sessionID != "" {
		return sessionID, nil
	}

	// Fallback to shell PID + TTY
	pid := os.Getpid()
	tty := os.Getenv("TTY")
	if tty == "" {
		// Try to get TTY from tty command
		tty = "unknown"
	}
	
	// Clean tty path to make it filesystem-safe
	tty = strings.ReplaceAll(tty, "/", "-")
	
	return fmt.Sprintf("%d_%s", pid, tty), nil
}

// InitializeSession creates a new session ID and stores it
func InitializeSession() (string, error) {
	// Generate a unique session ID
	hostname, _ := os.Hostname()
	pid := os.Getpid()
	timestamp := fmt.Sprintf("%d", os.Getpid())
	
	sessionID := fmt.Sprintf("%s_%d_%s", hostname, pid, timestamp)
	
	// Set it in environment for child processes
	os.Setenv("RRK_SESSION_ID", sessionID)
	
	// Also write to a session file for persistence
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	
	sessionFile := filepath.Join(homeDir, "rrk", "current_session")
	if err := os.WriteFile(sessionFile, []byte(sessionID), 0644); err != nil {
		return "", err
	}
	
	return sessionID, nil
}