package cmd

import "os"

// shortPath returns a shortened version of the path for display
func shortPath(path string) string {
	home, _ := os.UserHomeDir()
	if home != "" && len(path) > len(home) && path[:len(home)] == home {
		return "~" + path[len(home):]
	}
	if len(path) > 30 {
		return "..." + path[len(path)-27:]
	}
	return path
}

// shortSessionID returns a shortened version of the session ID for display
func shortSessionID(sessionID string) string {
	if len(sessionID) > 20 {
		return sessionID[:8] + "..." + sessionID[len(sessionID)-8:]
	}
	return sessionID
}