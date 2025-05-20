/*
Copyright © 2025 Cristian Oliveira licence@cristianoliveira.dev
*/
package cmd

import (
	"fmt"

	"github.com/cristianoliveira/aerospace-marks/internal/constants"
	"github.com/cristianoliveira/aerospace-marks/internal/logger"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
	"github.com/cristianoliveira/aerospace-marks/pkgs/aerospacecli"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
func ConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Displays aerospace-marks configuration",
		Long: `Displays the configuration of aerospace-marks.

This command allows you to view the current configurations for the aerospace-marks CLI.
It also displays help information about environment variables available.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logConfig := logger.GetDefaultLogger().GetConfig()
			dbConfig := storage.GetDatabaseConfig()
			socketPath, err := aerospacecli.GetSocketPath()
			if err != nil {
				return fmt.Errorf("failed to get socket path: %w", err)
			}

			cmd.Println(fmt.Sprintf(`Aerospace Marks CLI - Configuration

[Socket]
Path: %s

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
				dbConfig.DbName,
				dbConfig.DbPath,
				logConfig.Path,
				logConfig.Level,
			  aerospacecli.EnvAeroSpaceSock,
				constants.EnvAeroSpaceMarksDbPath,
				constants.EnvAeroSpaceMarksLogsLevel,
				constants.EnvAeroSpaceMarksLogsPath,
			))

			return nil
		},
	}
}
