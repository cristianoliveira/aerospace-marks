/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/cristianoliveira/aerospace-marks/internal/aerospace"
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
`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			mark := args[0]

			// Get window ID by mark
			windowID, err := storageClient.GetWindowIDByMark(mark)
			if err != nil {
				return
			}
			if windowID == "" {
				return
			}

			getWinId, _ := cmd.Flags().GetBool("window-id")
			if getWinId {
				fmt.Print(windowID)
				return
			}

			getWinTitle, _ := cmd.Flags().GetBool("window-title")
			getWinApp, _ := cmd.Flags().GetBool("window-app")

			window, err := aerospaceClient.GetWindowByID(windowID)
			if getWinTitle {
				fmt.Print(window.WindowTitle)
				return
			}

			if getWinApp {
				fmt.Print(window.AppName)
				return
			}

			fmt.Print(window.String())
		},
	}

	getCmd.Flags().BoolP("window-id", "i", false, "Window ID to get")
	getCmd.Flags().BoolP("window-title", "t", false, "Window title to get")
	getCmd.Flags().BoolP("window-app", "a", false, "Window app to get")

	return getCmd
}
