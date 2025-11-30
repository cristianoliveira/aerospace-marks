/*
Copyright © 2025 Cristian Oliveira license@cristianoliveira.dev
*/
package cmd

import (
	"github.com/cristianoliveira/aerospace-ipc/pkg/aerospace/windows"
	"github.com/cristianoliveira/aerospace-ipc/pkg/aerospace/workspaces"
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
		Args: cobra.MatchAll(
			cobra.ExactArgs(1),
			cli.ValidateArgIsNotEmpty,
		),
		Run: func(cmd *cobra.Command, args []string) {
			mark := args[0]

			shouldFocus, err := cmd.Flags().GetBool("focus")
			if err != nil {
				stdout.ErrorAndExit(err)
				return
			}

			// Get window ID by mark
			windowID, err := storageClient.GetWindowIDByMark(mark)
			if err != nil {
				stdout.ErrorAndExit(err)
				return
			}
			if windowID == 0 {
				stdout.ErrorAndExitf("no window found for mark '%s'", mark)
				return
			}

			workspace, err := aerospaceClient.Client().Workspaces().GetFocusedWorkspace()
			if err != nil {
				stdout.ErrorAndExit(err)
				return
			}

			// focus to window by ID
			err = aerospaceClient.Client().Workspaces().MoveWindowToWorkspaceWithOpts(
				workspaces.MoveWindowToWorkspaceArgs{
					WorkspaceName: workspace.Workspace,
				},
				workspaces.MoveWindowToWorkspaceOpts{
					WindowID: &windowID,
				},
			)
			if err != nil {
				stdout.ErrorAndExit(err)
				return
			}

			if shouldFocus {
				err := aerospaceClient.Client().Windows().SetFocusByWindowID(windows.SetFocusArgs{
					WindowID: windowID,
				})
				if err != nil {
					stdout.ErrorAndExit(err)
					return
				}
			}
		},
	}

	summonCmd.Flags().BoolP("focus", "f", false, "Focus the window after summoning")

	return summonCmd
}
