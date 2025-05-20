/*
Copyright © 2025 Cristian Oliveira licence@cristianoliveira.dev
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/cristianoliveira/aerospace-marks/internal/aerospace"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("No identifier provided to focus")
			}

			mark := args[0]

			// Get window ID by mark
			windowID, err := storageClient.GetWindowIDByMark(mark)
			if err != nil {
				return stdout.ErrorAndExit(err)
			}
			if windowID == "" {
				return stdout.ErrorAndExitf("no window found for mark '%s'", mark)
			}

			// Focus to window by ID
			intWindowID, err := strconv.Atoi(windowID)
			if err != nil {
				return stdout.ErrorAndExitf("invalid window ID '%s'", windowID)
			}
			err = aerospaceClient.Client().SetFocusByWindowID(intWindowID)
			if err != nil {
				return stdout.ErrorAndExit(err)
			}

			fmt.Printf("Focus moved to window ID %s\n", windowID)

			return nil
		},
	}
}
