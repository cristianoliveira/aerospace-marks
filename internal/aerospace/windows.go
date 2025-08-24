package aerospace

import (
	"fmt"
	"log"

	aerospacecli "github.com/cristianoliveira/aerospace-ipc"
	"github.com/cristianoliveira/aerospace-marks/internal/logger"
)

type AerosSpaceMarkWindows interface {
	// GetFocusedWindowID returns the ID of the currently focused window
	//
	// Returns the window ID of the currently focused window
	// or an error if the window ID is not found
	GetWindowByID(windowID int) (*aerospacecli.Window, error)

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
		log.Fatal("ASSERT: AeroSpaceWM client is not initialized")
	}

	return d.client
}

func (d *DefaultAeroSpaceWindows) GetWindowByID(windowID int) (*aerospacecli.Window, error) {
	logger := logger.GetDefaultLogger()
	windows, err := d.client.GetAllWindows()
	if err != nil {
		return nil, err
	}
	logger.LogDebug("Windows found: %d", len(windows))

	for _, window := range windows {
		if window.WindowID == windowID {
			return &window, nil
		}
	}

	return nil, fmt.Errorf("window with ID %d not found", windowID)
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
