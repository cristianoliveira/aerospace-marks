/*
Copyright © 2025 Cristian Oliveira license@cristianoliveira.dev
*/
package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cristianoliveira/aerospace-marks/internal/aerospace"
	"github.com/cristianoliveira/aerospace-marks/internal/cli"
	"github.com/cristianoliveira/aerospace-marks/internal/logger"
	"github.com/cristianoliveira/aerospace-marks/internal/stdout"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
	"github.com/spf13/cobra"
)

// focusCmd represents the focus command
func FocusCmd(
	storageClient storage.MarkStorage,
	aerospaceClient aerospace.AerosSpaceMarkWindows,
) *cobra.Command {
	return &cobra.Command{
		Use:   "focus <identifier> [flags]",
		Short: "Move focus to a window by mark (identifier)",
		Long: `Move focus to a window by mark (identifier)

Moves focus to the first window marked with the specified identifier.
	`,
		Args: cobra.MatchAll(
			cobra.ExactArgs(1),
			cli.ValidateArgIsNotEmpty,
		),
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.GetDefaultLogger()
			mark := args[0]
			logger.LogDebug("FocusCmd called", "mark", mark)

			windowID, err := storageClient.GetWindowIDByMark(mark)
			if err != nil {
				stdout.ErrorAndExit(err)
				return
			}

			if windowID == "" {
				stdout.ErrorAndExitf("empty window id for mark '%s'", mark)
				return
			}
			logger.LogDebug("Window found", "windowID", windowID)

			intWindowID, err := strconv.Atoi(windowID)
			if err != nil {
				stdout.ErrorAndExitf("invalid window ID '%s'", windowID)
				return
			}

			// The program is too fast, what a problem to have!
			// Sleep for half a second to ensure the window is ready
			time.Sleep(100 * time.Millisecond)

			err = aerospaceClient.Client().SetFocusByWindowID(intWindowID)
			if err != nil {
				stdout.ErrorAndExit(err)
				return
			}

			logger.LogDebug("Focus set", "windowID", windowID)
			fmt.Printf("Focus moved to window ID %s\n", windowID)
		},
	}
}
