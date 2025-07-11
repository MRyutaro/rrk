package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/MRyutaro/rrk/internal/updater"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update rrk to the latest version",
	Long:  `Download and install the latest version of rrk from GitHub releases.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := updateRrk(); err != nil {
			fmt.Fprintf(os.Stderr, "Update failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func updateRrk() error {
	// 現在の実行ファイルパスを取得
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	// ダウンロードURLを構築
	arch := runtime.GOARCH
	osName := runtime.GOOS
	if osName == "darwin" && arch == "arm64" {
		arch = "arm64"
	} else if arch == "amd64" {
		arch = "amd64"
	} else {
		return fmt.Errorf("unsupported architecture: %s", arch)
	}

	var binaryName string
	if osName == "windows" {
		binaryName = fmt.Sprintf("rrk-%s-%s.exe", osName, arch)
	} else {
		binaryName = fmt.Sprintf("rrk-%s-%s", osName, arch)
	}

	downloadURL := fmt.Sprintf("https://github.com/MRyutaro/rrk/releases/latest/download/%s", binaryName)

	fmt.Printf("Downloading latest version from %s...\n", downloadURL)

	// 新しいバイナリをダウンロード
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// 一時ファイルを作成
	tmpFile := execPath + ".tmp"
	out, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer out.Close()

	// ダウンロードしたコンテンツを一時ファイルにコピー
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("failed to write temp file: %v", err)
	}

	// 実行可能にする
	if err := os.Chmod(tmpFile, 0755); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("failed to make executable: %v", err)
	}

	// 現在のバイナリをバックアップ
	backupPath := execPath + ".backup"
	if err := os.Rename(execPath, backupPath); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("failed to backup current binary: %v", err)
	}

	// 新しいバイナリで置き換え
	if err := os.Rename(tmpFile, execPath); err != nil {
		// 置き換えに失敗した場合はバックアップを復元
		_ = os.Rename(backupPath, execPath)
		_ = os.Remove(tmpFile)
		return fmt.Errorf("failed to replace binary: %v", err)
	}

	// バックアップを削除
	_ = os.Remove(backupPath)

	// 新しいバイナリが動作するか検証
	cmd := exec.Command(execPath, "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("new binary verification failed: %v", err)
	}

	fmt.Println("✅ rrk has been successfully updated!")
	fmt.Printf("Updated binary location: %s\n", execPath)

	// 新バージョンの更新メッセージ表示を回避するためバージョンキャッシュをクリア
	updater.ClearCache()

	// 新バージョンを表示
	versionCmd := exec.Command(execPath, "--version")
	versionCmd.Stdout = os.Stdout
	_ = versionCmd.Run()

	return nil
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
