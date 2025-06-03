/*
Copyright © 2025 Cristian Oliveira license@cristianoliveira.dev
*/
package cmd

import (
	"os"

	"github.com/cristianoliveira/aerospace-marks/internal/aerospace"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
	"github.com/spf13/cobra"
)

func NewRootCmd(
	storage storage.MarkStorage,
	aerospaceClient aerospace.AerosSpaceMarkWindows,
) *cobra.Command {
	newRootCmd := &cobra.Command{
		Use:   "aerospace-marks [cmd] [flags] <identifier>",
		Short: "AeroSpace marks - Marks for Aerospace WM",
		Long: `AeroSpace marks is a command line tool to manage marks for the AeroSpace WM.

This CLI is heavily inspired by the marks feature of i3 and sway window managers.
		`,
		Version: VERSION,
	}

	// Required new Mark Cmd because of leaking context
	newRootCmd.AddCommand(MarkCmd(storage, aerospaceClient))
	newRootCmd.AddCommand(UnmarkCmd(storage))
	newRootCmd.AddCommand(FocusCmd(storage, aerospaceClient))
	newRootCmd.AddCommand(ListCmd(storage, aerospaceClient))
	newRootCmd.AddCommand(ConfigCmd())
	newRootCmd.AddCommand(SummonCmd(storage, aerospaceClient))
	newRootCmd.AddCommand(GetCmd(storage, aerospaceClient))

	return newRootCmd
}

func init() {
	// NOTE: add here global flags
	// rootCmd.Flags().BoolP("version", "v", false, "Print version information")
}

func Run(storage storage.MarkStorage, aerospaceClient aerospace.AerosSpaceMarkWindows) {
	rootCmd := NewRootCmd(storage, aerospaceClient)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// THIS IS GENERATED DON'T EDIT
// NOTE: to update VERSION to empty string
// and then run scripts/validate-version.sh
var VERSION = "v0.2.1"
