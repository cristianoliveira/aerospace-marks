package aerospace

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cristianoliveira/aerospace-marks/pkgs/aerospacecli"
)

// GetFocusedWindowID returns the ID of the currently focused window
func GetFocusedWindowID() (*aerospacecli.Window, error) {
	cli, err := aerospacecli.NewAeroSpaceConnection()
	if err != nil {
		return nil, err
	}
	defer cli.Conn.CloseConnection()

	response, err := cli.Conn.SendCommand("list-windows", []string{"--focused", "--json"})
	if err != nil {
		return nil, err
	}

	windows := []aerospacecli.Window{}
	err = json.Unmarshal([]byte(response.StdOut), &windows)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal windows\n%w", err)
	}

	if len(windows) == 0 {
		return nil, fmt.Errorf("no focused windows found")
	}

	return &windows[0], nil
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
func GetAllWindows() ([]aerospacecli.Window, error) {
	cli, err := aerospacecli.NewAeroSpaceConnection()
	if err != nil {
		return nil, err
	}
	defer cli.Conn.CloseConnection()
	// FIXME: use --json and return a struct instead
	res, err := cli.Conn.SendCommand("list-windows", []string{"--all", "--json"})
	if err != nil {
		return nil, err
	}

	var windows []aerospacecli.Window
	err = json.Unmarshal([]byte(res.StdOut), &windows)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal windows\n%w", err)
	}

	return windows, nil
}
