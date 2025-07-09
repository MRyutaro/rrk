package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Version information (set by build flags)
var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

// GitHubRelease represents a GitHub release
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
}

// GetLatestVersion fetches the latest version from GitHub
func GetLatestVersion() (string, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	
	resp, err := client.Get("https://api.github.com/repos/MRyutaro/rrk/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return "", err
	}
	
	return release.TagName, nil
}

// GetVersionInfo returns formatted version information
func GetVersionInfo() string {
	if Version == "dev" {
		// Try to get latest version from GitHub
		if latestVersion, err := GetLatestVersion(); err == nil {
			return fmt.Sprintf("rrk %s (development build)\nLatest release: %s", Version, latestVersion)
		}
		return fmt.Sprintf("rrk %s (development build)", Version)
	}
	
	return fmt.Sprintf("rrk %s", Version)
}