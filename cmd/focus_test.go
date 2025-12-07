package cmd_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/cristianoliveira/aerospace-marks/cmd"
	"github.com/cristianoliveira/aerospace-marks/internal/logger"
	"github.com/cristianoliveira/aerospace-marks/internal/mocks"
	"github.com/cristianoliveira/aerospace-marks/internal/stdout"
	"github.com/cristianoliveira/aerospace-marks/internal/testutils"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	aerospacecli "github.com/cristianoliveira/aerospace-ipc/pkg/client"
)

func TestFocusCmd(t *testing.T) {
	t.Run("validate missing identifier", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger.SetDefaultLogger(&logger.EmptyLogger{})

		_, strg := mocks.MockStorageDBClient(ctrl)

		_, aerospaceClient := mocks.MockAerospaceConnection(ctrl)

		args := []string{"focus"}
		rootCmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(rootCmd, args...)
		if out != "" {
			t.Fatal("output should be empty", out)
		}
		if err == nil {
			t.Fatal(err)
		}

		snapshot := testutils.RenderSnapshotSpec(testutils.SnapshotSpec{
			Command: testutils.CommandString(args...),
			Stdout:  out,
			Stderr:  err.Error(),
		})
		snaps.MatchSnapshot(t, snapshot)
	})

	t.Run("focus to a window by mark - `marks focus mark1`", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			GetWindowIDByMark("mark1").
			Return(1, nil).
			Times(1)

		mockAeroSpaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		mockAeroSpaceConnection.EXPECT().
			SendCommand("focus", []string{"--window-id", "1"}).
			Return(
				&aerospacecli.Response{
					ServerVersion: "1.0",
					StdOut:        "",
					StdErr:        "",
					ExitCode:      0,
				}, nil).Times(1)

		args := []string{"focus", "mark1"}
		rootCmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(rootCmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		snapshot := testutils.RenderSnapshotSpec(testutils.SnapshotSpec{
			Command: testutils.CommandString(args...),
			Stdout:  out,
		})
		snaps.MatchSnapshot(t, snapshot)
	})

	t.Run("focus using mark that does not exist", func(t *testing.T) {
		command := "focus"
		args := []string{command, "nonexistent-mark"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, strg := mocks.MockStorageDBClient(ctrl)

		strg.EXPECT().
			GetWindowIDByMark("nonexistent-mark").
			Return(0, nil).
			Times(1)

		_, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		//nolint:reassign // Test utility needs to modify package variable
		stdout.ShouldExit = false

		rootCmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(rootCmd, args...)
		if err == nil {
			t.Fatal(err)
		}
		if out != "" {
			t.Fatal("output should be empty", out)
		}

		snapshot := testutils.RenderSnapshotSpec(testutils.SnapshotSpec{
			Command: testutils.CommandString(args...),
			Stdout:  out,
			Stderr:  err.Error(),
		})
		snaps.MatchSnapshot(t, snapshot)
	})

	t.Run("outputs JSON format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			GetWindowIDByMark("mark1").
			Return(1, nil).
			Times(1)

		mockAeroSpaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		mockAeroSpaceConnection.EXPECT().
			SendCommand("focus", []string{"--window-id", "1"}).
			Return(
				&aerospacecli.Response{
					ServerVersion: "1.0",
					StdOut:        "",
					StdErr:        "",
					ExitCode:      0,
				}, nil).Times(1)

		args := []string{"focus", "mark1", "-o", "json"}
		rootCmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(rootCmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		// Should be valid JSON OutputEvent
		var jsonResult map[string]interface{}
		err = json.Unmarshal([]byte(result), &jsonResult)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "focus", jsonResult["command"])
		assert.Equal(t, "focus", jsonResult["action"])
		assert.InDelta(t, 1.0, jsonResult["window_id"], 0.0)
		assert.Contains(t, jsonResult["message"], "Focus moved")
	})

	t.Run("outputs CSV format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			GetWindowIDByMark("mark1").
			Return(1, nil).
			Times(1)

		mockAeroSpaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		mockAeroSpaceConnection.EXPECT().
			SendCommand("focus", []string{"--window-id", "1"}).
			Return(
				&aerospacecli.Response{
					ServerVersion: "1.0",
					StdOut:        "",
					StdErr:        "",
					ExitCode:      0,
				}, nil).Times(1)

		args := []string{"focus", "mark1", "-o", "csv"}
		rootCmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(rootCmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		lines := strings.Split(result, "\n")
		assert.Len(t, lines, 2)
		assert.Equal(
			t,
			"command,action,window_id,app_name,workspace,target_workspace,result,message",
			lines[0],
		)
		assert.Contains(t, lines[1], "focus,focus,1,")
		assert.Contains(t, lines[1], ",success,")
		assert.Contains(t, lines[1], "Focus moved")
	})
}
