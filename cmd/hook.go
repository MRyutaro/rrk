package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/MRyutaro/rrk/internal/history"
	"github.com/MRyutaro/rrk/internal/session"
	"github.com/MRyutaro/rrk/internal/storage"
	"github.com/spf13/cobra"
)

var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "Shell integration hooks",
	Long:  `Commands for shell integration to record history automatically.`,
}

var hookRecordCmd = &cobra.Command{
	Use:   "record <command>",
	Short: "Record a command to history",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		store, err := storage.New()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing storage: %v\n", err)
			os.Exit(1)
		}

		sessionID, err := session.GetCurrentSessionID()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting session ID: %v\n", err)
			os.Exit(1)
		}

		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
			os.Exit(1)
		}

		// Join all arguments to form the complete command
		command := ""
		for i, arg := range args {
			if i > 0 {
				command += " "
			}
			command += arg
		}

		// Skip recording cd commands as per requirements
		if len(command) >= 2 && command[:2] == "cd" {
			return
		}

		entry := &history.Entry{
			SessionID: sessionID,
			CWD:       cwd,
			Command:   command,
			Timestamp: time.Now(),
		}

		if err := store.Save(entry); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving history: %v\n", err)
			os.Exit(1)
		}
	},
}

var hookInitCmd = &cobra.Command{
	Use:   "init <shell>",
	Short: "Initialize shell integration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shell := args[0]
		switch shell {
		case "bash":
			fmt.Print(bashHook())
		case "zsh":
			fmt.Print(zshHook())
		default:
			fmt.Fprintf(os.Stderr, "Unsupported shell: %s\n", shell)
			os.Exit(1)
		}
	},
}

func bashHook() string {
	return `# rrk shell integration for bash
_rrk_hook() {
    local exit_code=$?
    local command=$(history 1 | sed 's/^[ ]*[0-9]*[ ]*//')
    if [ -n "$command" ]; then
        rrk hook record "$command" 2>/dev/null || true
    fi
    return $exit_code
}

# Set up the hook
if [ -z "$RRK_SESSION_ID" ]; then
    export RRK_SESSION_ID=$(rrk hook session-init 2>/dev/null || echo "unknown")
fi

# Install the hook
if [[ "$PROMPT_COMMAND" != *"_rrk_hook"* ]]; then
    PROMPT_COMMAND="${PROMPT_COMMAND:+$PROMPT_COMMAND; }_rrk_hook"
fi
`
}

func zshHook() string {
	return `# rrk shell integration for zsh
_rrk_hook() {
    local exit_code=$?
    local command=$(fc -ln -1)
    if [ -n "$command" ]; then
        rrk hook record "$command" 2>/dev/null || true
    fi
    return $exit_code
}

# Set up the hook
if [ -z "$RRK_SESSION_ID" ]; then
    export RRK_SESSION_ID=$(rrk hook session-init 2>/dev/null || echo "unknown")
fi

# Install the hook
autoload -U add-zsh-hook
add-zsh-hook precmd _rrk_hook
`
}

var hookSessionInitCmd = &cobra.Command{
	Use:   "session-init",
	Short: "Initialize a new session",
	Run: func(cmd *cobra.Command, args []string) {
		sessionID, err := session.InitializeSession()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing session: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(sessionID)
	},
}

func init() {
	rootCmd.AddCommand(hookCmd)
	hookCmd.AddCommand(hookRecordCmd)
	hookCmd.AddCommand(hookInitCmd)
	hookCmd.AddCommand(hookSessionInitCmd)
}
