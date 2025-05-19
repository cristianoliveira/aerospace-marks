/*
Copyright © 2025 Cristian Oliveira licence@cristianoliveira.dev
*/
package cmd

import (
	"os"

	"github.com/cristianoliveira/aerospace-marks/internal/storage"
	"github.com/spf13/cobra"
)

func NewRootCmd(storage storage.MarkStorage) *cobra.Command {
	newRootCmd := &cobra.Command{
		Use:   "aerospace-marks [cmd] [flags] <identifier>",
		Short: "AeroSpace marks - Marks for Aerospace WM",
		Long: `AeroSpace marks is a command line tool to manage marks for the AeroSpace WM.

This CLI is heavily inspired by the marks feature of i3 and sway window managers.
		`,
		Version: VERSION,
	}

	// Required new Mark Cmd because of leaking context
	newRootCmd.AddCommand(MarkCmd(storage))
	newRootCmd.AddCommand(UnmarkCmd(storage))
	newRootCmd.AddCommand(FocusCmd(storage))
	newRootCmd.AddCommand(ListCmd(storage))
	newRootCmd.AddCommand(ConfigCmd())

	return newRootCmd
}

func init() {
	// NOTE: add here global flags
	// rootCmd.Flags().BoolP("version", "v", false, "Print version information")
}

func Run(storage storage.MarkStorage) {
	rootCmd := NewRootCmd(storage)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// THIS IS GENERATED DON'T EDIT
// NOTE: to update VERSION to empty string 
// and then run scripts/validate-version.sh
var VERSION = "v0.0.1-20250518-16d72bb"
