/*
Copyright © 2025 Cristian Oliveira license@cristianoliveira.dev
*/
package cmd

import (
	"os"

	"github.com/cristianoliveira/aerospace-marks/internal/aerospace"
	"github.com/cristianoliveira/aerospace-marks/internal/format"
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
	newRootCmd.AddCommand(InfoCmd(storage, aerospaceClient))

	// Manage marks
	newRootCmd.AddCommand(MarkCmd(storage, aerospaceClient))
	newRootCmd.AddCommand(UnmarkCmd(storage))

	// Manage windows with marks
	newRootCmd.AddCommand(enableOutputFlag(FocusCmd(storage, aerospaceClient)))
	newRootCmd.AddCommand(enableOutputFlag(ListCmd(storage, aerospaceClient)))
	newRootCmd.AddCommand(enableOutputFlag(SummonCmd(storage, aerospaceClient)))
	newRootCmd.AddCommand(enableOutputFlag(GetCmd(storage, aerospaceClient)))

	return newRootCmd
}

//nolint:gochecknoinits // init function is required for cobra root command setup
func init() {
	// NOTE: add here global flags
	// rootCmd.Flags().BoolP("version", "v", false, "Print version information")
}

func enableOutputFlag(command *cobra.Command) *cobra.Command {
	// Add output flag
	command.Flags().StringP("output", "o", "text", "Output format: text, json, or csv")
	command.Flag("output").DefValue = string(format.OutputFormatText)
	return command
}

func Run(storage storage.MarkStorage, aerospaceClient aerospace.AerosSpaceMarkWindows) {
	rootCmd := NewRootCmd(storage, aerospaceClient)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// VERSION is the version of the application.
// THIS IS GENERATED DON'T EDIT
// NOTE: to update VERSION to empty string
// and then run scripts/validate-version.sh.
//
//nolint:gochecknoglobals // VERSION is a build-time constant
var VERSION = "v0.3.0"
