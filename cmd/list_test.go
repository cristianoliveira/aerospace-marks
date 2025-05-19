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
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestListCommand(t *testing.T) {
	t.Run("shows only the marked windows", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		storageDbClient, strg := mocks.MockStorageDbClient(ctrl)
		storageDbClient.EXPECT().
			QueryAll(gomock.Any()).
			Return(
				[]storage.Mark{
					{
						WindowID: "1",
						Mark:     "mark1",
					},
					{
						WindowID: "2",
						Mark:     "mark2",
					},
				}, nil,
			)

		mockAeroSpaceConnection, _ := mocks.MockAerospaceConnection(ctrl)
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
		mockAeroSpaceConnection.EXPECT().
			SendCommand("list-windows", []string{"--all", "--json"}).
			Return(
				&aerospacecli.Response{
					ServerVersion: "1.0",
					StdOut:        string(jsonData),
					StdErr:        "",
					ExitCode:      0,
				}, nil).Times(1)

		out, err := testutils.CmdExecute(NewRootCmd(strg), "list")
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

		storageDbClient, strg := mocks.MockStorageDbClient(ctrl)
		storageDbClient.EXPECT().
			QueryAll(gomock.Any()).
			Return(
				[]storage.Mark{
					{
						WindowID: "1",
						Mark:     "mark1",
					},
					{
						WindowID: "2",
						Mark:     "mark2",
					},
				}, nil,
			)

		mockAeroSpaceConnection, _ := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospacecli.Window{
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
			SendCommand("list-windows", []string{"--all", "--json"}).
			Return(
				&aerospacecli.Response{
					ServerVersion: "1.0",
					StdOut:        string(jsonData),
					StdErr:        "",
					ExitCode:      0,
				}, nil).Times(1)

		out, err := testutils.CmdExecute(NewRootCmd(strg), "list")
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		assert.Equal(t, result, "No marked window found")
	})

	t.Run("shows no marks found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		storageDbClient, strg := mocks.MockStorageDbClient(ctrl)
		storageDbClient.EXPECT().
			QueryAll(gomock.Any()).
			Return([]storage.Mark{}, nil)

		out, err := testutils.CmdExecute(NewRootCmd(strg), "list")
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
		storageDbClient, strg := mocks.MockStorageDbClient(ctrl)
		storageDbClient.EXPECT().
			QueryAll(gomock.Any()).
			Return(marks, nil)

		windows, err := mocks.LoadAeroWindowsFixtureRaw(
			"../internal/mocks/fixtures/aerospace/list-windows-all.json",
		)
		if err != nil {
			t.Fatal(err)
		}

		mockAeroSpaceConnection, _ := mocks.MockAerospaceConnection(ctrl)
		mockAeroSpaceConnection.EXPECT().
			SendCommand("list-windows", []string{"--all", "--json"}).
			Return(&aerospacecli.Response{
				ServerVersion: "1.0",
				StdOut:        string(windows),
				StdErr:        "",
				ExitCode:      0,
			}, nil).Times(1)

		out, err := testutils.CmdExecute(NewRootCmd(strg), "list")
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		snaps.MatchSnapshot(t,windows, marks, "result:\n", result)
	})

	t.Run("print all marked windows", func(t *testing.T) {
		// t.Skip("Skipping")
		command := "ls"
		args := []string{command, "--help"}

		connector := storage.MarksDatabaseConnector{}
		conn, err := connector.Connect()
		if err != nil {
			t.Fatal(err)
		}
		strg, err := storage.NewMarkClient(conn)
		if err != nil {
			t.Fatal(err)
		}

		out, err := testutils.CmdExecute(NewRootCmd(strg), args...)
		if err != nil {
			t.Fatal(err)
		}

		cmdAsString := "aerospace-marks " + strings.Join(args, " ") + "\n"
		snaps.MatchSnapshot(t, cmdAsString, out)
	})
}
