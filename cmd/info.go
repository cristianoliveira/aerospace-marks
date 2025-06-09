/*
Copyright © 2025 Cristian Oliveira license@cristianoliveira.dev
*/
package cmd

import (
	"fmt"

	"github.com/cristianoliveira/aerospace-marks/internal/aerospace"
	"github.com/cristianoliveira/aerospace-marks/internal/constants"
	"github.com/cristianoliveira/aerospace-marks/internal/logger"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
func InfoCmd(
	storageClient storage.MarkStorage,
	aerospaceClient aerospace.AerosSpaceMarkWindows,
) *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Displays aerospace-marks config information",
		Long: `Displays the config information of aerospace-marks.

This command allows you to view the current configurations for the aerospace-marks CLI.
It also displays help information about environment variables available.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logConfig := logger.GetDefaultLogger().GetConfig()
			dbConfig := storageClient.Client().GetStorageConfig()
			client := aerospaceClient.Client().Connection()
			socketPath, err := client.GetSocketPath()
			if err != nil {
				return fmt.Errorf("failed to get socket path: %w", err)
			}

			var validationInfo string
			if err = client.CheckServerVersion(); err != nil {
				validationInfo = "Incompatible. Reason: " + err.Error()
			} else {
				validationInfo = "Compatible."
			}
			serverVersion, err := client.GetServerVersion()
			if err != nil {
				return fmt.Errorf("failed to get server version: %w", err)
			}

			fmt.Printf(`Aerospace Marks CLI - Configuration

[Socket]
Path: %s
Version: %s
Status: %s

[Database]
Name: %s
Path: %s

[Logging]
Path: %s
Level: %s

Configure with ENV variables:
%s - Path to the socket file.
%s - Path to database directory.
%s - Log level [debug|info|warn|error] (default: disabled)
%s - Path to the logs file.
`,
				socketPath,
				serverVersion,
				validationInfo,

				// Database configuration
				dbConfig.DbName,
				dbConfig.DbPath,

				// logging configuration
				logConfig.Path,
				logConfig.Level,

				// Environment variables
				constants.EnvAeroSpaceSock,
				constants.EnvAeroSpaceMarksDbPath,
				constants.EnvAeroSpaceMarksLogsLevel,
				constants.EnvAeroSpaceMarksLogsPath,
			)

			return nil
		},
	}
}
