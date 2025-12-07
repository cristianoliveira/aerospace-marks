package cmd_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/cristianoliveira/aerospace-marks/cmd"
	"github.com/cristianoliveira/aerospace-marks/internal/mocks"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
	"github.com/cristianoliveira/aerospace-marks/internal/testutils"
	"github.com/gkampitakis/go-snaps/snaps"
	"go.uber.org/mock/gomock"
)

func TestInfoCmd(t *testing.T) {
	t.Run("Happy path - all compatible", func(tt *testing.T) {
		ctrl := gomock.NewController(tt)
		defer ctrl.Finish()

		aerospaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		dbClient, storageClient := mocks.MockStorageDBClient(ctrl)

		storageClient.
			EXPECT().
			Client().
			Return(dbClient).
			Times(1)

		dbClient.
			EXPECT().
			GetStorageConfig().
			Return(storage.StorageConfig{
				DBPath: "/tmp/database/",
				DBName: "foo.db",
			}).
			Times(1)

		gomock.InOrder(
			aerospaceConnection.
				EXPECT().
				GetSocketPath().
				Return("/tmp/foo.sock", nil).
				Times(1),

			aerospaceConnection.
				EXPECT().
				CheckServerVersion().
				Return(nil).
				Times(1),

			aerospaceConnection.
				EXPECT().
				GetServerVersion().
				Return("aerospace-ipc v0.1.0", nil).
				Times(1),
		)

		cmd := cmd.InfoCmd(
			storageClient,
			aerospaceClient,
		)
		out, err := testutils.CmdExecute(cmd)
		if err != nil {
			tt.Fatal(err)
		}

		snapshot := testutils.RenderSnapshotSpec(testutils.SnapshotSpec{
			Command: fmt.Sprintf("aerospace-marks %s", cmd.Use),
			Stdout:  out,
		})
		snaps.MatchSnapshot(tt, snapshot)
	})

	t.Run("Happy path - non compatible", func(tt *testing.T) {
		ctrl := gomock.NewController(tt)
		defer ctrl.Finish()

		aerospaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		dbClient, storageClient := mocks.MockStorageDBClient(ctrl)

		storageClient.
			EXPECT().
			Client().
			Return(dbClient).
			Times(1)

		dbClient.
			EXPECT().
			GetStorageConfig().
			Return(storage.StorageConfig{
				DBPath: "/tmp/database/",
				DBName: "foo.db",
			}).
			Times(1)

		gomock.InOrder(
			aerospaceConnection.
				EXPECT().
				GetSocketPath().
				Return("/tmp/foo.sock", nil).
				Times(1),

			aerospaceConnection.
				EXPECT().
				CheckServerVersion().
				Return(errors.New("incompatible version because reasons")).
				Times(1),

			aerospaceConnection.
				EXPECT().
				GetServerVersion().
				Return("aerospace-ipc v3.1.0", nil).
				Times(1),
		)

		cmd := cmd.InfoCmd(
			storageClient,
			aerospaceClient,
		)
		out, err := testutils.CmdExecute(cmd)
		if err != nil {
			tt.Fatal(err)
		}

		snapshot := testutils.RenderSnapshotSpec(testutils.SnapshotSpec{
			Command: fmt.Sprintf("aerospace-marks %s", cmd.Use),
			Stdout:  out,
		})
		snaps.MatchSnapshot(tt, snapshot)
	})

	t.Run("Failure - to retrieve socket path", func(tt *testing.T) {
		ctrl := gomock.NewController(tt)
		defer ctrl.Finish()

		aerospaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		dbClient, storageClient := mocks.MockStorageDBClient(ctrl)

		storageClient.
			EXPECT().
			Client().
			Return(dbClient).
			Times(1)

		dbClient.
			EXPECT().
			GetStorageConfig().
			Return(storage.StorageConfig{
				DBPath: "/tmp/database/",
				DBName: "foo.db",
			}).
			Times(1)

		gomock.InOrder(
			aerospaceConnection.
				EXPECT().
				GetSocketPath().
				Return("", errors.New("missing socket path")).
				Times(1),
		)

		cmd := cmd.InfoCmd(
			storageClient,
			aerospaceClient,
		)
		out, err := testutils.CmdExecute(cmd)
		if err == nil {
			tt.Fatal(err)
		}

		snapshot := testutils.RenderSnapshotSpec(testutils.SnapshotSpec{
			Command: fmt.Sprintf("aerospace-marks %s", cmd.Use),
			Stdout:  out,
			Stderr:  err.Error(),
		})
		snaps.MatchSnapshot(tt, snapshot)
	})
}
