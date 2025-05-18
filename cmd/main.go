/*
Copyright © 2025 Cristian Oliveira me@cristianoliveira.dev
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aerospace-marks [cmd] [flags] <identifier>",
	Short: "AeroSpace marks - Marks for Aerospace WM",
	Long: `AeroSpace marks is a command line tool to manage marks for the AeroSpace WM.

mark --add|--replace [--toggle] <identifier>

	Marks are arbitrary labels that can be used to identify certain windows 
	and then jump to them at a later time. Each identifier can only be 
	set on a single window at a time since they act as a unique identifier.
	By default, mark sets identifier as the only mark on a window. --add will 
	instead add identifier to the list of current marks for that window.
	If --toggle is specified mark will remove identifier if it is already marked.

See: man 5 sway
`,
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Run() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
