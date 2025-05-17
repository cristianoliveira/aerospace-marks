package aerospacecli

import (
	"encoding/json"
	"fmt"
)

/**
Example:
  {
    "window-id" : 7984,
    "window-title" : "WhatsApp",
    "app-name" : "WhatsApp"
  }
*/
type Window struct {
	WindowID    int    `json:"window-id"`
	WindowTitle string `json:"window-title"`
	AppName     string `json:"app-name"`
}

// GetAllWindows returns all windows
func (c *AeroSpaceDefaultConnection) GetAllWindows() ([]Window, error) {
	response, err := c.Conn.SendCommand("list-windows", []string{"--json", "--all"})
	if err != nil {
		return nil, err
	}
	var windows []Window
	err = json.Unmarshal([]byte(response.StdOut), &windows)
	if err != nil {
		return nil, err
	}
	return windows, nil
}

func (c *AeroSpaceDefaultConnection) GetWindowByID(windowID int) (Window, error) {
	windows, err := c.GetAllWindows()
	if err != nil {
		return Window{}, err
	}
	for _, window := range windows {
		if window.WindowID == windowID {
			return window, nil
		}
	}
	return Window{}, fmt.Errorf("window with ID %d not found", windowID)
}

