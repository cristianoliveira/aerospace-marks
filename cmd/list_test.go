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

func TestListCommand(t *testing.T) {
	t.Run("shows only the marked windows", func(t *testing.T) {
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
					{
						WindowID: "2",
						Mark:     "mark2",
					},
				}, nil,
			)

		mockAeroSpaceConnection, _ := mocks.MockAerospaceConnection(ctrl)
		aerospaceResponseOutput := []string{
			"1 | app1 | title1",
			"2 | app2 | title2",
			"3 | app3 | title3",
		}
		mockAeroSpaceConnection.EXPECT().
			SendCommand("list-windows", []string{"--all"}).
			Return(
				&aerospacecli.Response{
					ServerVersion: "1.0",
					StdOut:        strings.Join(aerospaceResponseOutput, "\n"),
					StdErr:        "",
					ExitCode:      0,
				}, nil).Times(1)

		out, err := testutils.CmdExecute(rootCmd, "list")
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		lines := strings.Split(result, "\n")
		assert.Equal(t, 2, len(lines))
		assert.Equal(t, lines, []string{
			"mark1| 1 | app1 | title1",
			"mark2| 2 | app2 | title2",
		})
	})
}
