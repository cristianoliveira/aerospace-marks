/*
Copyright © 2025 Cristian Oliveira license@cristianoliveira.dev
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/cristianoliveira/aerospace-marks/internal/aerospace"
	"github.com/cristianoliveira/aerospace-marks/internal/cli"
	"github.com/cristianoliveira/aerospace-marks/internal/format"
	"github.com/cristianoliveira/aerospace-marks/internal/logger"
	"github.com/cristianoliveira/aerospace-marks/internal/stdout"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
	"github.com/spf13/cobra"

	"github.com/cristianoliveira/aerospace-ipc/pkg/aerospace/windows"
	"github.com/cristianoliveira/aerospace-ipc/pkg/aerospace/workspaces"
)

// SummonCmd represents the summon command.
//
//nolint:gocognit,funlen // SummonCmd has high complexity and length due to workspace operations and output formatting
func SummonCmd(
	storageClient storage.MarkStorage,
	aerospaceClient aerospace.AerosSpaceMarkWindows,
) *cobra.Command {
	summonCmd := &cobra.Command{
		Use:   "summon <identifier> [flags]",
		Short: "Summon a marked window to current workspace",
		Long: `Summon a marked window to current workspace.

Similar to 'aerospace summon-workspace' but for marked windows to current workspace.
Output format can be controlled with --output flag (text, json, csv).
`,
		Args: cobra.MatchAll(
			cobra.ExactArgs(1),
			cli.ValidateArgIsNotEmpty,
		),
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.GetDefaultLogger()
			mark := args[0]
			logger.LogDebug("SummonCmd called", "mark", mark)

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

			// Move window to current workspace
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
				focusErr := aerospaceClient.Client().
					Windows().
					SetFocusByWindowID(windows.SetFocusArgs{
						WindowID: windowID,
					})
				if focusErr != nil {
					stdout.ErrorAndExit(focusErr)
					return
				}
			}

			logger.LogDebug(
				"Window summoned",
				"windowID",
				windowID,
				"workspace",
				workspace.Workspace,
			)

			// Format output using OutputEvent
			formatter, err := format.NewOutputEventFormatter(os.Stdout, outputFormat)
			if err != nil {
				stdout.ErrorAndExit(err)
				return
			}

			action := "summon"
			if shouldFocus {
				action = "summon_and_focus"
			}
			message := fmt.Sprintf(
				"Window %d summoned to workspace %s",
				windowID,
				workspace.Workspace,
			)
			if shouldFocus {
				message = fmt.Sprintf(
					"Window %d summoned to workspace %s and focused",
					windowID,
					workspace.Workspace,
				)
			}

			event := format.OutputEvent{
				Command:         "summon",
				Action:          action,
				WindowID:        windowID,
				Workspace:       workspace.Workspace,
				TargetWorkspace: workspace.Workspace,
				Result:          "success",
				Message:         message,
			}

			if formatErr := formatter.Format(event); formatErr != nil {
				stdout.ErrorAndExit(fmt.Errorf("failed to format output: %w", formatErr))
				return
			}
		},
	}

	summonCmd.Flags().BoolP("focus", "f", false, "Focus the window after summoning")
	summonCmd.Flags().
		StringP("output", "o", string(format.OutputFormatText), "Output format: text, json, or csv")
	summonCmd.Flag("output").DefValue = string(format.OutputFormatText)

	return summonCmd
}
