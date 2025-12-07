/*
Copyright © 2025 Cristian Oliveira license@cristianoliveira.dev
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/cristianoliveira/aerospace-marks/internal/aerospace"
	"github.com/cristianoliveira/aerospace-marks/internal/cli"
	"github.com/cristianoliveira/aerospace-marks/internal/format"
	"github.com/cristianoliveira/aerospace-marks/internal/logger"
	"github.com/cristianoliveira/aerospace-marks/internal/stdout"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
	"github.com/spf13/cobra"

	"github.com/cristianoliveira/aerospace-ipc/pkg/aerospace/windows"
)

//nolint:gochecknoglobals // focusDelay is a configuration constant
var focusDelay = 100 * time.Millisecond // Default delay to wait for the window to be ready

// FocusCmd represents the focus command.
func FocusCmd(
	storageClient storage.MarkStorage,
	aerospaceClient aerospace.AerosSpaceMarkWindows,
) *cobra.Command {
	focusCmd := &cobra.Command{
		Use:   "focus <identifier> [flags]",
		Short: "Move focus to a window by mark (identifier)",
		Long: `Move focus to a window by mark (identifier)

Moves focus to the first window marked with the specified identifier.
Output format can be controlled with --output flag (text, json, csv).
	`,
		Args: cobra.MatchAll(
			cobra.ExactArgs(1),
			cli.ValidateArgIsNotEmpty,
		),
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.GetDefaultLogger()
			mark := args[0]
			logger.LogDebug("FocusCmd called", "mark", mark)

			// Get and validate output format early
			outputFormat, err := cmd.Flags().GetString("output")
			if err != nil {
				stdout.ErrorAndExit(fmt.Errorf("failed to get output flag: %w", err))
				return
			}

			// Default to text if not specified
			if outputFormat == "" {
				outputFormat = string(format.OutputFormatText)
			}

			windowID, err := storageClient.GetWindowIDByMark(mark)
			if err != nil {
				stdout.ErrorAndExit(err)
				return
			}

			if windowID == 0 {
				stdout.ErrorAndExitf("empty window id for mark '%s'", mark)
				return
			}
			logger.LogDebug("Window found", "windowID", windowID)

			// The program is too fast, what a problem to have!
			// Delay setting focus to ensure the window is ready
			time.Sleep(focusDelay)
			err = aerospaceClient.Client().Windows().SetFocusByWindowID(windows.SetFocusArgs{
				WindowID: windowID,
			})
			if err != nil {
				stdout.ErrorAndExit(err)
				return
			}

			logger.LogDebug("Focus set", "windowID", windowID)

			// Format output using OutputEvent
			formatter, err := format.NewOutputEventFormatter(os.Stdout, outputFormat)
			if err != nil {
				stdout.ErrorAndExit(err)
				return
			}

			message := fmt.Sprintf("Focus moved to window ID %d", windowID)
			event := format.OutputEvent{
				Command:  "focus",
				Action:   "focus",
				WindowID: windowID,
				Result:   "success",
				Message:  message,
			}

			if formatErr := formatter.Format(event); formatErr != nil {
				stdout.ErrorAndExit(fmt.Errorf("failed to format output: %w", formatErr))
				return
			}
		},
	}

	// Add output flag
	focusCmd.Flags().
		StringP("output", "o", string(format.OutputFormatText), "Output format: text, json, or csv")
	focusCmd.Flag("output").DefValue = string(format.OutputFormatText)

	return focusCmd
}
