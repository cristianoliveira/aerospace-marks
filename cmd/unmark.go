/*
Copyright © 2025 Cristian Oliveira licence@cristianoliveira.dev
*/
package cmd

import (
	"fmt"

	"github.com/cristianoliveira/aerospace-marks/internal/storage"

	"github.com/spf13/cobra"
)

// unmarkCmd represents the unmark command
func UnmarkCmd() *cobra.Command {
	unmarkCmd := &cobra.Command{
		Use:   "unmark",
		Short: "Unmark one or more windows by identifier",
		Long: `Unmark one or more windows by identifier.

unmark [<identifier>]

unmark cmd will remove identifier from the list of current marks on a window. If identifier is omitted, all marks are removed.
	`,

		RunE: func(cmd *cobra.Command, args []string) error {
			markClient, err := storage.NewMarkClient()
			if err != nil {
				return err
			}
			defer markClient.Close()

			if len(args) == 0 {
				rowsEffected, err := markClient.DeleteAllMarks()
				if err != nil {
					return err
				}

				fmt.Printf("Removed %d marks\n", rowsEffected)

				return nil
			}

			var count int
			for _, identifier := range args {
				if _, err := markClient.DeleteByMark(identifier); err != nil {
					return err
				}
				count++
			}

			fmt.Printf("Removed %d marks\n", count)
			return nil
		},
	}

	return unmarkCmd
}
