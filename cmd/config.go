/*
Copyright © 2025 Cristian Oliveira licence@cristianoliveira.dev
*/
package cmd

import (
	"fmt"

	"github.com/cristianoliveira/aerospace-marks/internal/logger"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
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

			cmd.Println(fmt.Sprintf(`Aerospace Marks CLI - Configuration

[Database]
Name: %s
Path: %s

[Logging]
Path: %s
Level: %s

Configure with ENV variables:
AEROSPACE_MARKS_DB_PATH    - Path to database directory.
AEROSPACE_MARKS_LOGS_LEVEL - Log level [debug|info|warn|error] (default: disabled)
AEROSPACE_MARKS_LOGS_PATH  - Path to the logs file.
			`,
				dbConfig.DbName,
				dbConfig.DbPath,
				logConfig.Path,
				logConfig.Level,
			))

			return nil
		},
	}
}
