package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall rrk and remove shell integration",
	Long:  `Remove rrk shell integration and delete all history data.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Uninstalling rrk...")

		// Get flags
		autoConfirm, _ := cmd.Flags().GetBool("yes")

		// Get confirmation unless --yes flag is used
		if !autoConfirm {
			fmt.Print("This will remove rrk shell integration and delete all history data. Continue? [y/N]: ")
			var response string
			if _, err := fmt.Scanln(&response); err != nil {
				// Treat scan error as "no" response
				fmt.Println("Uninstall cancelled.")
				return
			}
			if response != "y" && response != "Y" && response != "yes" {
				fmt.Println("Uninstall cancelled.")
				return
			}
		}

		// Detect shell
		shell := detectShell()
		if shell == "" {
			fmt.Println("Could not detect shell. Manual cleanup may be required.")
		} else {
			// Remove shell integration
			if err := removeShellIntegration(shell); err != nil {
				fmt.Fprintf(os.Stderr, "Error removing shell integration: %v\n", err)
			} else {
				fmt.Printf("✅ Removed shell integration from ~/.%src\n", shell)
			}
		}

		// Remove hook file
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		} else {
			hookFile := filepath.Join(homeDir, ".rrk", "hook.sh")
			if err := os.Remove(hookFile); err != nil && !os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error removing hook file: %v\n", err)
			} else {
				fmt.Println("✅ Removed hook script")
			}
		}

		// Remove all rrk data
		if homeDir != "" {
			rrkDir := filepath.Join(homeDir, ".rrk")
			if err := os.RemoveAll(rrkDir); err != nil {
				fmt.Fprintf(os.Stderr, "Error removing data directory: %v\n", err)
			} else {
				fmt.Println("✅ Removed all rrk data")
			}
		}

		// Instructions for removing binary
		fmt.Println("\n📦 To complete uninstallation, remove the rrk binary:")
		fmt.Println("  sudo rm /usr/local/bin/rrk")
		fmt.Println("  # or")
		fmt.Println("  rm ~/.local/bin/rrk")

		fmt.Println("\n✨ Uninstall completed!")
		fmt.Println("Please restart your shell or run 'source ~/.zshrc' (or ~/.bashrc)")
	},
}

func removeShellIntegration(shell string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	var configFile string
	switch shell {
	case "bash":
		configFile = filepath.Join(homeDir, ".bashrc")
	case "zsh":
		configFile = filepath.Join(homeDir, ".zshrc")
	default:
		return fmt.Errorf("unsupported shell: %s", shell)
	}

	// 現在の設定を読み込み
	content, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	// rrk統合行を削除
	lines := strings.Split(string(content), "\n")
	var newLines []string
	skipNext := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// rrk関連行をスキップ
		if strings.Contains(trimmed, "rrk shell integration") ||
			strings.Contains(trimmed, "rrk hook init") ||
			(strings.HasPrefix(trimmed, "source") && strings.Contains(trimmed, "rrk/hook.sh")) {
			// コメントだけの場合は次の行もスキップ
			if i+1 < len(lines) && strings.TrimSpace(lines[i+1]) == "" {
				skipNext = true
			}
			continue
		}

		if skipNext {
			skipNext = false
			continue
		}

		newLines = append(newLines, line)
	}

	// クリーンアップした設定を書き戻し
	newContent := strings.Join(newLines, "\n")
	return os.WriteFile(configFile, []byte(newContent), 0644)
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
	uninstallCmd.Flags().BoolP("yes", "y", false, "Automatically confirm without prompting")
}
