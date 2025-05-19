/*
Copyright © 2025 Cristian Oliveira licence@cristianoliveira.dev
*/
package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cristianoliveira/aerospace-marks/internal/aerospace"
	"github.com/cristianoliveira/aerospace-marks/internal/format"
	"github.com/cristianoliveira/aerospace-marks/internal/stdout"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
	"github.com/cristianoliveira/aerospace-marks/pkgs/aerospacecli"
	"github.com/spf13/cobra"
)

func popWindow(windows []aerospacecli.Window, windowID string) (*aerospacecli.Window, error) {
	for i, window := range windows {
		if windowID == "" {
			return nil, fmt.Errorf("window ID not found")
		}
		winId := strconv.Itoa(window.WindowID)
		winId = strings.TrimSpace(winId)
		if windowID == winId {
			// Remove the window from the list
			windows = append(windows[:i], windows[i+1:]...)
			return &window, nil
		}
	}

	return nil, fmt.Errorf("window ID not found")
}

// listCmd represents the list command
func ListCmd(storageClient storage.MarkStorage) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all marked windows",
		Long: `List all marked windows

This command lists all marked windows with their respective marks.
Display in the following format:

<mark>|<window-id>|<window-title>|<window-app>
	`,
		RunE: func(cmd *cobra.Command, args []string) error {
			marks, err := storageClient.GetMarks()
			if err != nil {
				return stdout.ErrorAndExit(err)
			}
			if len(marks) == 0 {
				fmt.Println("No marks found")
				return nil
			}

			windows, err := aerospace.GetAllWindows()
			if err != nil {
				return stdout.ErrorAndExit(err)
			}

			lines := make([]string, 0)
			for _, mark := range marks {
				window, err := popWindow(windows, mark.WindowID)
				if err != nil {
					continue
				}

				line := fmt.Sprintf("%s|%s\r\n", mark.Mark, window)
				lines = append(lines, line)
			}

			if len(lines) == 0 {
				fmt.Println("No marked window found")
				return nil
			}

			formattedOutput := format.FormatTableList(lines)
			fmt.Println(formattedOutput)

			return nil
		},
	}
}
