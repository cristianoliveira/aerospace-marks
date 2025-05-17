package cmd

import (
	"bytes"
	"io"
	"os"
)

func cmdExecute(args... string) (string, error) {
	rootCmd.SetArgs(args)
	stdOut, err := captureStdOut(func() error {
		return rootCmd.Execute()
	})

	if err != nil {
		return "", err
	}

	return string(stdOut), nil
}

func captureStdOut(f func() error) (string, error) {
	var buf bytes.Buffer
	// Save original stdout
	old := os.Stdout
	// Redirect stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the function that prints to stdout
	err := f()
	if err != nil {
		return "", err
	}

	// Close writer and restore stdout
	w.Close()
	os.Stdout = old

	// Read output
	io.Copy(&buf, r)
	return buf.String(), nil
}

