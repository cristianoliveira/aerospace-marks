package cmd_test

import (
	"strings"
	"testing"

	"github.com/cristianoliveira/aerospace-marks/cmd"
	"github.com/cristianoliveira/aerospace-marks/internal/logger"
	"github.com/cristianoliveira/aerospace-marks/internal/mocks"
	"github.com/cristianoliveira/aerospace-marks/internal/testutils"
	"github.com/gkampitakis/go-snaps/snaps"
	"go.uber.org/mock/gomock"
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
}
