package session

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetCurrentSessionID 現在のシェルセッションIDを返す
func GetCurrentSessionID() (string, error) {
	// まず環境変数から取得を試行
	if sessionID := os.Getenv("RRK_SESSION_ID"); sessionID != "" {
		return sessionID, nil
	}

	// シェルPID + TTYにフォールバック
	pid := os.Getpid()
	tty := os.Getenv("TTY")
	if tty == "" {
		// ttyコマンドからTTYを取得しようとする
		tty = "unknown"
	}

	// ファイルシステム安全にするためttyパスをクリーンアップ
	tty = strings.ReplaceAll(tty, "/", "-")

	return fmt.Sprintf("%d_%s", pid, tty), nil
}

// InitializeSession 新しいセッションIDを作成して保存
func InitializeSession() (string, error) {
	// 一意のセッションIDを生成
	hostname, _ := os.Hostname()
	pid := os.Getpid()
	timestamp := fmt.Sprintf("%d", os.Getpid())

	sessionID := fmt.Sprintf("%s_%d_%s", hostname, pid, timestamp)

	// 子プロセス用に環境変数に設定
	os.Setenv("RRK_SESSION_ID", sessionID)

	// 永続化のためセッションファイルにも書き込み
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	sessionFile := filepath.Join(homeDir, ".rrk", "current_session")
	if err := os.WriteFile(sessionFile, []byte(sessionID), 0644); err != nil {
		return "", err
	}

	return sessionID, nil
}
