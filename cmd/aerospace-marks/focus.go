/*
Copyright © 2025 Cristian Oliveira me@cristianoliveira.dev
*/
package main

import (
	"fmt"

	"github.com/cristianoliveira/aerospace-marks/internal/aerospace"
	"github.com/cristianoliveira/aerospace-marks/internal/stdout"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
	"github.com/spf13/cobra"
)

// focusCmd represents the focus command
var focusCmd = &cobra.Command{
	Use:   "focus <mark> [flags]",
	Short: "Focus to a window by mark",
	Long: `Focus to a window by mark

USAGE:
aerospace-marks focus <mark> # Will focus to FIRST window marked with <mark>
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		markClient, err := storage.NewMarkClient()
		if err != nil {
			return stdout.ErrorAndExit(err)
		}
		defer markClient.Close()

		if len(args) < 1 {
			return stdout.ErrorAndExitf("No mark provided")
		}

		mark := args[0]

		// Get window ID by mark
		windowID, err := markClient.GetWindowIDByMark(mark)
		if err != nil {
			return stdout.ErrorAndExit(err)
		}
		if windowID == "" {
			return stdout.ErrorAndExitf("no window found for mark '%s'", mark)
		}

		// Focus to window by ID
		err = aerospace.SetFocusToWindowId(windowID)
		if err != nil {
			return stdout.ErrorAndExit(err)
		}

		fmt.Printf("Focus moved to window ID %s\n", windowID)

		return nil;
	},
}

func init() {
	rootCmd.AddCommand(focusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// focusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// focusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
