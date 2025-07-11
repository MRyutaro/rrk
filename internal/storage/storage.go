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

// Storage 履歴エントリの永続化ストレージを管理
type Storage struct {
	basePath string
	mu       sync.RWMutex
	nextID   int
}

// New 新しいStorageインスタンスを作成
func New() (*Storage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	basePath := filepath.Join(homeDir, ".rrk")
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create rrk directory: %w", err)
	}

	s := &Storage{
		basePath: basePath,
		nextID:   1,
	}

	// 既存エントリを読み込んで次のIDを決定
	if err := s.loadNextID(); err != nil {
		return nil, err
	}

	return s, nil
}

// historyFile 履歴ファイルのパスを返す
func (s *Storage) historyFile() string {
	return filepath.Join(s.basePath, "history.jsonl")
}

// loadNextID 既存エントリから次の使用可能IDを読み込み
func (s *Storage) loadNextID() error {
	file, err := os.Open(s.historyFile())
	if err != nil {
		if os.IsNotExist(err) {
			return nil // まだ履歴ファイルがない
		}
		return fmt.Errorf("failed to open history file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var maxID int

	for scanner.Scan() {
		var entry history.Entry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			continue // 無効なエントリをスキップ
		}
		if entry.ID > maxID {
			maxID = entry.ID
		}
	}

	s.nextID = maxID + 1
	return scanner.Err()
}

// Save 新しい履歴エントリを保存
func (s *Storage) Save(entry *history.Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// IDが設定されていない場合は割り当て
	if entry.ID == 0 {
		entry.ID = s.nextID
		s.nextID++
	}

	// ファイルを追加モードで開く
	file, err := os.OpenFile(s.historyFile(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open history file: %w", err)
	}
	defer file.Close()

	// エントリをJSON行として書き込み
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(entry); err != nil {
		return fmt.Errorf("failed to write entry: %w", err)
	}

	return nil
}

// Load フィルタ条件に基づいて履歴エントリを読み込み
func (s *Storage) Load(filter history.EntryFilter) ([]history.Entry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	file, err := os.Open(s.historyFile())
	if err != nil {
		if os.IsNotExist(err) {
			return []history.Entry{}, nil // まだ履歴がない
		}
		return nil, fmt.Errorf("failed to open history file: %w", err)
	}
	defer file.Close()

	var entries []history.Entry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var entry history.Entry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			continue // 無効なエントリをスキップ
		}

		// フィルタを適用
		if filter.SessionID != nil && entry.SessionID != *filter.SessionID {
			continue
		}
		if filter.CWD != nil && entry.CWD != *filter.CWD {
			continue
		}

		entries = append(entries, entry)

		// 制限を適用
		if filter.Limit > 0 && len(entries) >= filter.Limit {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read history file: %w", err)
	}

	return entries, nil
}

// GetByID IDにより特定の履歴エントリを取得
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

// ListSessions 全ての一意なセッションIDを返す
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

// ListDirectories 履歴を持つ全ての一意なディレクトリを返す
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
