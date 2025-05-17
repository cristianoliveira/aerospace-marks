package aerospacecli

import (
	"encoding/json"
	"fmt"
	"net"
)

type AeroSpaceDefaultConnection struct {
	SocketPath          string
	MinAerospaceVersion string
	Conn                *net.Conn
}

// NewAeroSpaceClient creates a new AeroSpaceClient with the default socket path.
// It checks for environment variable AEROSPACESOCK or uses the default socket path.
// which is usually /tmp/bobko.aerospace-<username>.sock
func NewAeroSpaceConnection() (*AeroSpaceDefaultConnection, error) {
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

// Close closes the connection to the AeroSpace socket.
func (c *AeroSpaceDefaultConnection) CloseConnection() error {
	if c.Conn != nil {
		err := (*c.Conn).Close()
		if err != nil {
			return fmt.Errorf("failed to close connection\n %w", err)
		}
	}
	return nil
}

// SendCommand sends a command to the AeroSpace window manager via Unix socket and returns the response.
func (c *AeroSpaceDefaultConnection) SendCommand(command string, args []string) (*Response, error) {
	if c.Conn == nil {
		return nil, fmt.Errorf("connection is not established")
	}

	// Merge command and arguments into the Command struct
	commandArgs := append([]string{command}, args...)
	cmd := Command{
		Command: "", // This field is deprecated and not used
		Args:    commandArgs,
		Stdin:   "",
	}

	cmdBytes, err := json.Marshal(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal command\n%w", err)
	}

	_, err = (*c.Conn).Write(cmdBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to send command\n%w", err)
	}

	buf := make([]byte, 4096)
	n, err := (*c.Conn).Read(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read response\n%w", err)
	}

	var response Response
	err = json.Unmarshal(buf[:n], &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response\n%w", err)
	}

	if response.ExitCode != 0 {
		return nil, fmt.Errorf("command failed with exit code %d\n%s", response.ExitCode, response.StdErr)
	}

	if response.StdErr != "" {
		return nil, fmt.Errorf("command error\n %s", response.StdErr)
	}

	return &response, nil
}
