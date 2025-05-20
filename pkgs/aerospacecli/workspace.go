package aerospacecli

import (
	"encoding/json"
	"fmt"
)

/**
 * Workspace represents the JSON structure for workspace commands.
 * Example:
[
  {
    "workspace" : "9"
  },
  {
    "workspace" : "scratchpad"
  }
]
*/
type Workspace struct {
	Workspace string    `json:"workspace"`
}

func (a *AeroSpaceWM) GetFocusedWorkspace() (*Workspace, error) {
	response, err := a.Conn.SendCommand("list-workspaces", []string{"--focused", "--json"})
	if err != nil {
		return nil, err
	}

	var workspaces []Workspace
	err = json.Unmarshal([]byte(response.StdOut), &workspaces)
	if err != nil {
		return nil, err
	}
	if len(workspaces) == 0 {
		return nil, fmt.Errorf("no workspace focused found")
	}

	return &workspaces[0], nil
}

func (a *AeroSpaceWM) MoveWindowToWorkspace(windowID int, workspaceName string) error {
	response, err := a.Conn.SendCommand("move-node-to-workspace", []string{workspaceName, "--window-id", fmt.Sprintf("%d", windowID)})
	if err != nil {
		return err
	}

	if response.ExitCode != 0 {
		return fmt.Errorf("failed to move window to workspace: %s", response.StdErr)
	}

	return nil
}
