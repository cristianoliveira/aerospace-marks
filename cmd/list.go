/*
Copyright © 2025 Cristian Oliveira license@cristianoliveira.dev
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"slices"

	"github.com/cristianoliveira/aerospace-marks/internal/aerospace"
	"github.com/cristianoliveira/aerospace-marks/internal/format"
	"github.com/cristianoliveira/aerospace-marks/internal/stdout"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
	"github.com/spf13/cobra"

	"github.com/cristianoliveira/aerospace-ipc/pkg/aerospace/windows"
)

func popWindow(windowsList []windows.Window, windowID int) (*windows.Window, error) {
	for i, window := range windowsList {
		if windowID == 0 {
			return nil, errors.New("window ID not found")
		}
		if windowID == window.WindowID {
			// Remove the window from the list
			// The assignment is intentional to show intent, even though the slice is passed by value
			//nolint:staticcheck,ineffassign,wastedassign // intentional: shows removal intent
			windowsList = slices.Delete(windowsList, i, i+1)
			return &window, nil
		}
	}

	return nil, errors.New("window ID not found")
}

// ListCmd represents the list command.
//
//nolint:gocognit // ListCmd has high complexity due to multiple formatting operations
func ListCmd(
	storageClient storage.MarkStorage,
	aerospaceClient aerospace.AerosSpaceMarkWindows,
) *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all marked windows",
		Long: `List all marked windows

This command lists all marked windows with their respective marks.
Display format can be controlled with --output flag (text, json, csv).

Default format (text):
<mark>|<window-id>|<app-name>|<window-title>|<workspace>|<app-bundle-id>
	`,
		Run: func(cmd *cobra.Command, args []string) {
			// Get and validate output format early
			outputFormat, err := cmd.Flags().GetString("output")
			if err != nil {
				stdout.ErrorAndExit(fmt.Errorf("failed to get output flag: %w", err))
				return
			}

			// Default to text if not specified
			if outputFormat == "" {
				outputFormat = "text"
			}

			// Validate format before any processing
			formatter, err := format.NewListOutputFormatter(os.Stdout, outputFormat)
			if err != nil {
				stdout.ErrorAndExit(err)
				return
			}

			// Get marks from storage
			marks, err := storageClient.GetMarks()
			if err != nil {
				stdout.ErrorAndExit(err)
				return
			}

			// Handle empty marks based on format
			if len(marks) == 0 {
				if formatErr := formatter.FormatEmpty("No marks found"); formatErr != nil {
					stdout.ErrorAndExit(fmt.Errorf("failed to format empty output: %w", formatErr))
					return
				}
				return
			}

			// Get windows from Aerospace
			windowsList, err := aerospaceClient.Client().Windows().GetAllWindows()
			if err != nil {
				stdout.ErrorAndExit(err)
				return
			}

			// Collect marked windows
			markedWindows := make([]format.MarkedWindow, 0)
			for _, mark := range marks {
				window, popErr := popWindow(windowsList, mark.WindowID)
				if popErr != nil {
					// Silently skip windows that no longer exist
					continue
				}

				markedWindows = append(markedWindows, format.MarkedWindow{
					Mark:        mark.Mark,
					WindowID:    window.WindowID,
					AppName:     window.AppName,
					WindowTitle: window.WindowTitle,
					Workspace:   window.Workspace,
					AppBundleID: window.AppBundleID,
				})
			}

			// Handle empty marked windows based on format
			if len(markedWindows) == 0 {
				if formatErr := formatter.FormatEmpty("No marked window found"); formatErr != nil {
					stdout.ErrorAndExit(fmt.Errorf("failed to format empty output: %w", formatErr))
					return
				}
				return
			}

			// Format and output
			if formatErr := formatter.Format(markedWindows); formatErr != nil {
				stdout.ErrorAndExit(fmt.Errorf("failed to format output: %w", formatErr))
				return
			}
		},
	}

	// Add output flag
	listCmd.Flags().StringP("output", "o", "text", "Output format: text, json, or csv")
	listCmd.Flag("output").DefValue = "text"

	return listCmd
}
