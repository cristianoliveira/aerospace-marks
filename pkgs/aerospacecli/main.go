package aerospacecli

import (
	"fmt"
)

type AeroSpaceClient interface {
	// Windows Methods
	
	// GetAllWindows returns all windows
	// 
	// Returns all windows from AeroSpaceWM
	// Same as `aerospace list-windows --all --json`
	GetAllWindows() ([]Window, error)

	// GetFocusedWindow returns the focused window
	//
	// Returns the focused window from AeroSpaceWM
	// Same as `aerospace list-windows --focused --json`
	GetFocusedWindow() (*Window, error)

	// SetFocusByWindowID sets the focused window
	//
	// Sets the focused window from AeroSpaceWM
	// Same as `aerospace focus --window-id <window-id>`
	SetFocusByWindowID(windowID int) error

	// GetFocusedWorkspace returns the current workspace
	//
	// Returns the current workspace from AeroSpaceWM
	// Same as:
	//
	// aerospace list-workspaces --focused --json
	GetFocusedWorkspace() (*Workspace, error)


	// MoveWindowToWorkspace moves the window to the workspace
	//
	// Moves the window to the workspace from AeroSpaceWM
	// Similar to:
	//
	// aerospace move-node-to-workspace <workspace> --window-id <window-id>
	MoveWindowToWorkspace(windowID int, workspaceName string) error

	// Connection Methods
	
	// CloseConnection
	// Closes the AeroSpaceWM connection and releases the resources
	CloseConnection() error
}

type AeroSpaceWM struct {
	MinAerospaceVersion string
	Conn                AeroSpaceSocketConn
}

func (a *AeroSpaceWM) CloseConnection() error {
	return a.Conn.CloseConnection()
}

// NewAeroSpaceClient creates a new AeroSpaceClient with the default socket path.
// It checks for environment variable AEROSPACESOCK or uses the default socket path.
// which is usually /tmp/bobko.aerospace-<username>.sock
func NewAeroSpaceConnection() (*AeroSpaceWM, error) {
	conn, err := DefaultConnector.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to socket\n %w", err)
	}

	client := &AeroSpaceWM{
		MinAerospaceVersion: "0.15.2-Beta",
		Conn:                conn,
	}

	return client, nil
}
