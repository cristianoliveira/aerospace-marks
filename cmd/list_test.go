package cmd

import (
	"testing"

	"github.com/cristianoliveira/aerospace-marks/internal/mocks/storage"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
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

		storage.DefaultConnector = mockStorage

		out, err := cmdExecute("ls")
		if err != nil {
			t.Fatal(err)
		}

		t.Fatalf("Output: %s", out)
	})
}
