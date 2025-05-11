/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var marks = make(map[string]string) // Map to store marks with window IDs
var markCmd = &cobra.Command{
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
		// Retrieve flags
		add, _ := cmd.Flags().GetBool("add")
		replace, _ := cmd.Flags().GetBool("replace")
		toggle, _ := cmd.Flags().GetBool("toggle")

		// Validate arguments
		if len(args) < 1 {
			fmt.Println("Error: identifier is required")
			return
		}
		identifier := args[0]

		// Call external command to get focused window ID
		out, err := exec.Command("aerospace", "list-windows", "--focused").Output()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		windowID := strings.Fields(string(out))[0]

		// Manage marks in memory
		if toggle {
			// Toggle logic
			if _, exists := marks[identifier]; exists {
				delete(marks, identifier)
				fmt.Printf("Removed mark: %s\n", identifier)
			} else {
				marks[identifier] = windowID
				fmt.Printf("Added mark: %s\n", identifier)
			}
		} else if add {
			// Add logic
			marks[identifier] = windowID
			fmt.Printf("Added mark: %s\n", identifier)
		} else if replace {
			// Replace logic
			for id, win := range marks {
				if win == windowID {
					delete(marks, id)
				}
			}
			marks[identifier] = windowID
			fmt.Printf("Replaced mark: %s\n", identifier)
		} else {
			// Default behavior
			for id, win := range marks {
				if win == windowID {
					delete(marks, id)
				}
			}
			marks[identifier] = windowID
			fmt.Printf("Set mark: %s\n", identifier)
		}
	},
}

func init() {
	rootCmd.AddCommand(markCmd)

	// Define flags and configuration settings
	markCmd.Flags().Bool("add", false, "Add a mark to the window")
	markCmd.Flags().Bool("replace", false, "Replace the mark on the window")
	markCmd.Flags().Bool("toggle", false, "Toggle the mark on the window")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// markCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// markCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
