/*
Copyright © 2025 Cristian Oliveira license@cristianoliveira.dev
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/cristianoliveira/aerospace-marks/internal/aerospace"
	"github.com/cristianoliveira/aerospace-marks/internal/cli"
	"github.com/cristianoliveira/aerospace-marks/internal/format"
	"github.com/cristianoliveira/aerospace-marks/internal/logger"
	"github.com/cristianoliveira/aerospace-marks/internal/stdout"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
	"github.com/spf13/cobra"
)

// formatSingleFieldOutput formats a single field output based on output format.
// For json/csv formats, uses OutputEvent. For text/default, outputs plain value.
//
//nolint:unparam // fieldName is used conditionally for app_name field
func formatSingleFieldOutput(
	outputFormat string,
	command string,
	windowID int,
	fieldValue string,
	fieldName string,
) error {
	// text format or default - output plain value (backward compatibility)
	if outputFormat == string(format.OutputFormatText) || outputFormat == "" {
		fmt.Fprint(os.Stdout, fieldValue)
		return nil
	}

	// json/csv formats - use OutputEvent
	formatter, formatErr := format.NewOutputEventFormatter(os.Stdout, outputFormat)
	if formatErr != nil {
		return formatErr
	}
	event := format.OutputEvent{
		Command:  command,
		WindowID: windowID,
		Result:   fieldValue,
		Message:  fieldValue,
	}
	if fieldName == "app_name" {
		event.AppName = fieldValue
	}
	if fmtErr := formatter.Format(event); fmtErr != nil {
		return fmt.Errorf("failed to format output: %w", fmtErr)
	}
	return nil
}

//nolint:gocognit,funlen // GetCmd has high complexity and length due to multiple field output options
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
		Args: cobra.MatchAll(
			cobra.ExactArgs(1),
			cli.ValidateArgIsNotEmpty,
		),
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.GetDefaultLogger()

			mark := args[0]
			logger.LogDebug("Getting window by mark: %s", mark)

			// Get and validate output format early
			outputFormat, err := cmd.Flags().GetString("output")
			if err != nil {
				stdout.ErrorAndExit(fmt.Errorf("failed to get output flag: %w", err))
				return
			}

			// Default to text if not specified
			if outputFormat == "" {
				outputFormat = string(format.OutputFormatText)
			}

			markedWindow, err := storageClient.GetWindowByMark(mark)
			if err != nil {
				stdout.ErrorAndExit(err)
				return
			}

			windowID := markedWindow.WindowID
			if windowID == 0 {
				stdout.ErrorAndExit(fmt.Errorf("no window found for mark %s", mark))
				return
			}

			getWinID, _ := cmd.Flags().GetBool("window-id")
			if getWinID {
				// Single field flag - respect --output flag for json/csv, plain value for text/default
				if formatErr := formatSingleFieldOutput(
					outputFormat,
					"get",
					windowID,
					strconv.Itoa(windowID),
					"",
				); formatErr != nil {
					stdout.ErrorAndExit(formatErr)
					return
				}
				return
			}

			window, err := aerospaceClient.GetWindowByID(windowID)
			if err != nil {
				stdout.ErrorAndExit(err)
				return
			}
			if window == nil {
				stdout.ErrorAndExit(fmt.Errorf("no window found for ID %d", windowID))
				return
			}
			logger.LogDebug(
				"Get window by ID",
				"windowID", windowID,
				"window", *window,
				"windowTitle", window.WindowTitle,
				"err", err,
			)

			getWinTitle, _ := cmd.Flags().GetBool("window-title")
			if getWinTitle {
				// Single field flag - respect --output flag for json/csv, plain value for text/default
				if formatErr := formatSingleFieldOutput(
					outputFormat,
					"get",
					windowID,
					window.WindowTitle,
					"",
				); formatErr != nil {
					stdout.ErrorAndExit(formatErr)
					return
				}
				return
			}

			getWinApp, _ := cmd.Flags().GetBool("app-name")
			if getWinApp {
				// Single field flag - respect --output flag for json/csv, plain value for text/default
				if formatErr := formatSingleFieldOutput(
					outputFormat,
					"get",
					windowID,
					window.AppName,
					"app_name",
				); formatErr != nil {
					stdout.ErrorAndExit(formatErr)
					return
				}
				return
			}

			getWinAppBundleID, _ := cmd.Flags().GetBool("app-bundle-id")
			if getWinAppBundleID {
				// Single field flag - respect --output flag for json/csv, plain value for text/default
				if formatErr := formatSingleFieldOutput(
					outputFormat,
					"get",
					windowID,
					window.AppBundleID,
					"",
				); formatErr != nil {
					stdout.ErrorAndExit(formatErr)
					return
				}
				return
			}

			// Full window output - use OutputEvent
			formatter, err := format.NewOutputEventFormatter(os.Stdout, outputFormat)
			if err != nil {
				stdout.ErrorAndExit(err)
				return
			}

			event := format.OutputEvent{
				Command:   "get",
				WindowID:  window.WindowID,
				AppName:   window.AppName,
				Workspace: window.Workspace,
				Message:   window.WindowTitle, // For backward compatibility, window_title goes in Message
				Result: fmt.Sprintf(
					"%d | %s | %s",
					window.WindowID,
					window.AppName,
					window.WindowTitle,
				),
			}

			logger.LogInfo(
				"Printing full window info",
				"windowID", windowID,
				"window", window,
			)

			if formatErr := formatter.Format(event); formatErr != nil {
				stdout.ErrorAndExit(fmt.Errorf("failed to format output: %w", formatErr))
				return
			}
		},
	}

	getCmd.Flags().BoolP("window-id", "i", false, "Get only window [i]D")
	getCmd.Flags().BoolP("window-title", "t", false, "Get only window [t]itle")
	getCmd.Flags().BoolP("app-name", "a", false, "Get only window [a]pp name")
	getCmd.Flags().BoolP("app-bundle-id", "b", false, "Get only window app [b]undle ID")

	return getCmd
}
