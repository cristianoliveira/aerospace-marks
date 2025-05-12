package aerospace

import (
	"fmt"
	exec "os/exec"
	"strings"
)

// GetFocusedWindowID returns the ID of the currently focused window
func GetFocusedWindowID() (string, error) {
	path, err := exec.LookPath("aerospace")
	if err != nil {
			return "", fmt.Errorf("aerospace binary not found: %w", err)
	}
	fmt.Printf("@@@@@@@@ path %+v \n", path)
	out, err := exec.Command("aerospace", "list-windows", "--focused").Output()
	if err != nil {
		return "", err
	}
	windowID := strings.Fields(string(out))[0]
	return windowID, nil
}

// GetWindowByID returns the window information for a given window ID
func GetWindowByID(windowID string) (string, error) {
	path, err := exec.LookPath("aerospace")
	if err != nil {
			return "", fmt.Errorf("aerospace binary not found: %w", err)
	}
	fmt.Printf("@@@@@@@@ path %+v \n", path)
	out, err := exec.Command(path, "list-windows", "--all").Output()
	if err != nil {
		return "", err
	}
	// Split the output into lines
	lines := strings.Split(string(out), "\n")
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
