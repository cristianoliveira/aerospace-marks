/*
Copyright © 2025 Cristian Oliveira license@cristianoliveira.dev
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/cristianoliveira/aerospace-marks/internal/aerospace"
	"github.com/cristianoliveira/aerospace-marks/internal/cli"
	"github.com/cristianoliveira/aerospace-marks/internal/logger"
	"github.com/cristianoliveira/aerospace-marks/internal/stdout"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
	"github.com/spf13/cobra"
)

func GetCmd(
	storageClient storage.MarkStorage,
	aerospaceClient aerospace.AerosSpaceMarkWindows,
) *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get <identifier>",
		Short: "Get a window by mark (identifier)",
		Long: `Get a window by mark (identifier)

This command retrieves a window by its mark (identifier). Print in the following format:

<window_id> | <window_title> | <app_name>
`,
		Args: cobra.MatchAll(
			cobra.ExactArgs(1),
			cli.ValidateArgIsNotEmpty,
		),
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.GetDefaultLogger()

			mark := args[0]
			logger.LogDebug("Getting window by mark: %s", mark)

			markedWindow, err := storageClient.GetWindowByMark(mark)
			if err != nil {
				stdout.ErrorAndExit(err)
				return
			}

			windowID := markedWindow.WindowID
			if windowID == 0 {
				stdout.ErrorAndExit(fmt.Errorf("no window found for mark %s", mark))
				return
			}

			getWinID, _ := cmd.Flags().GetBool("window-id")
			if getWinID {
				fmt.Fprint(os.Stdout, windowID)
				return
			}

			window, err := aerospaceClient.GetWindowByID(windowID)
			if err != nil {
				stdout.ErrorAndExit(err)
				return
			}
			if window == nil {
				stdout.ErrorAndExit(fmt.Errorf("no window found for ID %d", windowID))
				return
			}
			logger.LogDebug(
				"Get window by ID",
				"windowID", windowID,
				"window", *window,
				"windowTitle", window.WindowTitle,
				"err", err,
			)

			getWinTitle, _ := cmd.Flags().GetBool("window-title")
			if getWinTitle {
				fmt.Fprint(os.Stdout, window.WindowTitle)
				return
			}

			getWinApp, _ := cmd.Flags().GetBool("app-name")
			if getWinApp {
				fmt.Fprint(os.Stdout, window.AppName)
				return
			}

			getWinAppBundleID, _ := cmd.Flags().GetBool("app-bundle-id")
			if getWinAppBundleID {
				fmt.Fprint(os.Stdout, window.AppBundleID)
				return
			}

			logger.LogInfo(
				"Printing full window info",
				"windowID", windowID,
				"window", window,
			)
			fmt.Fprint(os.Stdout, window)
		},
	}

	getCmd.Flags().BoolP("window-id", "i", false, "Get only window [i]D")
	getCmd.Flags().BoolP("window-title", "t", false, "Get only window [t]itle")
	getCmd.Flags().BoolP("app-name", "a", false, "Get only window [a]pp name")
	getCmd.Flags().BoolP("app-bundle-id", "b", false, "Get only window app [b]undle ID")

	return getCmd
}
