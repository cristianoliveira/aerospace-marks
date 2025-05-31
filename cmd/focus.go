/*
Copyright © 2025 Cristian Oliveira license@cristianoliveira.dev
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/cristianoliveira/aerospace-marks/internal/aerospace"
	"github.com/cristianoliveira/aerospace-marks/internal/cli"
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
			mark := args[0]

			windowID, err := storageClient.GetWindowIDByMark(mark)
			if err != nil {
				stdout.ErrorAndExit(err)
				return
			}
			if windowID == "" {
				stdout.ErrorAndExitf("empty window id for mark '%s'", mark)
				return
			}

			intWindowID, err := strconv.Atoi(windowID)
			if err != nil {
				stdout.ErrorAndExitf("invalid window ID '%s'", windowID)
				return
			}
			err = aerospaceClient.Client().SetFocusByWindowID(intWindowID)
			if err != nil {
				stdout.ErrorAndExit(err)
				return
			}

			fmt.Printf("Focus moved to window ID %s\n", windowID)
		},
	}
}
