/*
Copyright © 2025 Cristian Oliveira license@cristianoliveira.dev
*/
package cmd

import (
	"fmt"

	"github.com/cristianoliveira/aerospace-marks/internal/storage"

	"github.com/spf13/cobra"
)

// unmarkCmd represents the unmark command
func UnmarkCmd(storageClient storage.MarkStorage) *cobra.Command {
	unmarkCmd := &cobra.Command{
		Use:   "unmark",
		Short: "Unmark one or more windows by identifier",
		Long: `Unmark one or more windows by identifier.

unmark [<identifier>]

unmark cmd will remove identifier from the list of current marks on a window. If identifier is omitted, all marks are removed.
	`,

		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				rowsEffected, err := storageClient.DeleteAllMarks()
				if err != nil {
					return err
				}

				fmt.Printf("Removed %d marks\n", rowsEffected)
				return nil
			}

			var count int
			for _, identifier := range args {
				rowsAffected, err := storageClient.DeleteByMark(identifier)
				if err != nil {
					return err
				}
				if rowsAffected == 0 {
					return fmt.Errorf("mark '%s' not found", identifier)
				}
				count += int(rowsAffected)
			}

			fmt.Printf("Removed %d marks\n", count)
			return nil
		},
	}

	return unmarkCmd
}
