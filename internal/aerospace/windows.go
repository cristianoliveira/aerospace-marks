package aerospace

import (
	exec "os/exec"
	"strings"
)

// GetFocusedWindowID returns the ID of the currently focused window
func GetFocusedWindowID() (string, error) {
	out, err := exec.Command("aerospace", "list-windows", "--focused").Output()
	if err != nil {
		return "", err
	}
	windowID := strings.Fields(string(out))[0]
	return windowID, nil
}
