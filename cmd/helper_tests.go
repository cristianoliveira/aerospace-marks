package cmd

import (
	"bytes"
	"io"
)

func cmdExecute(args... string) (string, error) {
	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	if err != nil {
		return "", err
	}
	out, err := io.ReadAll(b)

	if err != nil {
		return "", err
	}

	return string(out), nil
}
