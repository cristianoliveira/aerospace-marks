package aerospace

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cristianoliveira/aerospace-marks/pkgs/aerospacecli"
)

// GetFocusedWindowID returns the ID of the currently focused window
func GetFocusedWindowID() (string, error) {
	cli, err := aerospacecli.NewAeroSpaceConnection()
	if err != nil {
		return "", err
	}
	defer cli.Conn.CloseConnection()

	response, err := cli.Conn.SendCommand("list-windows", []string{"--focused"})
	if err != nil {
		return "", err
	}
	windowID := strings.Fields(response.StdOut)[0]
	return windowID, nil
}

// GetWindowByID returns the window information for a given window ID
func GetWindowByID(windowID string) (string, error) {
	intWindowID, err := strconv.Atoi(windowID)
	cli, err := aerospacecli.NewAeroSpaceConnection()
	if err != nil {
		return "", err
	}
	defer cli.Conn.CloseConnection()

	window, err := cli.GetWindowByID(intWindowID)
	if err != nil {
		return "", err
	}

	windowInfo := fmt.Sprintf("%d | %s | %s", window.WindowID, window.AppName, window.WindowTitle)
	return windowInfo, nil
}

// SetFocusToWindowId sets the focus to a window by id
func SetFocusToWindowId(windowID string) error {
	cli, err := aerospacecli.NewAeroSpaceConnection()
	if err != nil {
		return err
	}
	defer cli.Conn.CloseConnection()

	response, err := cli.Conn.SendCommand("focus", []string{"--window-id", windowID})
	if err != nil {
		return err
	}
	if response.ExitCode != 0 {
		return fmt.Errorf("failed to focus window: %s", response.StdErr)
	}
	return nil
}

// GetAllWindows returns all windows
func GetAllWindows() ([]string, error) {
	cli, err := aerospacecli.NewAeroSpaceConnection()
	if err != nil {
		return nil, err
	}
	defer cli.Conn.CloseConnection()
	// FIXME: use --json and return a struct instead
	response, err := cli.Conn.SendCommand("list-windows", []string{"--all"})
	if err != nil {
		return nil, err
	}
	// Split the output into lines
	lines := strings.Split(response.StdOut, "\n")
	return lines, nil
}
