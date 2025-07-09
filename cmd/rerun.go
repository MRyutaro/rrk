package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/MRyutaro/rrk/internal/storage"
	"github.com/spf13/cobra"
)

var rerunCmd = &cobra.Command{
	Use:   "rerun <history_id>",
	Short: "Re-execute a command from history",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		histID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid history ID: %s\n", args[0])
			os.Exit(1)
		}

		store, err := storage.New()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing storage: %v\n", err)
			os.Exit(1)
		}

		entry, err := store.GetByID(histID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding history entry: %v\n", err)
			os.Exit(1)
		}

		// Change to the original directory and execute the command
		shellCmd := fmt.Sprintf("cd %s && %s", shellescape(entry.CWD), entry.Command)
		execCmd := exec.Command("sh", "-c", shellCmd)
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		execCmd.Stdin = os.Stdin

		if err := execCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Command failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(rerunCmd)
}

// shellescape escapes a string for safe use in shell commands
func shellescape(s string) string {
	// Simple escaping - wrap in single quotes and escape any single quotes
	return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
}
