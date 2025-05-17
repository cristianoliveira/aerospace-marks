package cmd

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/cristianoliveira/aerospace-marks/internal/mocks"
	"github.com/cristianoliveira/aerospace-marks/internal/testutils"
	"github.com/cristianoliveira/aerospace-marks/pkgs/aerospacecli"
	"github.com/gkampitakis/go-snaps/snaps"
	"go.uber.org/mock/gomock"
)

func TestMarkCommand(t *testing.T) {
	t.Run("marks the focused window - `marks mark mark1`", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		storageDbClient, _ := mocks.MockStorageDbClient(ctrl)
		dbResult := mocks.MockStorageDbResult(ctrl, nil, &[]int64{1}[0])
		gomock.InOrder(
			storageDbClient.EXPECT().
				Execute(strings.TrimSpace(`
					DELETE FROM marks WHERE mark = ?
				`),
				"mark1").
				Return(dbResult, nil).
				Times(1),
			storageDbClient.EXPECT().
				Execute(strings.TrimSpace(`
					INSERT INTO marks (window_id, mark) VALUES (?, ?)
				`), "1", "mark1").
				Return(dbResult, nil).
				Times(1),
		)

		mockAeroSpaceConnection, _ := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospacecli.Window{
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
			SendCommand("list-windows", []string{"--focused", "--json"}).
			Return(
				&aerospacecli.Response{
					ServerVersion: "1.0",
					StdOut:        string(jsonData),
					StdErr:        "",
					ExitCode:      0,
				}, nil).Times(1)

		args := []string{"mark", "mark1"}
		out, err := testutils.CmdExecute(rootCmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		snaps.MatchSnapshot(t,windows, args, "result:\n", result)
	})

	t.Run("marks the focused window - `marks mark --add`", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		storageDbClient, _ := mocks.MockStorageDbClient(ctrl)
		gomock.InOrder(
			storageDbClient.EXPECT().
				Execute(strings.TrimSpace(`
					INSERT INTO marks (window_id, mark) VALUES (?, ?)
				`), "1", "mark1").
				Times(1),
		)

		mockAeroSpaceConnection, _ := mocks.MockAerospaceConnection(ctrl)
		windows := []aerospacecli.Window{
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
			SendCommand("list-windows", []string{"--focused", "--json"}).
			Return(
				&aerospacecli.Response{
					ServerVersion: "1.0",
					StdOut:        string(jsonData),
					StdErr:        "",
					ExitCode:      0,
				}, nil).Times(1)

		args := []string{"mark", "mark1", "--add"}
		out, err := testutils.CmdExecute(rootCmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		snaps.MatchSnapshot(t,windows, args, "result:\n", result)
	})
}
