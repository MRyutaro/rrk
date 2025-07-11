package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/MRyutaro/rrk/internal/history"
	"github.com/MRyutaro/rrk/internal/session"
	"github.com/MRyutaro/rrk/internal/storage"
	"github.com/spf13/cobra"
)

var sessionCmd = &cobra.Command{
	Use:     "session",
	Aliases: []string{"s"},
	Short:   "Manage shell sessions",
	Long:    `List, show, and manage shell session histories.`,
}

var sessionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all sessions",
	Run: func(cmd *cobra.Command, args []string) {
		store, err := storage.New()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing storage: %v\n", err)
			os.Exit(1)
		}

		sessions, err := store.ListSessions()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing sessions: %v\n", err)
			os.Exit(1)
		}

		if len(sessions) == 0 {
			fmt.Println("No sessions found.")
			return
		}

		currentSession, _ := session.GetCurrentSessionID()

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "SESSION_ID\tSTATUS")
		for _, s := range sessions {
			status := ""
			if s == currentSession {
				status = "(current)"
			}
			fmt.Fprintf(w, "%s\t%s\n", s, status)
		}
		w.Flush()
	},
}

var sessionShowCmd = &cobra.Command{
	Use:   "show [session_id]",
	Short: "Show session history",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		store, err := storage.New()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing storage: %v\n", err)
			os.Exit(1)
		}

		sessionID := ""
		if len(args) > 0 && args[0] != "current" {
			sessionID = args[0]
		} else {
			// 現在のセッションIDを取得
			sessionID, err = session.GetCurrentSessionID()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting current session: %v\n", err)
				os.Exit(1)
			}
		}

		filter := history.EntryFilter{
			SessionID: &sessionID,
		}

		entries, err := store.Load(filter)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading history: %v\n", err)
			os.Exit(1)
		}

		if len(entries) == 0 {
			fmt.Printf("No history found for session: %s\n", sessionID)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTIME\tDIRECTORY\tCOMMAND")
		for _, entry := range entries {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\n",
				entry.ID,
				entry.Timestamp.Format("15:04:05"),
				shortPath(entry.CWD),
				entry.Command)
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(sessionCmd)
	sessionCmd.AddCommand(sessionListCmd)
	sessionCmd.AddCommand(sessionShowCmd)
}
