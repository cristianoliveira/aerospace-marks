package aerospace

import (
	"github.com/cristianoliveira/aerospace-marks/pkgs/aerospacecli"
	"strings"
	"fmt"
)

// GetFocusedWindowID returns the ID of the currently focused window
func GetFocusedWindowID() (string, error) {
	response, err := aerospacecli.SendCommand("list-windows", []string{"--focused"})
	if err != nil {
		return "", err
	}
	windowID := strings.Fields(response.StdOut)[0]
	return windowID, nil
}

// GetWindowByID returns the window information for a given window ID
func GetWindowByID(windowID string) (string, error) {
	response, err := aerospacecli.SendCommand("list-windows", []string{"--all"})
	if err != nil {
		return "", err
	}
	// Split the output into lines
	lines := strings.Split(response.StdOut, "\n")
	// Find the line that contains the window ID
	var windowInfo string
	for _, line := range lines {
		id := strings.Split(line, "|")[0]
		id = strings.TrimSpace(id)

		if id == windowID {
			windowInfo = line
			break
		}
	}

	if windowInfo == "" {
		return "", fmt.Errorf("window with ID %s not found", windowID)
	}

	return windowInfo, nil
}

// SetFocusToWindowId sets the focus to a window by id
func SetFocusToWindowId(windowID string) error {
	response, err := aerospacecli.SendCommand("focus", []string{"--window-id", windowID})
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
	// FIXME: use --json and return a struct instead
	response, err := aerospacecli.SendCommand("list-windows", []string{"--all"})
	if err != nil {
		return nil, err
	}
	// Split the output into lines
	lines := strings.Split(response.StdOut, "\n")
	return lines, nil
}
