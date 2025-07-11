package cmd

import (
	"fmt"
	"os"

	"github.com/MRyutaro/rrk/internal/history"
	"github.com/MRyutaro/rrk/internal/storage"
	"github.com/MRyutaro/rrk/internal/tree"
	"github.com/MRyutaro/rrk/internal/updater"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rrk [path]",
	Short: "A shell history visualization tool",
	Long: `rrk (rireki) is a Go-based CLI tool that displays shell history
in directory tree format, making it easy to see which commands
were executed in each directory.`,
	Version: GetVersionInfo(),
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// フラグ値を取得
		maxCommands, _ := cmd.Flags().GetInt("number")
		
		// ストレージを初期化
		store, err := storage.New()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing storage: %v\n", err)
			os.Exit(1)
		}

		// 全履歴を読み込み
		filter := history.EntryFilter{}
		entries, err := store.Load(filter)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading history: %v\n", err)
			os.Exit(1)
		}

		if len(entries) == 0 {
			fmt.Println("No command history found.")
			fmt.Println("Run some commands to see them here, or run 'rrk setup' to enable history tracking.")
			return
		}

		// ツリーを構築
		builder := tree.NewTreeBuilder()
		root := builder.BuildTree(entries, maxCommands)

		// 指定されたパスがあるかチェック
		var targetPath string
		if len(args) > 0 {
			targetPath = args[0]
		}

		// ツリーを表示
		tree.PrintTree(root, targetPath, maxCommands)
	},
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
	rootCmd.Flags().IntP("number", "n", 0, "Maximum number of commands to show per directory (0 = show all)")
}
