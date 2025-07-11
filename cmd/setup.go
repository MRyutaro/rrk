package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup shell integration",
	Long:  `Automatic setup to configure shell integration for history recording.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Setting up rrk shell integration...")

		// シェルを検出
		shell := detectShell()
		if shell == "" {
			fmt.Println("Could not detect shell. Please specify with --shell flag.")
			os.Exit(1)
		}

		fmt.Printf("Detected shell: %s\n", shell)

		// --yesフラグが使用されていない場合は確認を取る
		autoConfirm, _ := cmd.Flags().GetBool("yes")
		if !autoConfirm {
			fmt.Printf("This will add rrk integration to your ~/.%src file. Continue? [y/N]: ", shell)
			var response string
			if _, err := fmt.Scanln(&response); err != nil {
				// スキャンエラーを"no"レスポンスとして扱う
				fmt.Println("Setup cancelled.")
				return
			}
			if response != "y" && response != "Y" && response != "yes" {
				fmt.Println("Setup cancelled.")
				return
			}
		}

		// フックスクリプトを取得
		hookScript := ""
		switch shell {
		case "bash":
			hookScript = bashHook()
		case "zsh":
			hookScript = zshHook()
		default:
			fmt.Printf("Unsupported shell: %s\n", shell)
			os.Exit(1)
		}

		// ホームディレクトリを取得
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		// シェル設定ファイルパス
		var shellConfigFile string
		switch shell {
		case "bash":
			shellConfigFile = filepath.Join(homeDir, ".bashrc")
		case "zsh":
			shellConfigFile = filepath.Join(homeDir, ".zshrc")
		}

		// rrkが既に設定されているかチェック
		if isAlreadyConfigured(shellConfigFile) {
			fmt.Println("rrk integration is already configured!")
			return
		}

		// フックスクリプトを設定ディレクトリに書き込み
		configDir := filepath.Join(homeDir, ".rrk")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating config directory: %v\n", err)
			os.Exit(1)
		}

		hookFile := filepath.Join(configDir, "hook.sh")
		if err := os.WriteFile(hookFile, []byte(hookScript), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing hook file: %v\n", err)
			os.Exit(1)
		}

		// シェル設定に追加
		hookLine := fmt.Sprintf("\n# rrk shell integration\nsource %s\n", hookFile)

		file, err := os.OpenFile(shellConfigFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening shell config file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		if _, err := file.WriteString(hookLine); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to shell config file: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ Setup complete!")
		fmt.Printf("rrk integration has been added to %s\n", shellConfigFile)
		fmt.Println("\nTo start using rrk, restart your shell or run:")
		fmt.Printf("source %s\n", shellConfigFile)
	},
}

func detectShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return ""
	}

	switch filepath.Base(shell) {
	case "bash":
		return "bash"
	case "zsh":
		return "zsh"
	default:
		return ""
	}
}

func init() {
	rootCmd.AddCommand(setupCmd)
	setupCmd.Flags().BoolP("yes", "y", false, "Automatically confirm setup without prompting")
}

// isAlreadyConfigured rrk統合が既に設定されているかチェック
func isAlreadyConfigured(configFile string) bool {
	content, err := os.ReadFile(configFile)
	if err != nil {
		return false
	}

	// rrk統合マーカーをチェック
	configStr := string(content)
	return strings.Contains(configStr, "rrk shell integration") ||
		strings.Contains(configStr, "rrk hook init")
}
