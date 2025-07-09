package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/MRyutaro/rrk/internal/history"
	"github.com/MRyutaro/rrk/internal/storage"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all command history",
	Long:  `Display all recorded command history across all sessions and directories.`,
	Run: func(cmd *cobra.Command, args []string) {
		store, err := storage.New()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing storage: %v\n", err)
			os.Exit(1)
		}

		// Get limit from flag
		limit, _ := cmd.Flags().GetInt("limit")
		
		filter := history.EntryFilter{
			Limit: limit,
		}

		entries, err := store.Load(filter)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading history: %v\n", err)
			os.Exit(1)
		}

		if len(entries) == 0 {
			fmt.Println("No history found.")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTIME\tDIRECTORY\tSESSION\tCOMMAND")
		for _, entry := range entries {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
				entry.ID,
				entry.Timestamp.Format("15:04:05"),
				shortPath(entry.CWD),
				shortSessionID(entry.SessionID),
				entry.Command)
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().IntP("limit", "n", 50, "Limit number of entries to show")
}