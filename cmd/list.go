/*
Copyright © 2025 Cristian Oliveira me@cristianoliveira.dev
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/cristianoliveira/aerospace-marks/internal/aerospace"
	"github.com/cristianoliveira/aerospace-marks/internal/format"
	"github.com/cristianoliveira/aerospace-marks/internal/stdout"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
	"github.com/spf13/cobra"
)

func popWindow(windows []string, windowID string) (string, error) {
	for i, window := range windows {
		if windowID == "" {
			return "", fmt.Errorf("window ID not found")
		}
		winId := window[:strings.Index(window, "|")]
		winId = strings.TrimSpace(winId)
		if windowID == winId {
			// Remove the window from the list
			windows = append(windows[:i], windows[i+1:]...)
			return window, nil
		}
	}
	return "", fmt.Errorf("window ID not found")
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all windows marked",
	RunE: func(cmd *cobra.Command, args []string) error {
		markClient, err := storage.NewMarkClient()
		if err != nil {
			return stdout.ErrorAndExit(err)
		}
		defer markClient.Close()

		marks, err := markClient.GetMarks()
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

		formattedOutput := format.FormatTableList(lines)
		fmt.Println(formattedOutput)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
