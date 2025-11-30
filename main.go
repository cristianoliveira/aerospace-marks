package main

import (
	"github.com/cristianoliveira/aerospace-marks/cmd"
	"github.com/cristianoliveira/aerospace-marks/internal/aerospace"
	"github.com/cristianoliveira/aerospace-marks/internal/logger"
	"github.com/cristianoliveira/aerospace-marks/internal/stdout"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
)

func main() {
	defaultLogger, err := logger.NewLogger()
	if err != nil {
		stdout.ErrorAndExit(err)
		return
	}
	defer func() {
		if closeErr := defaultLogger.Close(); closeErr != nil {
			stdout.ErrorAndExit(closeErr)
		}
	}()
	logger.SetDefaultLogger(defaultLogger)
	defaultLogger.LogInfo("Starting Aerospace Marks CLI")

	connector := storage.MarksDatabaseConnector{}
	conn, err := connector.Connect()
	if err != nil {
		stdout.ErrorAndExit(err)
	}
	markClient, err := storage.NewMarkClient(conn)
	if err != nil {
		stdout.ErrorAndExit(err)
	}
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			stdout.ErrorAndExit(closeErr)
		}
	}()

	aerospaceMarkClient, err := aerospace.NewAeroSpaceClient()
	if err != nil {
		stdout.ErrorAndExit(err)
	}

	cmd.Run(markClient, aerospaceMarkClient)
}
