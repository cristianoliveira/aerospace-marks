package cmd_test

import (
	"encoding/csv"
	"encoding/json"
	"strings"
	"testing"

	"github.com/cristianoliveira/aerospace-marks/cmd"
	"github.com/cristianoliveira/aerospace-marks/internal/logger"
	"github.com/cristianoliveira/aerospace-marks/internal/mocks"
	"github.com/cristianoliveira/aerospace-marks/internal/stdout"
	"github.com/cristianoliveira/aerospace-marks/internal/storage/db/queries"
	"github.com/cristianoliveira/aerospace-marks/internal/testutils"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	aerospace "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace"
	aerospacecli "github.com/cristianoliveira/aerospace-ipc/pkg/client"
)

//nolint:gocognit,gocyclo // TestListCommand has high complexity due to multiple test cases
func TestListCommand(t *testing.T) {
	t.Run("shows only the marked windows", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			GetMarks().
			Return(
				[]queries.Mark{
					{
						WindowID: 1,
						Mark:     "mark1",
					},
					{
						WindowID: 2,
						Mark:     "mark2",
					},
				}, nil,
			).Times(1)

		mockAeroSpaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospace.Window{
			{
				WindowID:    1,
				WindowTitle: "title1",
				AppName:     "app1",
				Workspace:   "",
				AppBundleID: "",
			},
			{
				WindowID:    2,
				WindowTitle: "title2",
				AppName:     "app2",
				Workspace:   "",
				AppBundleID: "",
			},
			{
				WindowID:    3,
				WindowTitle: "title3",
				AppName:     "app3",
				Workspace:   "",
				AppBundleID: "",
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

		args := []string{"list", "-o", "text"}
		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		lines := strings.Split(result, "\n")
		assert.Len(t, lines, 2)
		assert.Equal(t, []string{
			"mark1 | 1 | app1 | title1 | _ | _",
			"mark2 | 2 | app2 | title2 | _ | _",
		}, lines)
	})

	t.Run("shows no marked window found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			GetMarks().
			Return(
				[]queries.Mark{
					{
						WindowID: 1,
						Mark:     "mark1",
					},
					{
						WindowID: 2,
						Mark:     "mark2",
					},
				}, nil,
			).Times(1)

		mockAeroSpaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospace.Window{
			{
				WindowID:    111,
				WindowTitle: "title1",
				AppName:     "app1",
				Workspace:   "",
				AppBundleID: "",
			},
			{
				WindowID:    222,
				WindowTitle: "title2",
				AppName:     "app2",
				Workspace:   "",
				AppBundleID: "",
			},
			{
				WindowID:    333,
				WindowTitle: "title3",
				AppName:     "app3",
				Workspace:   "",
				AppBundleID: "",
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

		args := []string{"list", "-o", "text"}
		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		assert.Equal(t, "No marked window found", result)
	})

	t.Run("shows no marks found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			GetMarks().
			Return([]queries.Mark{}, nil).
			Times(1)

		_, aerospaceClient := mocks.MockAerospaceConnection(ctrl)

		args := []string{"list", "-o", "text"}
		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		assert.Equal(t, "No marks found", result)
	})

	t.Run("print all marked windows", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		marks, err := mocks.LoadMarksFixture(
			"../internal/mocks/fixtures/storage/list-marked-windows.json",
		)
		if err != nil {
			t.Fatal(err)
		}
		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			GetMarks().
			Return(marks, nil).
			Times(1)

		windows, err := mocks.LoadAeroWindowsFixtureRaw(
			"../internal/mocks/fixtures/aerospace/list-windows-all.json",
		)
		if err != nil {
			t.Fatal(err)
		}

		mockAeroSpaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		mockAeroSpaceConnection.EXPECT().
			SendCommand(
				"list-windows",
				[]string{
					"--all",
					"--json",
					"--format",
					"%{window-id} %{window-title} %{app-name} %{app-bundle-id} %{workspace}",
				}).
			Return(&aerospacecli.Response{
				ServerVersion: "1.0",
				StdOut:        windows,
				StdErr:        "",
				ExitCode:      0,
			}, nil).Times(1)

		args := []string{"list", "-o", "text"}
		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		cmdAsString := "\n$aerospace-marks " + strings.Join(args, " ") + "\n"
		snaps.MatchSnapshot(t, "context\n", windows, marks, cmdAsString, result)
	})

	t.Run("outputs JSON format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			GetMarks().
			Return(
				[]queries.Mark{
					{
						WindowID: 1,
						Mark:     "mark1",
					},
				}, nil,
			).Times(1)

		mockAeroSpaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospace.Window{
			{
				WindowID:    1,
				WindowTitle: "title1",
				AppName:     "app1",
				Workspace:   "workspace1",
				AppBundleID: "bundle1",
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

		args := []string{"list", "-o", "json"}
		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		// Verify it's valid JSON
		var jsonResult []map[string]interface{}
		err = json.Unmarshal([]byte(result), &jsonResult)
		require.NoError(t, err)
		assert.Len(t, jsonResult, 1)
		assert.Equal(t, "mark1", jsonResult[0]["mark"])
		assert.InDelta(t, 1.0, jsonResult[0]["window_id"], 0.0)
		assert.Equal(t, "app1", jsonResult[0]["app_name"])
	})

	t.Run("outputs CSV format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			GetMarks().
			Return(
				[]queries.Mark{
					{
						WindowID: 1,
						Mark:     "mark1",
					},
				}, nil,
			).Times(1)

		mockAeroSpaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospace.Window{
			{
				WindowID:    1,
				WindowTitle: "title1",
				AppName:     "app1",
				Workspace:   "workspace1",
				AppBundleID: "bundle1",
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

		args := []string{"list", "-o", "csv"}
		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		lines := strings.Split(result, "\n")
		assert.Len(t, lines, 2)
		assert.Equal(t, "mark,window_id,app_name,window_title,workspace,app_bundle_id", lines[0])
		assert.Equal(t, "mark1,1,app1,title1,workspace1,bundle1", lines[1])
	})

	t.Run("defaults to text format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			GetMarks().
			Return(
				[]queries.Mark{
					{
						WindowID: 1,
						Mark:     "mark1",
					},
				}, nil,
			).Times(1)

		mockAeroSpaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospace.Window{
			{
				WindowID:    1,
				WindowTitle: "title1",
				AppName:     "app1",
				Workspace:   "",
				AppBundleID: "",
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

		args := []string{"list"}
		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		// Should be text format (pipe-separated)
		assert.Contains(t, result, "|")
		assert.Contains(t, result, "mark1")
	})

	t.Run("empty marks returns empty JSON array", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			GetMarks().
			Return([]queries.Mark{}, nil).
			Times(1)

		_, aerospaceClient := mocks.MockAerospaceConnection(ctrl)

		args := []string{"list", "-o", "json"}
		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		assert.Equal(t, "[]", result)
	})

	t.Run("empty marks returns CSV header only", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			GetMarks().
			Return([]queries.Mark{}, nil).
			Times(1)

		_, aerospaceClient := mocks.MockAerospaceConnection(ctrl)

		args := []string{"list", "-o", "csv"}
		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		assert.Equal(t, "mark,window_id,app_name,window_title,workspace,app_bundle_id", result)
	})

	t.Run("invalid output format returns error", func(t *testing.T) {
		logger.SetDefaultLogger(&logger.EmptyLogger{})
		//nolint:reassign // Test utility needs to modify package variable
		stdout.ShouldExit = false
		defer func() {
			//nolint:reassign // Test utility needs to restore package variable
			stdout.ShouldExit = true
		}()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		_, aerospaceClient := mocks.MockAerospaceConnection(ctrl)

		args := []string{"list", "-o", "invalid"}
		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		_, err := testutils.CmdExecute(cmd, args...)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported output format")
	})

	t.Run("JSON output with special characters", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			GetMarks().
			Return(
				[]queries.Mark{
					{
						WindowID: 1,
						Mark:     "mark\"with\"quotes",
					},
				}, nil,
			).Times(1)

		mockAeroSpaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospace.Window{
			{
				WindowID:    1,
				WindowTitle: "title\nwith\nnewlines",
				AppName:     "app\twith\ttabs",
				Workspace:   "workspace",
				AppBundleID: "bundle.id",
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

		args := []string{"list", "-o", "json"}
		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		// Should be valid JSON
		var jsonResult []map[string]interface{}
		err = json.Unmarshal([]byte(result), &jsonResult)
		require.NoError(t, err)
		assert.Len(t, jsonResult, 1)
		assert.Equal(t, "mark\"with\"quotes", jsonResult[0]["mark"])
	})

	t.Run("CSV output with commas and quotes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			GetMarks().
			Return(
				[]queries.Mark{
					{
						WindowID: 1,
						Mark:     "mark,with,commas",
					},
				}, nil,
			).Times(1)

		mockAeroSpaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospace.Window{
			{
				WindowID:    1,
				WindowTitle: "title\"with\"quotes",
				AppName:     "app,name",
				Workspace:   "workspace",
				AppBundleID: "bundle.id",
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

		args := []string{"list", "-o", "csv"}
		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		lines := strings.Split(result, "\n")
		assert.Len(t, lines, 2)
		// CSV should properly escape commas and quotes
		assert.Contains(t, lines[1], "mark,with,commas")
		// Verify it's valid CSV
		reader := csv.NewReader(strings.NewReader(result))
		records, err := reader.ReadAll()
		require.NoError(t, err)
		assert.Len(t, records, 2)
	})

	t.Run("text output format preserves alignment", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			GetMarks().
			Return(
				[]queries.Mark{
					{
						WindowID: 1,
						Mark:     "a",
					},
					{
						WindowID: 2,
						Mark:     "very-long-mark-name-here",
					},
				}, nil,
			).Times(1)

		mockAeroSpaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospace.Window{
			{
				WindowID:    1,
				WindowTitle: "short",
				AppName:     "App",
				Workspace:   "ws",
				AppBundleID: "bundle",
			},
			{
				WindowID:    2,
				WindowTitle: "very long window title here",
				AppName:     "Very Long App Name",
				Workspace:   "very-long-workspace",
				AppBundleID: "very.long.bundle.id",
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

		args := []string{"list", "-o", "text"}
		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		lines := strings.Split(result, "\n")
		assert.Len(t, lines, 2)
		// Both lines should have pipe separators
		for _, line := range lines {
			assert.Contains(t, line, "|")
			fields := strings.Split(line, "|")
			assert.Len(t, fields, 6, "Each line should have 6 fields")
		}
	})

	t.Run("JSON output with multiple windows maintains array structure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			GetMarks().
			Return(
				[]queries.Mark{
					{WindowID: 1, Mark: "mark1"},
					{WindowID: 2, Mark: "mark2"},
					{WindowID: 3, Mark: "mark3"},
				}, nil,
			).Times(1)

		mockAeroSpaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospace.Window{
			{
				WindowID:    1,
				WindowTitle: "title1",
				AppName:     "app1",
				Workspace:   "ws1",
				AppBundleID: "bundle1",
			},
			{
				WindowID:    2,
				WindowTitle: "title2",
				AppName:     "app2",
				Workspace:   "ws2",
				AppBundleID: "bundle2",
			},
			{
				WindowID:    3,
				WindowTitle: "title3",
				AppName:     "app3",
				Workspace:   "ws3",
				AppBundleID: "bundle3",
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

		args := []string{"list", "-o", "json"}
		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		var jsonResult []map[string]interface{}
		err = json.Unmarshal([]byte(result), &jsonResult)
		require.NoError(t, err)
		assert.Len(t, jsonResult, 3)
		assert.Equal(t, "mark1", jsonResult[0]["mark"])
		assert.Equal(t, "mark2", jsonResult[1]["mark"])
		assert.Equal(t, "mark3", jsonResult[2]["mark"])
	})

	t.Run("CSV output maintains header consistency", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			GetMarks().
			Return(
				[]queries.Mark{
					{WindowID: 1, Mark: "mark1"},
					{WindowID: 2, Mark: "mark2"},
				}, nil,
			).Times(1)

		mockAeroSpaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospace.Window{
			{
				WindowID:    1,
				WindowTitle: "title1",
				AppName:     "app1",
				Workspace:   "ws1",
				AppBundleID: "bundle1",
			},
			{
				WindowID:    2,
				WindowTitle: "title2",
				AppName:     "app2",
				Workspace:   "ws2",
				AppBundleID: "bundle2",
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

		args := []string{"list", "-o", "csv"}
		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		reader := csv.NewReader(strings.NewReader(result))
		records, err := reader.ReadAll()
		require.NoError(t, err)
		assert.Len(t, records, 3) // Header + 2 rows
		expectedHeader := []string{
			"mark",
			"window_id",
			"app_name",
			"window_title",
			"workspace",
			"app_bundle_id",
		}
		assert.Equal(t, expectedHeader, records[0])
		// Verify all rows have same number of columns as header
		for i := 1; i < len(records); i++ {
			assert.Len(
				t,
				records[i],
				len(expectedHeader),
				"Row %d should have same number of columns as header",
				i,
			)
		}
	})

	t.Run("print all marked windows", func(t *testing.T) {
		// t.Skip("Skipping")
		command := "ls"
		args := []string{command, "--help"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)

		_, aerospaceClient := mocks.MockAerospaceConnection(ctrl)

		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		cmdAsString := "\n$aerospace-marks " + strings.Join(args, " ") + "\n"
		snaps.MatchSnapshot(t, cmdAsString, out)
	})

	t.Run("snapshot test for text output format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		marks, err := mocks.LoadMarksFixture(
			"../internal/mocks/fixtures/storage/list-marked-windows.json",
		)
		if err != nil {
			t.Fatal(err)
		}
		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			GetMarks().
			Return(marks, nil).
			Times(1)

		windows, err := mocks.LoadAeroWindowsFixtureRaw(
			"../internal/mocks/fixtures/aerospace/list-windows-all.json",
		)
		if err != nil {
			t.Fatal(err)
		}

		mockAeroSpaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		mockAeroSpaceConnection.EXPECT().
			SendCommand(
				"list-windows",
				[]string{
					"--all",
					"--json",
					"--format",
					"%{window-id} %{window-title} %{app-name} %{app-bundle-id} %{workspace}",
				}).
			Return(&aerospacecli.Response{
				ServerVersion: "1.0",
				StdOut:        windows,
				StdErr:        "",
				ExitCode:      0,
			}, nil).Times(1)

		args := []string{"list", "-o", "text"}
		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		cmdAsString := "\n$aerospace-marks " + strings.Join(args, " ") + "\n"
		snaps.MatchSnapshot(t, "context\n", windows, marks, cmdAsString, result)
	})

	t.Run("snapshot test for JSON output format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		marks, err := mocks.LoadMarksFixture(
			"../internal/mocks/fixtures/storage/list-marked-windows.json",
		)
		if err != nil {
			t.Fatal(err)
		}
		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			GetMarks().
			Return(marks, nil).
			Times(1)

		windows, err := mocks.LoadAeroWindowsFixtureRaw(
			"../internal/mocks/fixtures/aerospace/list-windows-all.json",
		)
		if err != nil {
			t.Fatal(err)
		}

		mockAeroSpaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		mockAeroSpaceConnection.EXPECT().
			SendCommand(
				"list-windows",
				[]string{
					"--all",
					"--json",
					"--format",
					"%{window-id} %{window-title} %{app-name} %{app-bundle-id} %{workspace}",
				}).
			Return(&aerospacecli.Response{
				ServerVersion: "1.0",
				StdOut:        windows,
				StdErr:        "",
				ExitCode:      0,
			}, nil).Times(1)

		args := []string{"list", "-o", "json"}
		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		cmdAsString := "\n$aerospace-marks " + strings.Join(args, " ") + "\n"
		snaps.MatchSnapshot(t, "context\n", windows, marks, cmdAsString, result)
	})

	t.Run("snapshot test for CSV output format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		marks, err := mocks.LoadMarksFixture(
			"../internal/mocks/fixtures/storage/list-marked-windows.json",
		)
		if err != nil {
			t.Fatal(err)
		}
		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			GetMarks().
			Return(marks, nil).
			Times(1)

		windows, err := mocks.LoadAeroWindowsFixtureRaw(
			"../internal/mocks/fixtures/aerospace/list-windows-all.json",
		)
		if err != nil {
			t.Fatal(err)
		}

		mockAeroSpaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		mockAeroSpaceConnection.EXPECT().
			SendCommand(
				"list-windows",
				[]string{
					"--all",
					"--json",
					"--format",
					"%{window-id} %{window-title} %{app-name} %{app-bundle-id} %{workspace}",
				}).
			Return(&aerospacecli.Response{
				ServerVersion: "1.0",
				StdOut:        windows,
				StdErr:        "",
				ExitCode:      0,
			}, nil).Times(1)

		args := []string{"list", "-o", "csv"}
		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		cmdAsString := "\n$aerospace-marks " + strings.Join(args, " ") + "\n"
		snaps.MatchSnapshot(t, "context\n", windows, marks, cmdAsString, result)
	})
}
