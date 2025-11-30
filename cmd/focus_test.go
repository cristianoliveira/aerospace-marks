package cmd_test

import (
	"strings"
	"testing"

	"github.com/cristianoliveira/aerospace-marks/cmd"
	"github.com/cristianoliveira/aerospace-marks/internal/logger"
	"github.com/cristianoliveira/aerospace-marks/internal/mocks"
	"github.com/cristianoliveira/aerospace-marks/internal/stdout"
	"github.com/cristianoliveira/aerospace-marks/internal/testutils"
	"github.com/gkampitakis/go-snaps/snaps"
	"go.uber.org/mock/gomock"

	aerospacecli "github.com/cristianoliveira/aerospace-ipc/pkg/client"
)

func TestFocusCmd(t *testing.T) {
	t.Run("validate missing identifier", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger.SetDefaultLogger(&logger.EmptyLogger{})

		_, strg := mocks.MockStorageDBClient(ctrl)

		_, aerospaceClient := mocks.MockAerospaceConnection(ctrl)

		args := []string{"focus"}
		rootCmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(rootCmd, args...)
		if out != "" {
			t.Fatal("output should be empty", out)
		}
		if err == nil {
			t.Fatal(err)
		}

		cmdAsString := "aerospace-marks " + strings.Join(args, " ") + "\n"
		snaps.MatchSnapshot(t, cmdAsString, err.Error())
	})

	t.Run("focus to a window by mark - `marks focus mark1`", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			GetWindowIDByMark("mark1").
			Return(1, nil).
			Times(1)

		mockAeroSpaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		mockAeroSpaceConnection.EXPECT().
			SendCommand("focus", []string{"--window-id", "1"}).
			Return(
				&aerospacecli.Response{
					ServerVersion: "1.0",
					StdOut:        "",
					StdErr:        "",
					ExitCode:      0,
				}, nil).Times(1)

		args := []string{"focus", "mark1"}
		rootCmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(rootCmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		cmdAsString := "aerospace-marks " + strings.Join(args, " ") + "\n"
		snaps.MatchSnapshot(t, cmdAsString, out)
	})

	t.Run("focus using mark that does not exist", func(t *testing.T) {
		command := "focus"
		args := []string{command, "nonexistent-mark"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, strg := mocks.MockStorageDBClient(ctrl)

		strg.EXPECT().
			GetWindowIDByMark("nonexistent-mark").
			Return(0, nil).
			Times(1)

		_, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		//nolint:reassign // Test utility needs to modify package variable
		stdout.ShouldExit = false

		rootCmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(rootCmd, args...)
		if err == nil {
			t.Fatal(err)
		}
		if out != "" {
			t.Fatal("output should be empty", out)
		}

		cmdAsString := "aerospace-marks " + strings.Join(args, " ") + "\n"
		snaps.MatchSnapshot(t, cmdAsString, err.Error())
	})
}
