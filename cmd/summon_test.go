package cmd_test

import (
	"strings"
	"testing"

	"github.com/cristianoliveira/aerospace-marks/cmd"
	"github.com/cristianoliveira/aerospace-marks/internal/logger"
	"github.com/cristianoliveira/aerospace-marks/internal/mocks"
	"github.com/cristianoliveira/aerospace-marks/internal/storage/db/queries"
	"github.com/cristianoliveira/aerospace-marks/internal/testutils"
	"github.com/gkampitakis/go-snaps/snaps"
	"go.uber.org/mock/gomock"

	aerospacecli "github.com/cristianoliveira/aerospace-ipc/pkg/client"
)

func TestSummonCmd(t *testing.T) {
	t.Run("validate missing identifier", func(t *testing.T) {
		command := "summon"
		args := []string{command}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger.SetDefaultLogger(&logger.EmptyLogger{})

		_, strg := mocks.MockStorageDBClient(ctrl)
		_, aerospaceClient := mocks.MockAerospaceConnection(ctrl)

		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if out != "" {
			t.Fatal("output should be empty", out)
		}
		if err == nil {
			t.Fatal(err)
		}

		cmdAsString := "aerospace-marks " + strings.Join(args, " ") + "\n"
		snaps.MatchSnapshot(t, cmdAsString, "Error", err.Error())
	})

	t.Run("validate empty identifier", func(t *testing.T) {
		command := "summon"
		args := []string{command, " "}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger.SetDefaultLogger(&logger.EmptyLogger{})

		_, strg := mocks.MockStorageDBClient(ctrl)
		_, aerospaceClient := mocks.MockAerospaceConnection(ctrl)

		cmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(cmd, args...)
		if out != "" {
			t.Fatal("output should be empty", out)
		}
		if err == nil {
			t.Fatal(err)
		}

		cmdAsString := "aerospace-marks " + strings.Join(args, " ") + "\n"
		snaps.MatchSnapshot(t, cmdAsString, "Error", err.Error())
	})

	t.Run("snapshot test for text output format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger.SetDefaultLogger(&logger.EmptyLogger{})

		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			GetWindowIDByMark("mark1").
			Return(1, nil).
			Times(1)

		mockAeroSpaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		// Mock list-workspaces command (used by GetFocusedWorkspace)
		workspaceJSON := `[{"workspace":"workspace1","is-visible":true,"is-focused":true}]`
		mockAeroSpaceConnection.EXPECT().
			SendCommand("list-workspaces", gomock.Any()).
			Return(
				&aerospacecli.Response{
					ServerVersion: "1.0",
					StdOut:        workspaceJSON,
					StdErr:        "",
					ExitCode:      0,
				}, nil).AnyTimes()
		// Mock move-window-to-workspace command
		mockAeroSpaceConnection.EXPECT().
			SendCommand("move-node-to-workspace", gomock.Any()).
			Return(
				&aerospacecli.Response{
					ServerVersion: "1.0",
					StdOut:        "",
					StdErr:        "",
					ExitCode:      0,
				}, nil).AnyTimes()

		args := []string{"summon", "mark1", "-o", "text"}
		rootCmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(rootCmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		cmdAsString := "aerospace-marks " + strings.Join(args, " ") + "\n"
		marks := []queries.Mark{{WindowID: 1, Mark: "mark1"}}
		snaps.MatchSnapshot(t, marks, cmdAsString, "result:\n", result)
	})

	t.Run("snapshot test for JSON output format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger.SetDefaultLogger(&logger.EmptyLogger{})

		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			GetWindowIDByMark("mark1").
			Return(1, nil).
			Times(1)

		mockAeroSpaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		workspaceJSON := `[{"workspace":"workspace1","is-visible":true,"is-focused":true}]`
		mockAeroSpaceConnection.EXPECT().
			SendCommand("list-workspaces", gomock.Any()).
			Return(
				&aerospacecli.Response{
					ServerVersion: "1.0",
					StdOut:        workspaceJSON,
					StdErr:        "",
					ExitCode:      0,
				}, nil).AnyTimes()
		mockAeroSpaceConnection.EXPECT().
			SendCommand("move-node-to-workspace", gomock.Any()).
			Return(
				&aerospacecli.Response{
					ServerVersion: "1.0",
					StdOut:        "",
					StdErr:        "",
					ExitCode:      0,
				}, nil).AnyTimes()

		args := []string{"summon", "mark1", "-o", "json"}
		rootCmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(rootCmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		cmdAsString := "aerospace-marks " + strings.Join(args, " ") + "\n"
		marks := []queries.Mark{{WindowID: 1, Mark: "mark1"}}
		snaps.MatchSnapshot(t, marks, cmdAsString, "result:\n", result)
	})

	t.Run("snapshot test for CSV output format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger.SetDefaultLogger(&logger.EmptyLogger{})

		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			GetWindowIDByMark("mark1").
			Return(1, nil).
			Times(1)

		mockAeroSpaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		workspaceJSON := `[{"workspace":"workspace1","is-visible":true,"is-focused":true}]`
		mockAeroSpaceConnection.EXPECT().
			SendCommand("list-workspaces", gomock.Any()).
			Return(
				&aerospacecli.Response{
					ServerVersion: "1.0",
					StdOut:        workspaceJSON,
					StdErr:        "",
					ExitCode:      0,
				}, nil).AnyTimes()
		mockAeroSpaceConnection.EXPECT().
			SendCommand("move-node-to-workspace", gomock.Any()).
			Return(
				&aerospacecli.Response{
					ServerVersion: "1.0",
					StdOut:        "",
					StdErr:        "",
					ExitCode:      0,
				}, nil).AnyTimes()

		args := []string{"summon", "mark1", "-o", "csv"}
		rootCmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(rootCmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		cmdAsString := "aerospace-marks " + strings.Join(args, " ") + "\n"
		marks := []queries.Mark{{WindowID: 1, Mark: "mark1"}}
		snaps.MatchSnapshot(t, marks, cmdAsString, "result:\n", result)
	})

	t.Run("snapshot test for JSON output format with focus flag", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger.SetDefaultLogger(&logger.EmptyLogger{})

		_, strg := mocks.MockStorageDBClient(ctrl)
		strg.EXPECT().
			GetWindowIDByMark("mark1").
			Return(1, nil).
			Times(1)

		mockAeroSpaceConnection, aerospaceClient := mocks.MockAerospaceConnection(ctrl)
		workspaceJSON := `[{"workspace":"workspace1","is-visible":true,"is-focused":true}]`
		mockAeroSpaceConnection.EXPECT().
			SendCommand("list-workspaces", gomock.Any()).
			Return(
				&aerospacecli.Response{
					ServerVersion: "1.0",
					StdOut:        workspaceJSON,
					StdErr:        "",
					ExitCode:      0,
				}, nil).AnyTimes()
		mockAeroSpaceConnection.EXPECT().
			SendCommand("move-node-to-workspace", gomock.Any()).
			Return(
				&aerospacecli.Response{
					ServerVersion: "1.0",
					StdOut:        "",
					StdErr:        "",
					ExitCode:      0,
				}, nil).AnyTimes()
		mockAeroSpaceConnection.EXPECT().
			SendCommand("focus", gomock.Any()).
			Return(
				&aerospacecli.Response{
					ServerVersion: "1.0",
					StdOut:        "",
					StdErr:        "",
					ExitCode:      0,
				}, nil).AnyTimes()

		args := []string{"summon", "mark1", "--focus", "-o", "json"}
		rootCmd := cmd.NewRootCmd(strg, aerospaceClient)
		out, err := testutils.CmdExecute(rootCmd, args...)
		if err != nil {
			t.Fatal(err)
		}

		result := strings.TrimSpace(out)
		cmdAsString := "aerospace-marks " + strings.Join(args, " ") + "\n"
		marks := []queries.Mark{{WindowID: 1, Mark: "mark1"}}
		snaps.MatchSnapshot(t, marks, cmdAsString, "result:\n", result)
	})
}
