package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"text/tabwriter"

	"github.com/MRyutaro/rrk/internal/history"
	"github.com/MRyutaro/rrk/internal/storage"
	"github.com/spf13/cobra"
)

var dirCmd = &cobra.Command{
	Use:     "dir",
	Aliases: []string{"d"},
	Short:   "Manage directory-based history",
	Long:    `List and show shell history for specific directories.`,
}

var dirShowCmd = &cobra.Command{
	Use:   "show [directory|dir-id]",
	Short: "Show history for a directory (by path or ID from 'rrk d list')",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		store, err := storage.New()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing storage: %v\n", err)
			os.Exit(1)
		}

		dir := ""
		if len(args) > 0 {
			inputArg := args[0]

			// 引数が数値のディレクトリ ID かチェック
			if dirID, parseErr := strconv.Atoi(inputArg); parseErr == nil {
				// 数値IDの場合、ディレクトリパスに解決
				directories, err := store.ListDirectories()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error listing directories: %v\n", err)
					os.Exit(1)
				}

				if dirID < 0 || dirID >= len(directories) {
					fmt.Fprintf(os.Stderr, "Invalid directory ID: %d. Use 'rrk d list' to see available IDs.\n", dirID)
					os.Exit(1)
				}

				dir = directories[dirID]
			} else {
				// ディレクトリパスの場合、".."のような相対パスを処理
				if !filepath.IsAbs(inputArg) {
					cwd, err := os.Getwd()
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
						os.Exit(1)
					}
					dir = filepath.Clean(filepath.Join(cwd, inputArg))
				} else {
					dir = inputArg
				}
			}
		} else {
			dir, err = os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
				os.Exit(1)
			}
		}

		filter := history.EntryFilter{
			CWD: &dir,
		}

		entries, err := store.Load(filter)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading history: %v\n", err)
			os.Exit(1)
		}

		if len(entries) == 0 {
			fmt.Printf("No history found for directory: %s\n", dir)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTIME\tSESSION\tCOMMAND")
		for _, entry := range entries {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\n",
				entry.ID,
				entry.Timestamp.Format("15:04:05"),
				shortSessionID(entry.SessionID),
				entry.Command)
		}
		w.Flush()
	},
}

var dirListCmd = &cobra.Command{
	Use:   "list",
	Short: "List directories with history",
	Run: func(cmd *cobra.Command, args []string) {
		store, err := storage.New()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing storage: %v\n", err)
			os.Exit(1)
		}

		directories, err := store.ListDirectories()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing directories: %v\n", err)
			os.Exit(1)
		}

		if len(directories) == 0 {
			fmt.Println("No directories with history found.")
			return
		}

		currentDir, _ := os.Getwd()

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tDIRECTORY\tSTATUS")
		for i, dir := range directories {
			status := ""
			if dir == currentDir {
				status = "(current)"
			}
			fmt.Fprintf(w, "%d\t%s\t%s\n", i, shortPath(dir), status)
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(dirCmd)
	dirCmd.AddCommand(dirShowCmd)
	dirCmd.AddCommand(dirListCmd)
}
