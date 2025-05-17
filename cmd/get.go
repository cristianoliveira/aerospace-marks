package cmd

import (
	"fmt"

	"github.com/cristianoliveira/aerospace-marks/internal/aerospace"
	"github.com/cristianoliveira/aerospace-marks/internal/stdout"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
	"github.com/spf13/cobra"
)

// getIdCmd represents the getId command
var getCmd = &cobra.Command{
	Use:   "get <mark> [flags]",
	Short: "Get windows marked with a mark",
	Long: `Get windows marked with a mark

USAGE:
aerospace-marks get <mark> # Will get the window ID | Name of the window marked with <mark>
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return stdout.ErrorAndExitf("No mark provided")
		}

		isGetWinID, _ := cmd.Flags().GetBool("window-id")

		markClient, err := storage.NewMarkClient()
		if err != nil {
			return stdout.ErrorAndExit(err)
		}
		defer markClient.Close()

		getID := args[0]

		// Get window ID by mark
		windowID, err := markClient.GetWindowIDByMark(getID)
		if err != nil {
			return stdout.ErrorAndExit(err)
		}
		if windowID == "" {
			return stdout.ErrorAndExitf("no window found for mark '%s'", getID)
		}

		if isGetWinID {
			fmt.Printf("%s", windowID)
			return nil
		}

		// Get window name by mark
		windowInfo, err := aerospace.GetWindowByID(windowID)
		if err != nil {
			return stdout.ErrorAndExit(err)
		}

		fmt.Printf("%s", windowInfo)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().BoolP("window-id", "w", false, "Get a Window ID by a mark")
}
