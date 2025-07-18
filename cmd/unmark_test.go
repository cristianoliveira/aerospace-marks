package cmd

import (
	"fmt"
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

		_, aerospaceClient := mocks.MockAerospaceConnection(ctrl)

		cmd := NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
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

		aerospaceClient := &testutils.MockEmptyAerspaceMarkWindows{}

		cmd := NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
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
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		storageDbClient, _ := mocks.MockStorageDbClient(ctrl)
		strg, err := storage.NewMarkClient(storageDbClient)
		if err != nil {
			t.Fatal(err)
		}

		aerospaceClient := &testutils.MockEmptyAerspaceMarkWindows{}

		cmd := NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		cmdAsString := "aerospace-marks " + strings.Join(args, " ") + "\n"
		snaps.MatchSnapshot(t, cmdAsString, out)
	})

	t.Run("fails when mark not found", func(t *testing.T) {
		// t.Skip("Skipping")
		command := "unmark"
		args := []string{command, "unkown"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		storageDbClient, strg := mocks.MockStorageDbClient(ctrl)
		dbResult := mocks.MockStorageDbResult(ctrl, nil, &[]int64{0}[0])
		gomock.InOrder(
			storageDbClient.EXPECT().
				Execute(strings.TrimSpace(`
					DELETE FROM marks WHERE mark = ?
				`),
					"unkown").
				Return(dbResult, nil).
				Times(1),
		)

		_, aerospaceClient := mocks.MockAerospaceConnection(ctrl)

		cmd := NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if err == nil {
			t.Fatal(err)
		}

		cmdAsString := "aerospace-marks " + strings.Join(args, " ") + "\n"
		expectedError := fmt.Sprintf("Error\n%+v", err)
		snaps.MatchSnapshot(t, cmdAsString, expectedError, out)
	})
}
