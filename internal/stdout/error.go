package stdout

import (
	"fmt"
	"os"
)

var ShuldExit = true

// ErrorAndExit is a function that prints an error message to stderr and exits the program with a non-zero status code.
func ErrorAndExit(err error) error {
	if err != nil {
		errorMessage := fmt.Errorf("Error: %v", err)
		fmt.Fprintln(os.Stderr, errorMessage)
		if ShuldExit {
			os.Exit(1)
		}

		fmt.Println(errorMessage)
	}

	return nil
}

func ErrorAndExitf(format string, a ...any) error {
	errorMessage := fmt.Errorf("Error: %v", fmt.Errorf(format, a...))
	fmt.Fprintln(os.Stderr, errorMessage)
	if ShuldExit {
		os.Exit(1)
	}

	return nil
}
