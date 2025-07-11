package cmd

import (
	"fmt"
	"os"

	"github.com/MRyutaro/rrk/internal/updater"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rrk",
	Short: "A shell history management tool",
	Long: `rrk (rireki) is a Go-based CLI tool that manages shell history
by session and directory, making past commands easily reusable.`,
	Version: GetVersionInfo(),
}

func Execute() {
	// コマンド実行前にアップデートをチェック
	if updateMsg := updater.CheckForUpdate(Version); updateMsg != "" {
		fmt.Fprintln(os.Stderr, updateMsg)
		fmt.Fprintln(os.Stderr)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
