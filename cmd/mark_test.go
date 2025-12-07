package cmd_test

import (
	"encoding/json"
	"testing"

	"github.com/cristianoliveira/aerospace-marks/cmd"
	"github.com/cristianoliveira/aerospace-marks/internal/logger"
	"github.com/cristianoliveira/aerospace-marks/internal/mocks"
	"github.com/cristianoliveira/aerospace-marks/internal/testutils"
	"github.com/gkampitakis/go-snaps/snaps"
	"go.uber.org/mock/gomock"

	aerospace "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace"
	aerospacecli "github.com/cristianoliveira/aerospace-ipc/pkg/client"
)

//nolint:gocognit // TestMarkCommand has high complexity due to multiple test scenarios
func TestMarkCommand(t *testing.T) {
	t.Run("marks the focused window - `marks mark mark1`", func(t *testing.T) {
		command := "mark"
		args := []string{command, "mark1"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			ReplaceAllMarks(1, "mark1").
			Return(int64(1), nil).
			Times(1)

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
					"%{window-id} %{window-title} %{app-name} %{app-bundle-id} %{workspace}",
				}).
			Return(
				&aerospacecli.Response{
					ServerVersion: "1.0",
					StdOut:        string(jsonData),
					StdErr:        "",
					ExitCode:      0,
				}, nil).Times(1)

		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		snapshot := testutils.RenderSnapshotSpec(testutils.SnapshotSpec{
			Command: testutils.CommandString(args...),
			Stdout:  out,
			Contexts: []testutils.SnapshotContext{
				testutils.Context("windows", windows),
			},
		})
		snaps.MatchSnapshot(t, snapshot)
	})

	t.Run("marks window by id - `marks mark mark1 --window-id 2`", func(t *testing.T) {
		command := "mark"
		args := []string{command, "mark1", "--window-id", "2"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			ReplaceAllMarks(2, "mark1").
			Return(int64(1), nil).
			Times(1)

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
					"%{window-id} %{window-title} %{app-name} %{app-bundle-id} %{workspace}",
				}).
			Return(
				&aerospacecli.Response{
					ServerVersion: "1.0",
					StdOut:        string(jsonData),
					StdErr:        "",
					ExitCode:      0,
				}, nil).Times(1)

		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		snapshot := testutils.RenderSnapshotSpec(testutils.SnapshotSpec{
			Command: testutils.CommandString(args...),
			Stdout:  out,
			Contexts: []testutils.SnapshotContext{
				testutils.Context("windows", windows),
			},
		})
		snaps.MatchSnapshot(t, snapshot)
	})

	t.Run("marks the focused window - `marks mark --add`", func(t *testing.T) {
		command := "mark"
		args := []string{command, "mark2", "--add"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			AddMark(1, "mark2").
			Return(nil).
			Times(1)

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
					"%{window-id} %{window-title} %{app-name} %{app-bundle-id} %{workspace}",
				}).
			Return(
				&aerospacecli.Response{
					ServerVersion: "1.0",
					StdOut:        string(jsonData),
					StdErr:        "",
					ExitCode:      0,
				}, nil).Times(1)

		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		snapshot := testutils.RenderSnapshotSpec(testutils.SnapshotSpec{
			Command: testutils.CommandString(args...),
			Stdout:  out,
			Contexts: []testutils.SnapshotContext{
				testutils.Context("windows", windows),
			},
		})
		snaps.MatchSnapshot(t, snapshot)
	})

	t.Run("validates missing identifier - `marks mark`", func(t *testing.T) {
		command := "mark"
		args := []string{command} // Missing identifier

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		logger.SetDefaultLogger(&logger.EmptyLogger{})

		_, aerospaceClient := mocks.MockAerospaceConnection(ctrl)

		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err == nil {
			t.Fatal(err)
		}
		if out != "" {
			t.Fatal("output should be empty", out)
		}

		snapshot := testutils.RenderSnapshotSpec(testutils.SnapshotSpec{
			Command: testutils.CommandString(args...),
			Stdout:  out,
			Stderr:  err.Error(),
		})
		snaps.MatchSnapshot(t, snapshot)
	})

	t.Run("marks toggles mark (remove) - `marks foobar --toggle`", func(t *testing.T) {
		command := "mark"
		args := []string{command, "foobar", "--toggle"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			ToggleMark(2, "foobar").
			Return(nil).
			Times(1)

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
					"%{window-id} %{window-title} %{app-name} %{app-bundle-id} %{workspace}",
				}).
			Return(
				&aerospacecli.Response{
					ServerVersion: "1.0",
					StdOut:        string(jsonData),
					StdErr:        "",
					ExitCode:      0,
				}, nil).Times(1)

		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		snapshot := testutils.RenderSnapshotSpec(testutils.SnapshotSpec{
			Command: testutils.CommandString(args...),
			Stdout:  out,
			Contexts: []testutils.SnapshotContext{
				testutils.Context("windows", windows),
			},
		})
		snaps.MatchSnapshot(t, snapshot)
	})

	t.Run("marks toggles mark (adding) - `marks foobar --toggle`", func(t *testing.T) {
		command := "mark"
		args := []string{command, "foobar", "--toggle"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			ToggleMark(2, "foobar").
			Return(nil).
			Times(1)

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
					"%{window-id} %{window-title} %{app-name} %{app-bundle-id} %{workspace}",
				}).
			Return(
				&aerospacecli.Response{
					ServerVersion: "1.0",
					StdOut:        string(jsonData),
					StdErr:        "",
					ExitCode:      0,
				}, nil).Times(1)

		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		snapshot := testutils.RenderSnapshotSpec(testutils.SnapshotSpec{
			Command: testutils.CommandString(args...),
			Stdout:  out,
			Contexts: []testutils.SnapshotContext{
				testutils.Context("windows", windows),
			},
		})
		snaps.MatchSnapshot(t, snapshot)
	})

	t.Run("fails when empty identifier - `marks ''`", func(t *testing.T) {
		command := "mark"
		args := []string{command, ""}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		_, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospace.Window{
			{
				WindowID:    2,
				WindowTitle: "title2",
				AppName:     "app2",
			},
		}

		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if out != "" {
			t.Fatal("output should be empty", out)
		}
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		snapshot := testutils.RenderSnapshotSpec(testutils.SnapshotSpec{
			Command: testutils.CommandString(args...),
			Stdout:  out,
			Stderr:  err.Error(),
			Contexts: []testutils.SnapshotContext{
				testutils.Context("windows", windows),
			},
		})
		snaps.MatchSnapshot(t, snapshot)
	})
}
