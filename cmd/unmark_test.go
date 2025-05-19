package cmd

import (
	"strings"
	"testing"

	"github.com/cristianoliveira/aerospace-marks/internal/mocks"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
	"github.com/cristianoliveira/aerospace-marks/internal/testutils"
	"github.com/gkampitakis/go-snaps/snaps"
	"go.uber.org/mock/gomock"
)

func TestUnmarkCommand(t *testing.T) {
	t.Run("unmarks a mark from a window - `marks unmark mark1`", func(t *testing.T) {
		// t.Skip("Skipping")
		command := "unmark"
		args := []string{command, "mark1"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		storageDbClient, strg := mocks.MockStorageDbClient(ctrl)
		dbResult := mocks.MockStorageDbResult(ctrl, nil, &[]int64{1}[0])
		gomock.InOrder(
			storageDbClient.EXPECT().
				Execute(strings.TrimSpace(`
					DELETE FROM marks WHERE mark = ?
				`),
				"mark1").
				Return(dbResult, nil).
				Times(1),
		)

		out, err := testutils.CmdExecute(NewRootCmd(strg), args...)
		if err != nil {
			t.Fatal(err)
		}

		cmdAsString := "aerospace-marks " + strings.Join(args, " ") + "\n"
		snaps.MatchSnapshot(t, cmdAsString, out)
	})

	t.Run("unmarks all marks from a window - `marks unmark`", func(t *testing.T) {
		// t.Skip("Skipping")
		command := "unmark"
		args := []string{command}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		storageDbClient, strg := mocks.MockStorageDbClient(ctrl)
		dbResult := mocks.MockStorageDbResult(ctrl, nil, &[]int64{2}[0])
		gomock.InOrder(
			storageDbClient.EXPECT().
				Execute(`DELETE FROM marks`).
				Return(dbResult, nil).
				Times(1),
		)

		out, err := testutils.CmdExecute(NewRootCmd(strg), args...)
		if err != nil {
			t.Fatal(err)
		}

		cmdAsString := "aerospace-marks " + strings.Join(args, " ") + "\n"
		snaps.MatchSnapshot(t, cmdAsString, out)
	})

	t.Run("unmarks --help", func(t *testing.T) {
		// t.Skip("Skipping")
		command := "unmark"
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
