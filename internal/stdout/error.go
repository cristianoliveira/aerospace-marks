package stdout

import (
	"fmt"
	"os"
)

var ShuldExit = true

// ErrorAndExit is a function that prints an error message to stderr and exits the program with a non-zero status code.
func ErrorAndExit(err error) {
	if err != nil {
		errorMessage := fmt.Errorf("error: %v", err)
		fmt.Fprintln(os.Stderr, err.Error())
		if ShuldExit {
			os.Exit(1)
		}

		fmt.Println(errorMessage)
	}
}

func ErrorAndExitf(format string, a ...any) {
	fmt.Fprintln(os.Stderr, fmt.Errorf(format, a...))
	if ShuldExit {
		os.Exit(1)
	}
}
