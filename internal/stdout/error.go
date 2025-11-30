package stdout

import (
	"fmt"
	"os"

	"github.com/cristianoliveira/aerospace-marks/internal/logger"
)

//nolint:gochecknoglobals // ShouldExit is a test configuration flag
var ShouldExit = true

// ErrorAndExit is a function that prints an error message to stderr and exits the program with a non-zero status code.
func ErrorAndExit(err error) {
	if err != nil {
		logger := logger.GetDefaultLogger()
		logger.LogError("ERROR:", "msg", err)
		errorMessage := fmt.Errorf("error: %w", err)
		fmt.Fprintln(os.Stderr, err.Error())
		if ShouldExit {
			os.Exit(1)
		}

		fmt.Fprintln(os.Stdout, errorMessage)
	}
}

func ErrorAndExitf(format string, a ...any) {
	logger := logger.GetDefaultLogger()
	logger.LogError("ERROR", "msg", fmt.Errorf(format, a...))
	fmt.Fprintln(os.Stderr, fmt.Errorf(format, a...))
	if ShouldExit {
		os.Exit(1)
	}
}
