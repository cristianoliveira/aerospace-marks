/*
Copyright © 2025 Cristian Oliveira license@cristianoliveira.dev
*/
package cmd

import (
	"fmt"

	aerospacecli "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace"
	"github.com/cristianoliveira/aerospace-marks/internal/aerospace"
	"github.com/cristianoliveira/aerospace-marks/internal/format"
	"github.com/cristianoliveira/aerospace-marks/internal/stdout"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
	"github.com/spf13/cobra"
	"slices"
)

func popWindow(windows []aerospacecli.Window, windowID int) (*aerospacecli.Window, error) {
	for i, window := range windows {
		if windowID == 0 {
			return nil, fmt.Errorf("window ID not found")
		}
		if windowID == window.WindowID {
			// Remove the window from the list
			//nolint:staticcheck,ineffassign
			windows = slices.Delete(windows, i, i+1)
			return &window, nil
		}
	}

	return nil, fmt.Errorf("window ID not found")
}

// listCmd represents the list command
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
				fmt.Println("No marks found")
				return
			}

			windows, err := aerospaceClient.Client().Windows().GetAllWindows()
			if err != nil {
				stdout.ErrorAndExit(err)
				return
			}

			lines := make([]string, 0)
			for _, mark := range marks {
				window, err := popWindow(windows, mark.WindowID)
				if err != nil {
					continue
				}

				// Format window fields, using "_" for empty fields
				windowID := fmt.Sprintf("%d", window.WindowID)
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
				fmt.Println("No marked window found")
				return
			}

			formattedOutput := format.FormatTableList(lines)
			fmt.Println(formattedOutput)
		},
	}
}
