package main

import (
	"github.com/cristianoliveira/aerospace-marks/cmd"
	"github.com/cristianoliveira/aerospace-marks/internal/logger"
)

func main() {
	defaultLogger, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	defer defaultLogger.Close()
	logger.SetDefaultLogger(defaultLogger)
	defaultLogger.LogInfo("Starting Aerospace Marks CLI")

	cmd.Run()
}
