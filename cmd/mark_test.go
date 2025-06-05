package cmd

import (
	"encoding/json"
	"strings"
	"testing"

	aerospace "github.com/cristianoliveira/aerospace-ipc"
	aerospacecli "github.com/cristianoliveira/aerospace-ipc/pkg/client"
	"github.com/cristianoliveira/aerospace-marks/internal/logger"
	"github.com/cristianoliveira/aerospace-marks/internal/mocks"
	"github.com/cristianoliveira/aerospace-marks/internal/testutils"
	"github.com/gkampitakis/go-snaps/snaps"
	"go.uber.org/mock/gomock"
)

func TestMarkCommand(t *testing.T) {
	t.Run("marks the focused window - `marks mark mark1`", func(t *testing.T) {
		command := "mark"
		args := []string{command, "mark1"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		storageDbClient, strg := mocks.MockStorageDbClient(ctrl)
		dbResult := mocks.MockStorageDbResult(ctrl, nil, &[]int64{1}[0])
		gomock.InOrder(
			storageDbClient.EXPECT().
				Execute(strings.TrimSpace(`
					DELETE FROM marks WHERE window_id = ? OR mark = ?
				`),
					"1", "mark1").
				Return(dbResult, nil).
				Times(1),
			storageDbClient.EXPECT().
				Execute(strings.TrimSpace(`
					INSERT INTO marks (window_id, mark) VALUES (?, ?)
				`), "1", "mark1").
				Return(dbResult, nil).
				Times(1),
		)

		mockAeroSpaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospace.Window{
			{
				WindowID:    1,
				WindowTitle: "title1",
				AppName:     "app1",
			},
		}
		jsonData, err := json.Marshal(windows)
		if err != nil {
			t.Fatal(err)
		}
		mockAeroSpaceConnection.EXPECT().
			SendCommand(
				"list-windows",
				[]string{
					"--focused",
					"--json",
					"--format",
					"%{window-id} %{app-name} %{app-bundle-id} %{workspace}",
				}).
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
		snaps.MatchSnapshot(t, windows, cmdAsString, out)
	})

	t.Run("marks window by id - `marks mark mark1 --window-id 2`", func(t *testing.T) {
		command := "mark"
		args := []string{command, "mark1", "--window-id", "2"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		storageDbClient, strg := mocks.MockStorageDbClient(ctrl)
		dbResult := mocks.MockStorageDbResult(ctrl, nil, &[]int64{1}[0])
		gomock.InOrder(
			storageDbClient.EXPECT().
				Execute(strings.TrimSpace(`
					DELETE FROM marks WHERE window_id = ? OR mark = ?
				`),
					"2", "mark1").
				Return(dbResult, nil).
				Times(1),
			storageDbClient.EXPECT().
				Execute(strings.TrimSpace(`
					INSERT INTO marks (window_id, mark) VALUES (?, ?)
				`), "2", "mark1").
				Return(dbResult, nil).
				Times(1),
		)

		mockAeroSpaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospace.Window{
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
		}
		jsonData, err := json.Marshal(windows)
		if err != nil {
			t.Fatal(err)
		}
		mockAeroSpaceConnection.EXPECT().
			SendCommand(
				"list-windows",
				[]string{
					"--all",
					"--json",
					"--format",
					"%{window-id} %{app-name} %{app-bundle-id} %{workspace}",
				}).
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
		snaps.MatchSnapshot(t, windows, cmdAsString, out)
	})

	t.Run("marks the focused window - `marks mark --add`", func(t *testing.T) {
		command := "mark"
		args := []string{command, "mark2", "--add"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		storageDbClient, strg := mocks.MockStorageDbClient(ctrl)
		gomock.InOrder(
			storageDbClient.EXPECT().
				Execute(strings.TrimSpace(`
					INSERT INTO marks (window_id, mark) VALUES (?, ?)
				`), "1", "mark2").
				Times(1),
		)

		mockAeroSpaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospace.Window{
			{
				WindowID:    1,
				WindowTitle: "title1",
				AppName:     "app1",
			},
		}
		jsonData, err := json.Marshal(windows)
		if err != nil {
			t.Fatal(err)
		}
		mockAeroSpaceConnection.EXPECT().
			SendCommand(
				"list-windows",
				[]string{
					"--focused",
					"--json",
					"--format",
					"%{window-id} %{app-name} %{app-bundle-id} %{workspace}",
				}).
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
		snaps.MatchSnapshot(t, windows, cmdAsString, out)
	})

	t.Run("validates missing identifier - `marks mark`", func(t *testing.T) {
		command := "mark"
		args := []string{command} // Missing identifier

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDbClient(ctrl)
		logger.SetDefaultLogger(&logger.EmptyLogger{})

		_, aerospaceClient := mocks.MockAerospaceConnection(ctrl)

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

	t.Run("marks toggles mark (remove) - `marks foobar --toggle`", func(t *testing.T) {
		command := "mark"
		args := []string{command, "foobar", "--toggle"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		storageDbClient, strg := mocks.MockStorageDbClient(ctrl)
		dbResult := mocks.MockStorageDbResult(ctrl, nil, &[]int64{1}[0])
		gomock.InOrder(
			storageDbClient.EXPECT().
				Execute(strings.TrimSpace(`
					DELETE FROM marks WHERE mark = ?
				`), "foobar").
				Return(dbResult, nil).
				Times(1),
			storageDbClient.EXPECT().
				Execute(strings.TrimSpace(`
					INSERT INTO marks (window_id, mark) VALUES (?, ?)
				`), "1", "mark1").
				Return(dbResult, nil).
				Times(0),
		)

		mockAeroSpaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospace.Window{
			{
				WindowID:    2,
				WindowTitle: "title2",
				AppName:     "app2",
			},
		}
		jsonData, err := json.Marshal(windows)
		if err != nil {
			t.Fatal(err)
		}
		mockAeroSpaceConnection.EXPECT().
			SendCommand(
				"list-windows",
				[]string{
					"--focused",
					"--json",
					"--format",
					"%{window-id} %{app-name} %{app-bundle-id} %{workspace}",
				}).
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
		snaps.MatchSnapshot(t, windows, cmdAsString, out)
	})

	t.Run("marks toggles mark (adding) - `marks foobar --toggle`", func(t *testing.T) {
		command := "mark"
		args := []string{command, "foobar", "--toggle"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		storageDbClient, strg := mocks.MockStorageDbClient(ctrl)
		dbResult := mocks.MockStorageDbResult(ctrl, nil, &[]int64{0}[0])
		gomock.InOrder(
			storageDbClient.EXPECT().
				Execute(strings.TrimSpace(`
					DELETE FROM marks WHERE mark = ?
				`), "foobar").
				Return(dbResult, nil).
				Times(1),
			storageDbClient.EXPECT().
				Execute(strings.TrimSpace(`
					INSERT INTO marks (window_id, mark) VALUES (?, ?)
				`), "2", "foobar").
				Return(dbResult, nil).
				Times(1),
		)

		mockAeroSpaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospace.Window{
			{
				WindowID:    2,
				WindowTitle: "title2",
				AppName:     "app2",
			},
		}
		jsonData, err := json.Marshal(windows)
		if err != nil {
			t.Fatal(err)
		}
		mockAeroSpaceConnection.EXPECT().
			SendCommand(
				"list-windows",
				[]string{
					"--focused",
					"--json",
					"--format",
					"%{window-id} %{app-name} %{app-bundle-id} %{workspace}",
				}).
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
		snaps.MatchSnapshot(t, windows, cmdAsString, out)
	})

	t.Run("fails when empty identifier - `marks ''`", func(t *testing.T) {
		command := "mark"
		args := []string{command, ""}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDbClient(ctrl)
		_, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospace.Window{
			{
				WindowID:    2,
				WindowTitle: "title2",
				AppName:     "app2",
			},
		}

		cmd := NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if out != "" {
			t.Fatal("output should be empty", out)
		}
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		cmdAsString := "aerospace-marks " + strings.Join(args, " ") + "\n"
		snaps.MatchSnapshot(t, windows, cmdAsString, "Error", err.Error())
	})
}
