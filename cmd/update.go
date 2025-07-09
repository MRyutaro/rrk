package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"

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
	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	// Build download URL
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

	// Download the new binary
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// Create temporary file
	tmpFile := execPath + ".tmp"
	out, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer out.Close()

	// Copy downloaded content to temp file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("failed to write temp file: %v", err)
	}

	// Make executable
	if err := os.Chmod(tmpFile, 0755); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("failed to make executable: %v", err)
	}

	// Backup current binary
	backupPath := execPath + ".backup"
	if err := os.Rename(execPath, backupPath); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("failed to backup current binary: %v", err)
	}

	// Replace with new binary
	if err := os.Rename(tmpFile, execPath); err != nil {
		// Restore backup if replacement fails
		_ = os.Rename(backupPath, execPath)
		_ = os.Remove(tmpFile)
		return fmt.Errorf("failed to replace binary: %v", err)
	}

	// Remove backup
	_ = os.Remove(backupPath)

	// Verify the new binary works
	cmd := exec.Command(execPath, "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("new binary verification failed: %v", err)
	}

	fmt.Println("✅ rrk has been successfully updated!")
	fmt.Printf("Updated binary location: %s\n", execPath)

	// Show new version
	versionCmd := exec.Command(execPath, "--version")
	versionCmd.Stdout = os.Stdout
	_ = versionCmd.Run()

	return nil
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
