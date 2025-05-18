/*
Copyright © 2025 Cristian Oliveira licence@cristianoliveira.dev
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	newRootCmd := &cobra.Command{
		Use:   "aerospace-marks [cmd] [flags] <identifier>",
		Short: "AeroSpace marks - Marks for Aerospace WM",
		Long: `AeroSpace marks is a command line tool to manage marks for the AeroSpace WM.

This CLI is heavily inspired by the marks feature of i3 and sway window managers.
		`,
		Version: VERSION,
	}

	// Required new Mark Cmd because of leaking context
	newRootCmd.AddCommand(MarkCmd())
	newRootCmd.AddCommand(UnmarkCmd())
	newRootCmd.AddCommand(FocusCmd())
	newRootCmd.AddCommand(ListCmd())
	newRootCmd.AddCommand(ConfigCmd())

	return newRootCmd
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = NewRootCmd()

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
// NOTE: to update VERSION to empty string 
// and then run scripts/validate-version.sh
var VERSION = "v0.0.1-20250518-16d72bb"
