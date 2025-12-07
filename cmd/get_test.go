package cmd_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/cristianoliveira/aerospace-marks/cmd"
	"github.com/cristianoliveira/aerospace-marks/internal/logger"
	"github.com/cristianoliveira/aerospace-marks/internal/mocks"
	"github.com/cristianoliveira/aerospace-marks/internal/storage/db/queries"
	"github.com/cristianoliveira/aerospace-marks/internal/testutils"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	aerospace "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace"
	aerospacecli "github.com/cristianoliveira/aerospace-ipc/pkg/client"
)

//nolint:gocognit // TestGetCommand has high complexity due to multiple test cases
func TestGetCommand(t *testing.T) {
	logger.SetDefaultLogger(&logger.EmptyLogger{})

	t.Run("shows only the marked windows", func(t *testing.T) {
		args := []string{"get", "mark1"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		marks := []queries.Mark{
			{
				WindowID: 1,
				Mark:     "mark1",
			},
		}

		strg.EXPECT().
			GetWindowByMark("mark1").
			Return(&marks[0], nil).
			Times(1)

		connectionMock, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
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

		cmdAsString := "aerospace-marks " + strings.Join(args, " ") + "\n"
		snaps.MatchSnapshot(t, marks, windows, cmdAsString, out)
	})

	t.Run("shows only the marked windows id", func(t *testing.T) {
		// t.Skip("Skipping")
		args := []string{"get", "mark1", "--window-id"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		marks := []queries.Mark{
			{
				WindowID: 1,
				Mark:     "mark1",
			},
		}

		strg.EXPECT().
			GetWindowByMark("mark1").
			Return(&marks[0], nil).
			Times(1)

		_, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
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
			{
				WindowID:    3,
				WindowTitle: "title3",
				AppName:     "app3",
			},
		}

		cmd := cmd.NewRootCmd(strg, aerospaceClient)
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

		_, strg := mocks.MockStorageDBClient(ctrl)
		marks := []queries.Mark{
			{
				WindowID: 1,
				Mark:     "mark1",
			},
		}

		strg.EXPECT().
			GetWindowByMark("mark1").
			Return(&marks[0], nil).
			Times(1)

		connectionMock, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
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

		cmdAsString := "aerospace-marks " + strings.Join(args, " ") + "\n"
		snaps.MatchSnapshot(t, marks, windows, cmdAsString, out)
	})

	t.Run("shows only the marked windows app-bundle-id", func(t *testing.T) {
		// t.Skip("Skipping")
		args := []string{"get", "mark1", "--app-bundle-id"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		marks := []queries.Mark{
			{
				WindowID: 1,
				Mark:     "mark1",
			},
		}

		strg.EXPECT().
			GetWindowByMark("mark1").
			Return(&marks[0], nil).
			Times(1)

		connectionMock, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospace.Window{
			{
				WindowID:    1,
				WindowTitle: "title1",
				AppName:     "app1",
				AppBundleID: "bundle1",
			},
			{
				WindowID:    2,
				WindowTitle: "title2",
				AppName:     "app2",
				AppBundleID: "bundle2",
			},
			{
				WindowID:    3,
				WindowTitle: "title3",
				AppName:     "app3",
				AppBundleID: "bundle3",
			},
		}
		jsonData, err := json.Marshal(windows)
		if err != nil {
			t.Fatal(err)
		}
		connectionMock.EXPECT().
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

		cmdAsString := "aerospace-marks " + strings.Join(args, " ") + "\n"
		snaps.MatchSnapshot(t, marks, windows, cmdAsString, out)
	})

	t.Run("outputs JSON format", func(t *testing.T) {
		args := []string{"get", "mark1", "-o", "json"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		marks := []queries.Mark{
			{
				WindowID: 1,
				Mark:     "mark1",
			},
		}

		strg.EXPECT().
			GetWindowByMark("mark1").
			Return(&marks[0], nil).
			Times(1)

		connectionMock, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospace.Window{
			{
				WindowID:    1,
				WindowTitle: "title1",
				AppName:     "app1",
				Workspace:   "ws1",
				AppBundleID: "bundle1",
			},
		}
		jsonData, err := json.Marshal(windows)
		if err != nil {
			t.Fatal(err)
		}
		connectionMock.EXPECT().
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
		require.NoError(t, err)

		result := strings.TrimSpace(out)
		// Should be valid JSON OutputEvent
		var jsonResult map[string]interface{}
		err = json.Unmarshal([]byte(result), &jsonResult)
		require.NoError(t, err)
		assert.Equal(t, "get", jsonResult["command"])
		assert.InDelta(t, 1.0, jsonResult["window_id"], 0.0)
		assert.Equal(t, "app1", jsonResult["app_name"])
		assert.Equal(t, "title1", jsonResult["message"])
	})

	t.Run("outputs CSV format", func(t *testing.T) {
		args := []string{"get", "mark1", "-o", "csv"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		marks := []queries.Mark{
			{
				WindowID: 1,
				Mark:     "mark1",
			},
		}

		strg.EXPECT().
			GetWindowByMark("mark1").
			Return(&marks[0], nil).
			Times(1)

		connectionMock, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospace.Window{
			{
				WindowID:    1,
				WindowTitle: "title1",
				AppName:     "app1",
				Workspace:   "ws1",
				AppBundleID: "bundle1",
			},
		}
		jsonData, err := json.Marshal(windows)
		if err != nil {
			t.Fatal(err)
		}
		connectionMock.EXPECT().
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
		require.NoError(t, err)

		result := strings.TrimSpace(out)
		lines := strings.Split(result, "\n")
		assert.Len(t, lines, 2)
		assert.Equal(
			t,
			"command,action,window_id,app_name,workspace,target_workspace,result,message",
			lines[0],
		)
		assert.Contains(t, lines[1], "get,")
		assert.Contains(t, lines[1], ",1,")
		assert.Contains(t, lines[1], ",app1,")
		assert.Contains(t, lines[1], ",ws1,")
		assert.Contains(t, lines[1], ",title1")
	})

	t.Run("single field with JSON format uses OutputEvent", func(t *testing.T) {
		// Single field flags (-i, -t, -a, -b) with --output json/csv use OutputEvent format
		args := []string{"get", "mark1", "--window-id", "-o", "json"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		marks := []queries.Mark{
			{
				WindowID: 1,
				Mark:     "mark1",
			},
		}

		strg.EXPECT().
			GetWindowByMark("mark1").
			Return(&marks[0], nil).
			Times(1)

		_, aerospaceClient := mocks.MockAerospaceConnection(ctrl)

		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		require.NoError(t, err)

		result := strings.TrimSpace(out)
		// Should output OutputEvent JSON structure
		var jsonResult map[string]interface{}
		err = json.Unmarshal([]byte(result), &jsonResult)
		require.NoError(t, err)
		assert.Equal(t, "get", jsonResult["command"])
		assert.InDelta(t, 1.0, jsonResult["window_id"], 0.0)
		assert.Equal(t, "1", jsonResult["result"])
	})

	t.Run("single field without output flag outputs plain value", func(t *testing.T) {
		// Single field flags (-i, -t, -a, -b) without --output flag output plain value (backward compatibility)
		args := []string{"get", "mark1", "--window-id"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		marks := []queries.Mark{
			{
				WindowID: 1,
				Mark:     "mark1",
			},
		}

		strg.EXPECT().
			GetWindowByMark("mark1").
			Return(&marks[0], nil).
			Times(1)

		_, aerospaceClient := mocks.MockAerospaceConnection(ctrl)

		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		require.NoError(t, err)

		result := strings.TrimSpace(out)
		// Should output plain value (backward compatibility)
		assert.Equal(t, "1", result)
	})

	t.Run("snapshot test for text output format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		marks := []queries.Mark{
			{
				WindowID: 1,
				Mark:     "mark1",
			},
		}

		strg.EXPECT().
			GetWindowByMark("mark1").
			Return(&marks[0], nil).
			Times(1)

		connectionMock, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospace.Window{
			{
				WindowID:    1,
				WindowTitle: "Test Window Title",
				AppName:     "Test App",
				Workspace:   "workspace1",
				AppBundleID: "com.test.app",
			},
			{
				WindowID:    2,
				WindowTitle: "Another Window",
				AppName:     "Another App",
				Workspace:   "workspace2",
				AppBundleID: "com.another.app",
			},
		}
		jsonData, err := json.Marshal(windows)
		if err != nil {
			t.Fatal(err)
		}
		connectionMock.EXPECT().
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

		args := []string{"get", "mark1", "-o", "text"}
		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		require.NoError(t, err)

		result := strings.TrimSpace(out)
		cmdAsString := "aerospace-marks " + strings.Join(args, " ") + "\n"
		snaps.MatchSnapshot(t, marks, windows, cmdAsString, "result:\n", result)
	})

	t.Run("snapshot test for JSON output format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		marks := []queries.Mark{
			{
				WindowID: 1,
				Mark:     "mark1",
			},
		}

		strg.EXPECT().
			GetWindowByMark("mark1").
			Return(&marks[0], nil).
			Times(1)

		connectionMock, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospace.Window{
			{
				WindowID:    1,
				WindowTitle: "Test Window Title",
				AppName:     "Test App",
				Workspace:   "workspace1",
				AppBundleID: "com.test.app",
			},
			{
				WindowID:    2,
				WindowTitle: "Another Window",
				AppName:     "Another App",
				Workspace:   "workspace2",
				AppBundleID: "com.another.app",
			},
		}
		jsonData, err := json.Marshal(windows)
		if err != nil {
			t.Fatal(err)
		}
		connectionMock.EXPECT().
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

		args := []string{"get", "mark1", "-o", "json"}
		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		require.NoError(t, err)

		result := strings.TrimSpace(out)
		cmdAsString := "aerospace-marks " + strings.Join(args, " ") + "\n"
		snaps.MatchSnapshot(t, marks, windows, cmdAsString, "result:\n", result)
	})

	t.Run("snapshot test for CSV output format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		marks := []queries.Mark{
			{
				WindowID: 1,
				Mark:     "mark1",
			},
		}

		strg.EXPECT().
			GetWindowByMark("mark1").
			Return(&marks[0], nil).
			Times(1)

		connectionMock, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospace.Window{
			{
				WindowID:    1,
				WindowTitle: "Test Window Title",
				AppName:     "Test App",
				Workspace:   "workspace1",
				AppBundleID: "com.test.app",
			},
			{
				WindowID:    2,
				WindowTitle: "Another Window",
				AppName:     "Another App",
				Workspace:   "workspace2",
				AppBundleID: "com.another.app",
			},
		}
		jsonData, err := json.Marshal(windows)
		if err != nil {
			t.Fatal(err)
		}
		connectionMock.EXPECT().
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

		args := []string{"get", "mark1", "-o", "csv"}
		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		require.NoError(t, err)

		result := strings.TrimSpace(out)
		cmdAsString := "aerospace-marks " + strings.Join(args, " ") + "\n"
		snaps.MatchSnapshot(t, marks, windows, cmdAsString, "result:\n", result)
	})

	t.Run("snapshot test for single field JSON output", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		marks := []queries.Mark{
			{
				WindowID: 1,
				Mark:     "mark1",
			},
		}

		strg.EXPECT().
			GetWindowByMark("mark1").
			Return(&marks[0], nil).
			Times(1)

		connectionMock, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospace.Window{
			{
				WindowID:    1,
				WindowTitle: "Test Window Title",
				AppName:     "Test App",
				Workspace:   "workspace1",
				AppBundleID: "com.test.app",
			},
		}
		jsonData, err := json.Marshal(windows)
		if err != nil {
			t.Fatal(err)
		}
		connectionMock.EXPECT().
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

		args := []string{"get", "mark1", "--window-title", "-o", "json"}
		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		require.NoError(t, err)

		result := strings.TrimSpace(out)
		cmdAsString := "aerospace-marks " + strings.Join(args, " ") + "\n"
		snaps.MatchSnapshot(t, marks, windows, cmdAsString, "result:\n", result)
	})

	t.Run("snapshot test for single field CSV output", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		marks := []queries.Mark{
			{
				WindowID: 1,
				Mark:     "mark1",
			},
		}

		strg.EXPECT().
			GetWindowByMark("mark1").
			Return(&marks[0], nil).
			Times(1)

		connectionMock, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospace.Window{
			{
				WindowID:    1,
				WindowTitle: "Test Window Title",
				AppName:     "Test App",
				Workspace:   "workspace1",
				AppBundleID: "com.test.app",
			},
		}
		jsonData, err := json.Marshal(windows)
		if err != nil {
			t.Fatal(err)
		}
		connectionMock.EXPECT().
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

		args := []string{"get", "mark1", "--app-name", "-o", "csv"}
		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		require.NoError(t, err)

		result := strings.TrimSpace(out)
		cmdAsString := "aerospace-marks " + strings.Join(args, " ") + "\n"
		snaps.MatchSnapshot(t, marks, windows, cmdAsString, "result:\n", result)
	})
}
