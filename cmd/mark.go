/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/cristianoliveira/aerospace-ext/internal/storage"
	"github.com/cristianoliveira/aerospace-ext/internal/aerospace"
	"github.com/spf13/cobra"
)

var markClient *storage.MarkClient
var markCmd = &cobra.Command{
	// aerospace mark 
	Use:   "mark",
	Short: "Manage marks like in i3wm and Sway",
	Long: `Manage marks like in i3wm and Sway

See: man 5 sway

mark --add|--replace [--toggle] <identifier>

Marks are arbitrary labels that can be used to identify certain windows 
and then jump to them at a later time. Each identifier can only be 
set on a single window at a time since they act as a unique identifier.
By default, mark sets identifier as the only mark on a window. --add will 
instead add identifier to the list of current marks for that window.
If --toggle is specified mark will remove identifier if it is already marked.
`,
	Run: func(cmd *cobra.Command, args []string) {
		markClient, err := storage.NewMarkClient()
		if err != nil {
			fmt.Printf("Error initializing mark client: %v\n", err)
			return
		}
		defer markClient.Close()

		getID, _ := cmd.Flags().GetString("get-id")
		add, _ := cmd.Flags().GetBool("add")
		winArgID, _ := cmd.Flags().GetString("window")
		replace, _ := cmd.Flags().GetBool("replace")
		toggle, _ := cmd.Flags().GetBool("toggle")

		if getID != "" {
			// Get window ID by mark
			windowID, err := markClient.GetWindowIDByMark(getID)
			if err != nil {
				fmt.Printf("Error retrieving window ID for mark '%s': %v\n", getID, err)
				return
			}
			if windowID == "" {
				fmt.Printf("No window found for mark: %s\n", getID)
			} else {
				fmt.Printf("Window ID for mark '%s': %s\n", getID, windowID)
			}
			return
		} else if replace || toggle {
			panic("replace and toggle are not implemented yet")
		}

		// Validate arguments
		var identifier string
		if add {
			if len(args) < 1 {
				fmt.Println("Error: identifier is required")
				return
			} else {
				identifier = args[0]
			}
		}

		// Get the window ID from the command line argument
		windowID := strings.TrimSpace(winArgID)

		// Manage marks using MarkClient
		if add {
			// Add logic
			if winArgID == "" {
				windowQuery, err := aerospace.GetFocusedWindowID()
				if err != nil {
					fmt.Printf("Error getting focused window ID: %v\n", err)
					return
				}

				windowID = strings.TrimSpace(windowQuery)
			} else {
				windowID = strings.TrimSpace(winArgID)
			}

			err = markClient.AddMark(windowID, identifier)
			if err != nil {
				fmt.Printf("Error adding mark: %v\n", err)
				return
			}
			fmt.Printf("Added mark: %s\n", identifier)
		} else {
			// List the marks
			var marks []storage.Mark
			// If no identifier is provided, get all marks
			if winArgID == "" {
				marks, err = markClient.GetMarks()
				if err != nil {
					fmt.Printf("Error getting marks: %v\n", err)
					return
				}
			} else {
				// Get marks for the specified window ID
				marks, err = markClient.GetMarksByWindowID(windowID)
				if err != nil {
					fmt.Printf("Error getting marks: %v\n", err)
					return
				}
			}

			if len(marks) == 0 {
				fmt.Println("No marks found")
				return
			}

			fmt.Println("Marks:")
			for _, mark := range marks {
				fmt.Printf(" - %s\n", mark.Mark)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(markCmd)

	// Define flags and configuration settings
	markCmd.Flags().String("get-id", "", "Get the window ID associated with the specified mark")
	markCmd.Flags().Bool("add", false, "Add a mark to the window")
	markCmd.Flags().Bool("replace", false, "Replace all marks on the window")
	markCmd.Flags().Bool("toggle", false, "Toggle the mark on the window")
	markCmd.Flags().StringP("window", "w", "", "Window ID to mark (default: focused window)")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// markCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// markCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
