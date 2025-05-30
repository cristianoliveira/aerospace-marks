/*
Copyright © 2025 Cristian Oliveira license@cristianoliveira.dev
*/
package cmd

import (
	"strconv"

	"github.com/cristianoliveira/aerospace-marks/internal/aerospace"
	"github.com/cristianoliveira/aerospace-marks/internal/cli"
	"github.com/cristianoliveira/aerospace-marks/internal/stdout"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
	"github.com/spf13/cobra"
)

// SummonCmd represents the summon command
func SummonCmd(
	storageClient storage.MarkStorage,
	aerospaceClient aerospace.AerosSpaceMarkWindows,
) *cobra.Command {
	summonCmd := &cobra.Command{
		Use:   "summon <identifier> [flags]",
		Short: "Summon a marked window to current workspace",
		Long: `Summon a marked window to current workspace.

Similar to 'aerospace summon-workspace' but for marked windows to current workspace.
`,
		Args:	cobra.MatchAll(
			cobra.ExactArgs(1),
			cli.ValidateArgIsNotEmpty,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			mark := args[0]

			shouldFocus, err := cmd.Flags().GetBool("focus")
			if err != nil {
				return stdout.ErrorAndExit(err)
			}

			// Get window ID by mark
			windowID, err := storageClient.GetWindowIDByMark(mark)
			if err != nil {
				return stdout.ErrorAndExit(err)
			}
			if windowID == "" {
				return stdout.ErrorAndExitf("no window found for mark '%s'", mark)
			}

			workspace, err := aerospaceClient.Client().GetFocusedWorkspace()
			if err != nil {
				return stdout.ErrorAndExit(err)
			}

			// FIXME: windowsID as number
			intWindowID, err := strconv.Atoi(windowID)
			if err != nil {
				return stdout.ErrorAndExitf("invalid window ID '%s'", windowID)
			}

			// focus to window by ID
			err = aerospaceClient.Client().MoveWindowToWorkspace(intWindowID, workspace.Workspace)
			if err != nil {
				return stdout.ErrorAndExit(err)
			}

			if shouldFocus {
				aerospaceClient.Client().SetFocusByWindowID(intWindowID)
			}

			return nil
		},
	}

	summonCmd.Flags().BoolP("focus", "f", false, "Focus the window after summoning")

	return summonCmd
}
