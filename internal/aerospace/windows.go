package aerospace

import (
	"fmt"

	"github.com/cristianoliveira/aerospace-marks/internal/logger"

	aerospacecli "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace"
	"github.com/cristianoliveira/aerospace-ipc/pkg/aerospace/windows"
)

type AerosSpaceMarkWindows interface {
	// GetFocusedWindowID returns the ID of the currently focused window
	//
	// Returns the window ID of the currently focused window
	// or an error if the window ID is not found
	GetWindowByID(windowID int) (*windows.Window, error)

	// Client returns the AeroSpaceWM client
	//
	// Returns the AeroSpaceWM client
	// or panics if the client is not initialized
	Client() *aerospacecli.AeroSpaceWM
}

type DefaultAeroSpaceWindows struct {
	client *aerospacecli.AeroSpaceWM
}

func NewAeroSpaceClient() (*DefaultAeroSpaceWindows, error) {
	cli, err := aerospacecli.NewClient()
	if err != nil {
		return nil, err
	}

	return &DefaultAeroSpaceWindows{
		client: cli,
	}, nil
}

func (d *DefaultAeroSpaceWindows) Client() *aerospacecli.AeroSpaceWM {
	if d.client == nil {
		logger := logger.GetDefaultLogger()
		logger.LogError("ASSERT: AeroSpaceWM client is not initialized", nil)
		panic("AeroSpaceWM client is not initialized")
	}

	return d.client
}

func (d *DefaultAeroSpaceWindows) GetWindowByID(windowID int) (*windows.Window, error) {
	logger := logger.GetDefaultLogger()
	windowsList, err := d.client.Windows().GetAllWindows()
	if err != nil {
		return nil, err
	}
	logger.LogDebug("Windows found: %d", len(windowsList))

	for _, window := range windowsList {
		if window.WindowID == windowID {
			return &window, nil
		}
	}

	return nil, fmt.Errorf("window with ID %d not found", windowID)
}
