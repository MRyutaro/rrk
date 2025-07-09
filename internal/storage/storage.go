package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/MRyutaro/rrk/internal/history"
)

// Storage handles persistent storage of history entries
type Storage struct {
	basePath string
	mu       sync.RWMutex
	nextID   int
}

// New creates a new Storage instance
func New() (*Storage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	basePath := filepath.Join(homeDir, "rrk")
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create rrk directory: %w", err)
	}

	s := &Storage{
		basePath: basePath,
		nextID:   1,
	}

	// Load existing entries to determine next ID
	if err := s.loadNextID(); err != nil {
		return nil, err
	}

	return s, nil
}

// historyFile returns the path to the history file
func (s *Storage) historyFile() string {
	return filepath.Join(s.basePath, "history.jsonl")
}

// loadNextID loads the next available ID from existing entries
func (s *Storage) loadNextID() error {
	file, err := os.Open(s.historyFile())
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No history file yet
		}
		return fmt.Errorf("failed to open history file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var maxID int

	for scanner.Scan() {
		var entry history.Entry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			continue // Skip invalid entries
		}
		if entry.ID > maxID {
			maxID = entry.ID
		}
	}

	s.nextID = maxID + 1
	return scanner.Err()
}

// Save saves a new history entry
func (s *Storage) Save(entry *history.Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Assign ID if not set
	if entry.ID == 0 {
		entry.ID = s.nextID
		s.nextID++
	}

	// Open file in append mode
	file, err := os.OpenFile(s.historyFile(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open history file: %w", err)
	}
	defer file.Close()

	// Write entry as JSON line
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(entry); err != nil {
		return fmt.Errorf("failed to write entry: %w", err)
	}

	return nil
}

// Load loads history entries based on filter criteria
func (s *Storage) Load(filter history.EntryFilter) ([]history.Entry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	file, err := os.Open(s.historyFile())
	if err != nil {
		if os.IsNotExist(err) {
			return []history.Entry{}, nil // No history yet
		}
		return nil, fmt.Errorf("failed to open history file: %w", err)
	}
	defer file.Close()

	var entries []history.Entry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var entry history.Entry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			continue // Skip invalid entries
		}

		// Apply filters
		if filter.SessionID != nil && entry.SessionID != *filter.SessionID {
			continue
		}
		if filter.CWD != nil && entry.CWD != *filter.CWD {
			continue
		}

		entries = append(entries, entry)

		// Apply limit
		if filter.Limit > 0 && len(entries) >= filter.Limit {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read history file: %w", err)
	}

	return entries, nil
}

// GetByID retrieves a specific history entry by ID
func (s *Storage) GetByID(id int) (*history.Entry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	file, err := os.Open(s.historyFile())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("history entry not found")
		}
		return nil, fmt.Errorf("failed to open history file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var entry history.Entry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			continue
		}
		if entry.ID == id {
			return &entry, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read history file: %w", err)
	}

	return nil, fmt.Errorf("history entry not found")
}

// ListSessions returns all unique session IDs
func (s *Storage) ListSessions() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	file, err := os.Open(s.historyFile())
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to open history file: %w", err)
	}
	defer file.Close()

	sessionMap := make(map[string]bool)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var entry history.Entry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			continue
		}
		sessionMap[entry.SessionID] = true
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read history file: %w", err)
	}

	var sessions []string
	for session := range sessionMap {
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// ListDirectories returns all unique directories with history
func (s *Storage) ListDirectories() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	file, err := os.Open(s.historyFile())
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to open history file: %w", err)
	}
	defer file.Close()

	dirMap := make(map[string]bool)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var entry history.Entry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			continue
		}
		dirMap[entry.CWD] = true
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read history file: %w", err)
	}

	var dirs []string
	for dir := range dirMap {
		dirs = append(dirs, dir)
	}

	return dirs, nil
}
