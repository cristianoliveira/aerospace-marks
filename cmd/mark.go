/*
Copyright © 2025 Cristian Oliveira licence@cristianoliveira.dev
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

func MarkCmd(
	storageClient storage.MarkStorage,
	aerospaceClient aerospace.AerosSpaceMarkWindows,
) *cobra.Command {
	newMarkCmd := &cobra.Command{
		// aerospace mark
		Use:   "mark <identifier> [flags]",
		Short: "Mark a window with a specific identifier",
		Long: `Mark a window with a specific identifier

mark --add|--replace [--toggle] <identifier>

Marks are arbitrary labels that can be used to identify certain windows 
and then jump to them at a later time. Each identifier can only be 
set on a single window at a time since they act as a unique identifier.
By default, mark sets identifier as the only mark on a window. --add will 
instead add identifier to the list of current marks for that window.
If --toggle is specified mark will remove identifier if it is already marked.

See: in sway manual page for more information.

Example:

aerospace-marks mark first # Will set the mark first on the current window [first]
aerospace-marks mark --add sec # Will add the mark sec to the current window [first sec]
`,

		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("No identifier provided to mark")
			}
			identifier := args[0]

			add, _ := cmd.Flags().GetBool("add")
			replace, _ := cmd.Flags().GetBool("replace")
			winArgID, _ := cmd.Flags().GetString("window-id")

			// Get the window ID from the command line argument
			windowID := strings.TrimSpace(winArgID)
			if winArgID == "" {
				window, err := aerospaceClient.Client().GetFocusedWindow()
				if err != nil {
					return stdout.ErrorAndExit(err)
				}
				windowID = fmt.Sprintf("%d", window.WindowID)
			} else {
				window, err := aerospaceClient.GetWindowByID(windowID)
				if err != nil {
					return stdout.ErrorAndExit(err)
				}
				windowID = fmt.Sprintf("%d", window.WindowID)
			}

			// Manage marks using MarkClient
			if add && !replace {
				err := storageClient.AddMark(windowID, identifier)
				if err != nil {
					return stdout.ErrorAndExit(err)
				}

				fmt.Printf("Added mark: %s\n", identifier)
				return nil
			}

			if toggle, err := cmd.Flags().GetBool("toggle"); toggle {
				if err != nil {
					return stdout.ErrorAndExit(err)
				}

				err := storageClient.ToggleMark(windowID, identifier)
				if err != nil {
					return stdout.ErrorAndExit(err)
				}

				fmt.Printf("Toggling mark: %s\n", identifier)

				return nil
			}

			hasBeenDeleted, err := storageClient.ReplaceAllMarks(windowID, identifier)
			if err != nil {
				return stdout.ErrorAndExit(err)
			}

			if hasBeenDeleted > 0 {
				fmt.Printf("Replaced all marks with '%s'\n", identifier)
			} else {
				fmt.Printf("Marked window with '%s'\n", identifier)
			}

			return nil
		},
	}

	// Define flags and configuration settings
	newMarkCmd.Flags().Bool("add", false, "Add a mark to the window")
	newMarkCmd.Flags().Bool("replace", false, "Replace all marks on the window with the new mark")
	newMarkCmd.Flags().Bool("toggle", false, "Toggle the mark on the window")
	newMarkCmd.Flags().String("window-id", "", "Window ID to mark (default: focused window)")

	return newMarkCmd
}
