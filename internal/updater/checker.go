package updater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

// compareVersions compares two semantic versions
// Returns: 1 if v1 > v2, -1 if v1 < v2, 0 if v1 == v2
func compareVersions(v1, v2 string) int {
	// Remove 'v' prefix if present
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")

	// Handle dev versions (always consider them older)
	if v1 == "dev" && v2 != "dev" {
		return -1
	}
	if v1 != "dev" && v2 == "dev" {
		return 1
	}
	if v1 == "dev" && v2 == "dev" {
		return 0
	}

	// Split versions into components
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	// Ensure both have at least 3 parts (major.minor.patch)
	for len(parts1) < 3 {
		parts1 = append(parts1, "0")
	}
	for len(parts2) < 3 {
		parts2 = append(parts2, "0")
	}

	// Compare each component
	for i := 0; i < 3; i++ {
		num1, err1 := strconv.Atoi(parts1[i])
		num2, err2 := strconv.Atoi(parts2[i])

		// If parsing fails, fall back to string comparison
		if err1 != nil || err2 != nil {
			if parts1[i] > parts2[i] {
				return 1
			} else if parts1[i] < parts2[i] {
				return -1
			}
			continue
		}

		if num1 > num2 {
			return 1
		} else if num1 < num2 {
			return -1
		}
	}

	return 0
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
	if cache.Latest != "" && compareVersions(cache.Latest, currentVersion) > 0 {
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
	_ = json.Unmarshal(data, &cache)
	return cache
}

func saveCache(cache VersionCache) {
	cachePath := getCacheFilePath()
	if cachePath == "" {
		return
	}

	// Create directory if it doesn't exist
	_ = os.MkdirAll(filepath.Dir(cachePath), 0755)

	data, err := json.Marshal(cache)
	if err != nil {
		return
	}

	_ = os.WriteFile(cachePath, data, 0644)
}

// ClearCache removes the version cache file
func ClearCache() {
	cachePath := getCacheFilePath()
	if cachePath != "" {
		_ = os.Remove(cachePath)
	}
}
