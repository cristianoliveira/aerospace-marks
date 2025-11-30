/*
Copyright © 2025 Cristian Oliveira license@cristianoliveira.dev
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"

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
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all marked windows",
		Long: `List all marked windows

This command lists all marked windows with their respective marks.
Display in the following format:

<mark>|<window-id>|<app-name>|<window-title>|<workspace>|<app-bundle-id>
	`,
		Run: func(cmd *cobra.Command, args []string) {
			marks, err := storageClient.GetMarks()
			if err != nil {
				stdout.ErrorAndExit(err)
				return
			}
			if len(marks) == 0 {
				fmt.Fprintln(os.Stdout, "No marks found")
				return
			}

			windowsList, err := aerospaceClient.Client().Windows().GetAllWindows()
			if err != nil {
				stdout.ErrorAndExit(err)
				return
			}

			lines := make([]string, 0)
			for _, mark := range marks {
				window, popErr := popWindow(windowsList, mark.WindowID)
				if popErr != nil {
					continue
				}

				// Format window fields, using "_" for empty fields
				windowID := strconv.Itoa(window.WindowID)
				appName := window.AppName
				if appName == "" {
					appName = "_"
				}
				windowTitle := window.WindowTitle
				if windowTitle == "" {
					windowTitle = "_"
				}
				workspace := window.Workspace
				if workspace == "" {
					workspace = "_"
				}
				appBundleID := window.AppBundleID
				if appBundleID == "" {
					appBundleID = "_"
				}

				line := fmt.Sprintf("%s|%s|%s|%s|%s|%s\r\n",
					mark.Mark, windowID, appName, windowTitle, workspace, appBundleID)
				lines = append(lines, line)
			}

			if len(lines) == 0 {
				fmt.Fprintln(os.Stdout, "No marked window found")
				return
			}

			formattedOutput := format.FormatTableList(lines)
			fmt.Fprintln(os.Stdout, formattedOutput)
		},
	}
}
