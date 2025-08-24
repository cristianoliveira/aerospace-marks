package testutils

// This module contains test utilities for CLI commands.
// - Shell output
// - Cobra command execution
// - Standard input/output capturing

import (
	"bytes"
	"fmt"
	"io"
	"os"

	aerospacecli "github.com/cristianoliveira/aerospace-ipc"
	"github.com/spf13/cobra"
)

// CmdExecuteWithStdin executes a Cobra command with the provided arguments and stdin input.
func CmdExecuteWithStdin(cmd *cobra.Command, stdinInput string, args ...string) (string, error) {
	cmd.SetArgs(args)
	cmd.SetIn(bytes.NewBufferString(stdinInput))
	stdOut, err := CaptureStdOut(func() error {
		return cmd.Execute()
	})

	if err != nil {
		return "", err
	}

	return string(stdOut), nil
}

func CmdExecute(cmd *cobra.Command, args ...string) (string, error) {
	cmd.SetArgs(args)
	stdOut, err := CaptureStdOut(func() error {
		return cmd.Execute()
	})

	if err != nil {
		return "", err
	}

	return string(stdOut), nil
}

func CaptureStdOut(f func() error) (string, error) {
	old := os.Stdout
	defer func() {
		os.Stdout = old
	}()
	errOld := os.Stderr
	defer func() {
		os.Stdout = errOld
	}()

	r, w, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = wErr // Redirect stderr as well

	// Run the function that prints to stdout
	err := f()
	if err != nil {
		return "", err
	}

	// Close writer and restore stdout
	err = w.Close()
	if err != nil {
		return "", err
	}

	err = wErr.Close()
	if err != nil {
		return "", err
	}

	// Read output
	stdString, err := readFileToString(r)
	if err != nil {
		return "", err
	}
	stdErrString, err := readFileToString(rErr)
	if err != nil {
		return "", err
	}

	if stdErrString != "" {
		return "", fmt.Errorf("error: %s", stdErrString)
	}

	return stdString, nil
}

func readFileToString(f *os.File) (string, error) {
	bytes, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

type MockEmptyAerspaceMarkWindows struct{}

func (d *MockEmptyAerspaceMarkWindows) Client() *aerospacecli.AeroSpaceWM {
	return &aerospacecli.AeroSpaceWM{}
}

func (d *MockEmptyAerspaceMarkWindows) GetWindowByID(windowID int) (*aerospacecli.Window, error) {
	fmt.Println("Mocked GetWindowByID called with windowID:", windowID)
	return &aerospacecli.Window{}, nil
}
