package cmd

import (
	"strings"
	"testing"

	aerospacecli "github.com/cristianoliveira/aerospace-ipc/pkg/client"
	"github.com/cristianoliveira/aerospace-marks/internal/logger"
	"github.com/cristianoliveira/aerospace-marks/internal/mocks"
	"github.com/cristianoliveira/aerospace-marks/internal/stdout"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
	"github.com/cristianoliveira/aerospace-marks/internal/testutils"
	"github.com/gkampitakis/go-snaps/snaps"
	"go.uber.org/mock/gomock"
)

func TestFocusCmd(t *testing.T) {
	t.Run("validate missing identifier", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger.SetDefaultLogger(&logger.EmptyLogger{})

		storageDbClient, _ := mocks.MockStorageDbClient(ctrl)
		strg, err := storage.NewMarkClient(storageDbClient)
		if err != nil {
			t.Fatal(err)
		}

		_, aerospaceClient := mocks.MockAerospaceConnection(ctrl)

		args := []string{"focus"}
		cmd := NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
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

		storageDbClient, strg := mocks.MockStorageDbClient(ctrl)
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
		cmd := NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
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
		storageDbClient, strg := mocks.MockStorageDbClient(ctrl)

		gomock.InOrder(
			storageDbClient.EXPECT().
				QueryOne(strings.TrimSpace(`
					SELECT * FROM marks WHERE mark = ?
				`), "nonexistent-mark").
				Return(nil, nil).
				Times(1),
		)
		_, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		stdout.ShouldExit = false

		cmd := NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
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
