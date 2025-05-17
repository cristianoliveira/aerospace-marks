/*
Copyright © 2025 Cristian Oliveira me@cristianoliveira.dev
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/cristianoliveira/aerospace-marks/internal/aerospace"
	"github.com/cristianoliveira/aerospace-marks/internal/stdout"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
	"github.com/spf13/cobra"
)

var markClient *storage.StorageClient
var markCmd = &cobra.Command{
	// aerospace mark
	Use:   "mark <identifier> [flags]",
	Short: "Mark a window with a specific identifier",
	Long: `Manage marks like in i3wm and Sway

Example:
aerospace-marks mark first # Will set the mark first on the current window [first]
aerospace-marks mark --add sec # Will add the mark sec to the current window [first sec]
`,

	RunE: func(cmd *cobra.Command, args []string) error {
		markClient, err := storage.NewMarkClient()
		if err != nil {
			return stdout.ErrorAndExit(err)
		}
		defer markClient.Close()


		if len(args) < 1 {
			return stdout.ErrorAndExitf("no identifier provided")
		}

		identifier := args[0]

		add, _ := cmd.Flags().GetBool("add")
		winArgID, _ := cmd.Flags().GetString("window-id")

		// Get the window ID from the command line argument
		windowID := strings.TrimSpace(winArgID)
		if winArgID == "" {
			window, err := aerospace.GetFocusedWindowID()
			if err != nil {
				return stdout.ErrorAndExit(err)
			}
			windowID = fmt.Sprintf("%d", window.WindowID)
		} else {
			window, err := aerospace.GetWindowByID(windowID)
			if err != nil {
				return stdout.ErrorAndExit(err)
			}
			windowID = fmt.Sprintf("%d", window.WindowID)
		}

		// Manage marks using MarkClient
		if add {
			err = markClient.AddMark(windowID, identifier)
			if err != nil {
				return stdout.ErrorAndExit(err)
			}

			fmt.Printf("Added mark: %s\n", identifier)
		} else {
			hasBeenDeleted, err := markClient.ReplaceAllMarks(windowID, identifier)
			if err != nil {
				return stdout.ErrorAndExit(err)
			}

			if hasBeenDeleted > 0 {
				fmt.Printf("Replaced all marks with '%s'\n", identifier)
			} else {
				fmt.Printf("Marked window with '%s'\n", identifier)
			}

		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(markCmd)

	// Define flags and configuration settings
	markCmd.Flags().Bool("add", false, "Add a mark to the window")
	markCmd.Flags().String("window-id", "", "Window ID to mark (default: focused window)")
}
