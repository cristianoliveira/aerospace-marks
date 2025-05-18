package cmd

import (
	"strings"
	"testing"

	"github.com/cristianoliveira/aerospace-marks/internal/logger"
	"github.com/cristianoliveira/aerospace-marks/internal/mocks"
	"github.com/cristianoliveira/aerospace-marks/internal/stdout"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
	"github.com/cristianoliveira/aerospace-marks/internal/testutils"
	"github.com/cristianoliveira/aerospace-marks/pkgs/aerospacecli"
	"github.com/gkampitakis/go-snaps/snaps"
	"go.uber.org/mock/gomock"
)

func TestFocusCmd(t *testing.T) {
	t.Run("validate missing identifier", func(t *testing.T) {
		logger.SetDefaultLogger(&logger.EmptyLogger{})

		args := []string{"focus"}
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

		storageDbClient, _ := mocks.MockStorageDbClient(ctrl)
		gomock.InOrder(
			storageDbClient.EXPECT().
				QueryOne(strings.TrimSpace(`
					SELECT * FROM marks WHERE mark = ?
				`), "mark1").
				Return(&storage.Mark{
					WindowID: "1",
					Mark:     "mark1",
				}, nil).
				Times(1),
		)

		mockAeroSpaceConnection, _ := mocks.MockAerospaceConnection(ctrl)
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
		out, err := testutils.CmdExecute(NewRootCmd(), args...)
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

		storageDbClient, _ := mocks.MockStorageDbClient(ctrl)
		gomock.InOrder(
			storageDbClient.EXPECT().
				QueryOne(strings.TrimSpace(`
					SELECT * FROM marks WHERE mark = ?
				`), "nonexistent-mark").
				Return(nil, nil).
				Times(1),
		)

		stdout.ShuldExit = false
		out, err := testutils.CmdExecute(rootCmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		cmdAsString := "aerospace-marks " + strings.Join(args, " ") + "\n"
		snaps.MatchSnapshot(t, cmdAsString, out)
	})
}
