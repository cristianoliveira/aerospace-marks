/*
Copyright © 2025 Cristian Oliveira license@cristianoliveira.dev
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/cristianoliveira/aerospace-ipc/pkg/aerospace/windows"
	"github.com/cristianoliveira/aerospace-marks/internal/aerospace"
	"github.com/cristianoliveira/aerospace-marks/internal/cli"
	"github.com/cristianoliveira/aerospace-marks/internal/logger"
	"github.com/cristianoliveira/aerospace-marks/internal/stdout"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
	"github.com/spf13/cobra"
)

var focusDelay = 100 * time.Millisecond // Default delay to wait for the window to be ready

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
			fmt.Printf("Focus moved to window ID %d\n", windowID)
		},
	}
}
