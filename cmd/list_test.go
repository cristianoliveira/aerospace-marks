package cmd

import (
	"strings"
	"testing"

	"github.com/cristianoliveira/aerospace-marks/internal/mocks"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
	"github.com/cristianoliveira/aerospace-marks/internal/testutils"
	"github.com/cristianoliveira/aerospace-marks/pkgs/aerospacecli"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestOutputSetter(t *testing.T) {
	t.Run("SetOutputPosition", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		storageDbClient, _ := mocks.MockStorageDbClient(ctrl)
		storageDbClient.EXPECT().
			QueryAll(gomock.Any()).
				Return(
					[]storage.Mark{
						{
							WindowID: "1",
							Mark:     "mark1",
						},
					}, nil,
				)

		mockAeroSpaceConnection, _ := mocks.MockAerospaceConnection(ctrl)
		mockAeroSpaceConnection.EXPECT().
			SendCommand("list-windows", []string{"--all"}).
			Return(
				&aerospacecli.Response{
					ServerVersion: "1.0",
					StdOut:        "1 | app1 | title1",
					StdErr:        "",
					ExitCode:      0,
				}, nil).Times(1)

		out, err := testutils.CmdExecute(rootCmd, "list")
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		assert.Equal(t, "mark1| 1 | app1 | title1", result)
	})
}
