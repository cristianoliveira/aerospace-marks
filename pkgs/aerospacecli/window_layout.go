package aerospacecli

import "fmt"

// Window layouts related functions

func (a *AeroSpaceWM) SetLayout(windowID int, layout string) error {
	windowStr := fmt.Sprintf("%d", windowID)
	args := []string{layout, "--window-id", windowStr}
	if res, err := a.Conn.SendCommand("layout", args); err != nil {
		fmt.Println("Error setting layout:", res, err)
		return fmt.Errorf("failed to set layout for window %d: %w", windowID, err)
	}
	return nil
}
