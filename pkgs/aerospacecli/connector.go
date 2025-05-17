package aerospacecli

import (
	"fmt"
	"net"
)

// Connector should return AeroSpaceConnectiono

// Command represents the JSON structure for AeroSpace socket commands.
type Command struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Stdin   string   `json:"stdin"`
}

// Response represents the JSON structure from AeroSpace socket response.
type Response struct {
	ServerVersion string `json:"serverVersion"`
	StdErr        string `json:"stderr"`
	StdOut        string `json:"stdout"`
	ExitCode      int32  `json:"exitCode"`
}

type AeroSpaceConnection interface {
	// CloseConnection closes the connection to the AeroSpace socket.
	CloseConnection() error

	// SendCommand sends a command to the AeroSpace socket and returns the response.
	SendCommand(command string, args []string) (*Response, error)

	// GetAllWindows returns all windows
	GetAllWindows() ([]Window, error)
}


type AeroSpaceConnector interface {
	// Connect to the AeroSpace Socket and return client
	Connect() (*AeroSpaceConnection, error)
}

type AeroSpaceDefaultConnector struct {}

func (c *AeroSpaceDefaultConnector) Connect() (AeroSpaceConnection, error) {
	socketPath, err := GetSocketPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get socket path\n %w", err)
	}

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to socket\n %w", err)
	}

	client := &AeroSpaceDefaultConnection{
		SocketPath:          socketPath,
		MinAerospaceVersion: "0.15.2-Beta",
		Conn:                &conn,
	}

	return client, nil
}
