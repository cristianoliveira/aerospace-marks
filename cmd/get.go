/*
Copyright © 2025 Cristian Oliveira licence@cristianoliveira.dev
*/
package cmd

import (
	"fmt"

	"github.com/cristianoliveira/aerospace-marks/internal/aerospace"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("missing mark (identifier)")
			}

			mark := args[0]

			// Get window ID by mark
			windows, err := storageClient.GetMarksByWindowID(mark)
			if err != nil {
				return stdout.ErrorAndExit(err)
			}
			if len(windows) == 0 {
				err = fmt.Errorf("no window found for mark %s", mark)
				return stdout.ErrorAndExit(err)
			}

			windowID := windows[0].WindowID

			getWinId, _ := cmd.Flags().GetBool("window-id")
			if getWinId {
				fmt.Print(windowID)
				return nil
			}

			getWinTitle, _ := cmd.Flags().GetBool("window-title")
			getWinApp, _ := cmd.Flags().GetBool("app-name")

			window, err := aerospaceClient.GetWindowByID(windowID)
			if err != nil {
				return stdout.ErrorAndExit(err)
			}
			if window == nil {
				return stdout.ErrorAndExit(fmt.Errorf("no window found for ID %s", windowID))
			}

			if getWinTitle {
				fmt.Print(window.WindowTitle)
				return nil
			}

			if getWinApp {
				fmt.Print(window.AppName)
				return nil
			}

			fmt.Print(window.String())
			return nil
		},
	}

	getCmd.Flags().BoolP("window-id", "i", false, "Get only window [i]D")
	getCmd.Flags().BoolP("window-title", "t", false, "Get only window [t]itle")
	getCmd.Flags().BoolP("app-name", "a", false, "Get only window [a]pp name")

	return getCmd
}
