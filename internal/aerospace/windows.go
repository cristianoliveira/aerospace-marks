package aerospace

import (
	"fmt"
	"strconv"

	aerospacecli "github.com/cristianoliveira/aerospace-ipc"
)

type AerosSpaceMarkWindows interface {
	// GetFocusedWindowID returns the ID of the currently focused window
	//
	// Returns the window ID of the currently focused window
	// or an error if the window ID is not found
	GetWindowByID(windowID string) (*aerospacecli.Window, error)

	// Client returns the AeroSpaceWM client
	//
	// Returns the AeroSpaceWM client
	// or panics if the client is not initialized
	Client() *aerospacecli.AeroSpaceWM
}

type DefaultAeroSpaceWindows struct {
	client *aerospacecli.AeroSpaceWM
}

func (d *DefaultAeroSpaceWindows) Client() *aerospacecli.AeroSpaceWM {
	if d.client == nil {
		panic("ASSERTION: AeroSpaceWM client is not initialized")
	}

	return d.client
}

func (d *DefaultAeroSpaceWindows) GetWindowByID(windowID string) (*aerospacecli.Window, error) {
	windows, err := d.client.GetAllWindows()
	if err != nil {
		return nil, err
	}

	intWindowID, err := strconv.Atoi(windowID)
	if err != nil {
		return nil, fmt.Errorf("invalid window ID '%s': %w", windowID, err)
	}

	for _, window := range windows {
		if window.WindowID == intWindowID {
			return &window, nil
		}
	}

	return nil, fmt.Errorf("window with ID %s not found", windowID)
}

func NewAeroSpaceClient() (*DefaultAeroSpaceWindows, error) {
	cli, err := aerospacecli.NewAeroSpaceClient()
	if err != nil {
		return nil, err
	}

	return &DefaultAeroSpaceWindows{
		client: cli,
	}, nil
}
