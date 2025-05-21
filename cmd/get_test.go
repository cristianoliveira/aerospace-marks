package cmd

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/cristianoliveira/aerospace-marks/internal/mocks"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
	"github.com/cristianoliveira/aerospace-marks/internal/testutils"
	"github.com/cristianoliveira/aerospace-marks/pkgs/aerospacecli"
	"github.com/gkampitakis/go-snaps/snaps"
	"go.uber.org/mock/gomock"
)

func TestGetCommand(t *testing.T) {
	t.Run("shows only the marked windows", func(t *testing.T) {
		args := []string{"get", "mark1"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		storageMock, strg := mocks.MockStorageDbClient(ctrl)
		marks := []storage.Mark{
			{
				WindowID: "1",
				Mark:     "mark1",
			},
		}

		storageMock.EXPECT().
			QueryOne("SELECT * FROM marks WHERE mark = ?", "mark1").
			Return(&marks[0], nil)

		connectionMock, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospacecli.Window{
			{
				WindowID:    1,
				WindowTitle: "title1",
				AppName:     "app1",
			},
			{
				WindowID:    2,
				WindowTitle: "title2",
				AppName:     "app2",
			},
			{
				WindowID:    3,
				WindowTitle: "title3",
				AppName:     "app3",
			},
		}
		jsonData, err := json.Marshal(windows)
		if err != nil {
			t.Fatal(err)
		}
		connectionMock.EXPECT().
			SendCommand("list-windows", []string{"--all", "--json"}).
			Return(
				&aerospacecli.Response{
					ServerVersion: "1.0",
					StdOut:        string(jsonData),
					StdErr:        "",
					ExitCode:      0,
				}, nil).Times(1)

		cmd := NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		cmdAsString := "aerospace-marks " + strings.Join(args, " ") + "\n"
		snaps.MatchSnapshot(t, marks, windows, cmdAsString, out)
	})

	t.Run("shows only the marked windows id", func(t *testing.T) {
		// t.Skip("Skipping")
		args := []string{"get", "mark1", "--window-id"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		storageMock, strg := mocks.MockStorageDbClient(ctrl)
		marks := []storage.Mark{
			{
				WindowID: "1",
				Mark:     "mark1",
			},
		}

		storageMock.EXPECT().
			QueryOne("SELECT * FROM marks WHERE mark = ?", "mark1").
			Return(&marks[0], nil)

		_, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospacecli.Window{
			{
				WindowID:    1,
				WindowTitle: "title1",
				AppName:     "app1",
			},
			{
				WindowID:    2,
				WindowTitle: "title2",
				AppName:     "app2",
			},
			{
				WindowID:    3,
				WindowTitle: "title3",
				AppName:     "app3",
			},
		}

		cmd := NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		cmdAsString := "aerospace-marks " + strings.Join(args, " ") + "\n"
		snaps.MatchSnapshot(t, marks, windows, cmdAsString, out)
	})

	t.Run("shows only the marked windows app-name", func(t *testing.T) {
		// t.Skip("Skipping")
		args := []string{"get", "mark1", "-a"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		storageMock, strg := mocks.MockStorageDbClient(ctrl)
		marks := []storage.Mark{
			{
				WindowID: "1",
				Mark:     "mark1",
			},
		}

		storageMock.EXPECT().
			QueryOne("SELECT * FROM marks WHERE mark = ?", "mark1").
			Return(&marks[0], nil)

		connectionMock, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospacecli.Window{
			{
				WindowID:    1,
				WindowTitle: "title1",
				AppName:     "app1",
			},
			{
				WindowID:    2,
				WindowTitle: "title2",
				AppName:     "app2",
			},
			{
				WindowID:    3,
				WindowTitle: "title3",
				AppName:     "app3",
			},
		}
		jsonData, err := json.Marshal(windows)
		if err != nil {
			t.Fatal(err)
		}
		connectionMock.EXPECT().
			SendCommand("list-windows", []string{"--all", "--json"}).
			Return(
				&aerospacecli.Response{
					ServerVersion: "1.0",
					StdOut:        string(jsonData),
					StdErr:        "",
					ExitCode:      0,
				}, nil).Times(1)

		cmd := NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		cmdAsString := "aerospace-marks " + strings.Join(args, " ") + "\n"
		snaps.MatchSnapshot(t, marks, windows, cmdAsString, out)
	})
}
