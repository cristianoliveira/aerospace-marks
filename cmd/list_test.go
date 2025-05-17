package cmd

import (
	"strings"
	"testing"

	"github.com/cristianoliveira/aerospace-marks/internal/mocks/aerospacecli"
	storage_mock "github.com/cristianoliveira/aerospace-marks/internal/mocks/storage"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
	"github.com/cristianoliveira/aerospace-marks/internal/testutils"
	"github.com/cristianoliveira/aerospace-marks/pkgs/aerospacecli"
	"go.uber.org/mock/gomock"
)

func TestOutputSetter(t *testing.T) {
	t.Run("SetOutputPosition", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		storageDbClient := storage_mock.NewMockStorageDbClient(ctrl)
		storageDbClient.EXPECT().Close().Return(nil).Times(1)
		storageDbClient.EXPECT().QueryAll(gomock.Any()).Return(
			[]storage.Mark{
				{
					WindowID: "1",
					Mark:    "mark1",
				},
			}, nil,
		)

		mockStorage := storage_mock.NewMockDatabaseConnector(ctrl)
		mockStorage.EXPECT().Connect().Return(storageDbClient, nil).Times(1)

		mockAeroSpaceConnection := aerospacecli_mock.NewMockAeroSpaceSocketConn(ctrl)
		mockAeroSpaceConnetor := aerospacecli_mock.NewMockAeroSpaceConnector(ctrl)
		mockAeroSpaceConnetor.EXPECT().Connect().Return(mockAeroSpaceConnection, nil).Times(1)
		mockAeroSpaceConnection.EXPECT().SendCommand("list-windows", []string{"--all"}).Return(&aerospacecli.Response{
			ServerVersion: "1.0",
			StdOut:        "1 | app1 | title1",
			StdErr:        "",
			ExitCode:      0,
		}, nil).Times(1)
		mockAeroSpaceConnection.EXPECT().CloseConnection().Return(nil).Times(1)

		storage.DefaultConnector = mockStorage
		aerospacecli.DefaultConnector = mockAeroSpaceConnetor

		out, err := testutils.CmdExecute(rootCmd, "list")
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		expectedOutput := "mark1| 1 | app1 | title1"
		if result != expectedOutput {
			t.Errorf("\r\nexpected %s \n     got  %s", expectedOutput, out)
		}
	})
}
