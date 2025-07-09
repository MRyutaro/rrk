package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup shell integration",
	Long:  `Interactive setup to configure shell integration for automatic history recording.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Setting up rrk shell integration...")
		
		// Detect shell
		shell := detectShell()
		if shell == "" {
			fmt.Println("Could not detect shell. Please specify with --shell flag.")
			os.Exit(1)
		}
		
		fmt.Printf("Detected shell: %s\n", shell)
		
		// Generate hook script
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
		
		// Write to config file
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			os.Exit(1)
		}
		
		configDir := filepath.Join(homeDir, "rrk")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating config directory: %v\n", err)
			os.Exit(1)
		}
		
		hookFile := filepath.Join(configDir, "hook.sh")
		if err := os.WriteFile(hookFile, []byte(hookScript), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing hook file: %v\n", err)
			os.Exit(1)
		}
		
		// Instructions
		fmt.Printf("\nSetup complete! To enable rrk integration, add the following to your shell config:\n\n")
		
		switch shell {
		case "bash":
			fmt.Printf("echo 'source %s' >> ~/.bashrc\n", hookFile)
			fmt.Printf("source ~/.bashrc\n")
		case "zsh":
			fmt.Printf("echo 'source %s' >> ~/.zshrc\n", hookFile)
			fmt.Printf("source ~/.zshrc\n")
		}
		
		fmt.Println("\nOr run this command to add it automatically:")
		fmt.Printf("rrk hook init %s >> ~/.%src && source ~/.%src\n", shell, shell, shell)
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
}