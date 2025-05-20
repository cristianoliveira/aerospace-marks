package aerospacecli

import (
	"fmt"
	"os"
)

const (
  // Environment variables

  // EnvAeroSpaceSock is the environment variable for the AeroSpace socket path
  // default: `/tmp/bobko.aerospace-$USER.sock`
  EnvAeroSpaceSock string = "AEROSPACESOCK"
)

// GetSocketPath returns the socket path
// It checks for environment variable AEROSPACESOCK or uses the default socket path
// which is usually /tmp/bobko.aerospace-<username>.sock
// See: https://github.com/nikitabobko/AeroSpace/blob/f12ee6c9d914f7b561ff7d5c64909882c67061cd/Sources/Cli/_main.swift#L47
func GetSocketPath() (string, error) {
	socketPath := fmt.Sprintf("/tmp/bobko.%s-%s.sock", "aerospace", os.Getenv("USER"))
	if os.Getenv(EnvAeroSpaceSock) != "" {
		socketPath = os.Getenv(EnvAeroSpaceSock)
	} else {
		socketPath = fmt.Sprintf("/tmp/bobko.%s-%s.sock", "aerospace", os.Getenv("USER"))
	}

	if _, err := os.Stat(socketPath); os.IsNotExist(err) {
		return "", fmt.Errorf("failed to access socket path %s\r reason: %w", socketPath, err)
	}

	return socketPath, nil
}
