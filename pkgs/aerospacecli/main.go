package aerospacecli

import (
	"fmt"
)

type AeroSpaceDefaultConnection struct {
	MinAerospaceVersion string
	Conn                AeroSpaceSocketConn
}

// NewAeroSpaceClient creates a new AeroSpaceClient with the default socket path.
// It checks for environment variable AEROSPACESOCK or uses the default socket path.
// which is usually /tmp/bobko.aerospace-<username>.sock
func NewAeroSpaceConnection() (*AeroSpaceDefaultConnection, error) {
	conn, err := DefaultConnector.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to socket\n %w", err)
	}

	client := &AeroSpaceDefaultConnection{
		MinAerospaceVersion: "0.15.2-Beta",
		Conn:                conn,
	}

	return client, nil
}
