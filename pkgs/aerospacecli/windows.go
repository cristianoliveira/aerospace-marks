package aerospacecli

import (
	"encoding/json"
	"fmt"
)

/*
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

func (w Window) String() string {
	return fmt.Sprintf("%d | %s | %s", w.WindowID, w.AppName, w.WindowTitle)
}

func (c *AeroSpaceWM) GetAllWindows() ([]Window, error) {
	response, err := c.Conn.SendCommand("list-windows", []string{"--all", "--json"})
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

func (c *AeroSpaceWM) GetAllWindowsByWorkspace(workspaceName string) ([]Window, error) {
	response, err := c.Conn.SendCommand("list-windows", []string{"--workspace", workspaceName, "--json"})
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

func (c *AeroSpaceWM) GetFocusedWindow() (*Window, error) {
	response, err := c.Conn.SendCommand("list-windows", []string{"--focused", "--json"})
	if err != nil {
		return nil, err
	}

	var windows []Window
	err = json.Unmarshal([]byte(response.StdOut), &windows)
	if err != nil {
		return nil, err
	}
	if len(windows) == 0 {
		return nil, fmt.Errorf("no windows focused found")
	}

	return &windows[0], nil
}

func (c *AeroSpaceWM) SetFocusByWindowID(windowID int) error {
	response, err := c.Conn.SendCommand("focus", []string{"--window-id", fmt.Sprintf("%d", windowID)})
	if err != nil {
		return err
	}

	if response.ExitCode != 0 {
		return fmt.Errorf("failed to focus window with ID %d\n%s", windowID, response.StdErr)
	}

	return nil
}
