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
	Version: VERSION,
}

func init() {
	// NOTE: add here global flags
	// rootCmd.Flags().BoolP("version", "v", false, "Print version information")
}

func Run() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// THIS IS GENERATED DON'T EDIT
// NOTE: to update delete the var and
// run scripts/validate-version.sh
var VERSION = "v0.0.1-20250518-75aee46"
