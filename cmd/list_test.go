package cmd

import (
	"encoding/json"
	"strings"
	"testing"

	aerospace "github.com/cristianoliveira/aerospace-ipc"
	aerospacecli "github.com/cristianoliveira/aerospace-ipc/pkg/client"
	"github.com/cristianoliveira/aerospace-marks/internal/mocks"
	"github.com/cristianoliveira/aerospace-marks/internal/storage/db/queries"
	"github.com/cristianoliveira/aerospace-marks/internal/testutils"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestListCommand(t *testing.T) {
	t.Run("shows only the marked windows", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDbClient(ctrl)
		strg.EXPECT().
			GetMarks().
			Return(
				[]queries.Mark{
					{
						WindowID: "1",
						Mark:     "mark1",
					},
					{
						WindowID: "2",
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
		cmd := NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		lines := strings.Split(result, "\n")
		assert.Equal(t, 2, len(lines))
		assert.Equal(t, lines, []string{
			"mark1 | 1 | app1 | title1",
			"mark2 | 2 | app2 | title2",
		})
	})

	t.Run("shows no marked window found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDbClient(ctrl)
		strg.EXPECT().
			GetMarks().
			Return(
				[]queries.Mark{
					{
						WindowID: "1",
						Mark:     "mark1",
					},
					{
						WindowID: "2",
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
			},
			{
				WindowID:    222,
				WindowTitle: "title2",
				AppName:     "app2",
			},
			{
				WindowID:    333,
				WindowTitle: "title3",
				AppName:     "app3",
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
		cmd := NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		assert.Equal(t, result, "No marked window found")
	})

	t.Run("shows no marks found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDbClient(ctrl)
		strg.EXPECT().
			GetMarks().
			Return([]queries.Mark{}, nil).
			Times(1)

		_, aerospaceClient := mocks.MockAerospaceConnection(ctrl)

		args := []string{"list"}
		cmd := NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		assert.Equal(t, result, "No marks found")
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
		_, strg := mocks.MockStorageDbClient(ctrl)
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
				StdOut:        string(windows),
				StdErr:        "",
				ExitCode:      0,
			}, nil).Times(1)

		args := []string{"list"}
		cmd := NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		snaps.MatchSnapshot(t, windows, marks, "result:\n", result)
	})

	t.Run("print all marked windows", func(t *testing.T) {
		// t.Skip("Skipping")
		command := "ls"
		args := []string{command, "--help"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDbClient(ctrl)

		_, aerospaceClient := mocks.MockAerospaceConnection(ctrl)

		cmd := NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		cmdAsString := "aerospace-marks " + strings.Join(args, " ") + "\n"
		snaps.MatchSnapshot(t, cmdAsString, out)
	})
}
