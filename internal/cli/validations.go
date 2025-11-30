package cli

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
)

// ValidateArgIsNotEmpty validates that an argument is not empty or whitespace.
func ValidateArgIsNotEmpty(cmd *cobra.Command, args []string) error {
	if strings.TrimSpace(args[0]) == "" {
		return errors.New("argument cannot be empty or whitespace")
	}

	return nil
}
