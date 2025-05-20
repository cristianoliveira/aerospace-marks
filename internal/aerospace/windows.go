package aerospace

import (
	"fmt"
	"strconv"

	"github.com/cristianoliveira/aerospace-marks/pkgs/aerospacecli"
)

type AerosSpaceWindows interface {
	// GetFocusedWindowID returns the ID of the currently focused window
	//
	// Returns the window ID of the currently focused window
	// or an error if the window ID is not found
	GetWindowByID(windowID int) (*aerospacecli.Window, error)
}

type DefaultAeroSpaceWindows struct {
	cli *aerospacecli.AeroSpaceWM
}

// GetFocusedWindowID returns the ID of the currently focused window
func GetFocusedWindowID() (*aerospacecli.Window, error) {
	cli, err := aerospacecli.NewAeroSpaceConnection()
	if err != nil {
		return nil, err
	}
	defer cli.Conn.CloseConnection()

	return cli.GetFocusedWindow()
}

// GetWindowByID returns the window information for a given window ID
func GetWindowByID(windowID string) (*aerospacecli.Window, error) {
	intWindowID, err := strconv.Atoi(windowID)
	cli, err := aerospacecli.NewAeroSpaceConnection()
	if err != nil {
		return nil, err
	}
	defer cli.Conn.CloseConnection()

	windows, err := cli.GetAllWindows()
	if err != nil {
		return nil, err
	}

	for _, window := range windows {
		if window.WindowID == intWindowID {
			return &window, nil
		}
	}

	return nil, fmt.Errorf("window with ID %d not found", intWindowID)
}

// SetFocusToWindowId sets the focus to a window by id
func SetFocusToWindowId(windowID string) error {
	intWindowID, err := strconv.Atoi(windowID)
	cli, err := aerospacecli.NewAeroSpaceConnection()
	if err != nil {
		return err
	}
	defer cli.Conn.CloseConnection()

	return cli.SetFocusByWindowID(intWindowID)
}

// GetAllWindows returns all windows
func GetAllWindows() ([]aerospacecli.Window, error) {
	cli, err := aerospacecli.NewAeroSpaceConnection()
	if err != nil {
		return nil, err
	}
	defer cli.Conn.CloseConnection()

	return cli.GetAllWindows()
}
