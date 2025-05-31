package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// Validate argument is not empty or whitespace
func ValidateArgIsNotEmpty(cmd *cobra.Command, args []string) error {
	if strings.TrimSpace(args[0]) == "" {
		return fmt.Errorf("argument cannot be empty or whitespace")
	}

	return nil
}
