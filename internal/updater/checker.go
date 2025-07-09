package updater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
}

const (
	checkInterval = 24 * time.Hour
	cacheFile     = ".rrk_version_cache"
)

type VersionCache struct {
	LastCheck time.Time `json:"last_check"`
	Latest    string    `json:"latest"`
}

// CheckForUpdate checks if a newer version is available and returns update message if needed
func CheckForUpdate(currentVersion string) string {
	// Skip check if version is dev or contains commit hash
	if strings.Contains(currentVersion, "dev") || strings.Contains(currentVersion, "-") {
		return ""
	}

	cache := loadCache()
	
	// Check if we need to fetch latest version
	if time.Since(cache.LastCheck) > checkInterval {
		latest := fetchLatestVersion()
		if latest != "" {
			cache.Latest = latest
			cache.LastCheck = time.Now()
			saveCache(cache)
		}
	}

	// Compare versions and return message if update available
	if cache.Latest != "" && cache.Latest != currentVersion && !strings.HasPrefix(cache.Latest, currentVersion) {
		return fmt.Sprintf("🚀 A new version of rrk is available: %s (current: %s)\n   Run 'rrk update' to upgrade.", cache.Latest, currentVersion)
	}

	return ""
}

func fetchLatestVersion() string {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/MRyutaro/rrk/releases/latest")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return ""
	}

	return strings.TrimPrefix(release.TagName, "v")
}

func getCacheFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, ".rrk", cacheFile)
}

func loadCache() VersionCache {
	cachePath := getCacheFilePath()
	if cachePath == "" {
		return VersionCache{}
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		return VersionCache{}
	}

	var cache VersionCache
	json.Unmarshal(data, &cache)
	return cache
}

func saveCache(cache VersionCache) {
	cachePath := getCacheFilePath()
	if cachePath == "" {
		return
	}

	// Create directory if it doesn't exist
	os.MkdirAll(filepath.Dir(cachePath), 0755)

	data, err := json.Marshal(cache)
	if err != nil {
		return
	}

	os.WriteFile(cachePath, data, 0644)
}