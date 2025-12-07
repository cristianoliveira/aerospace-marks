package format_test

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"strings"
	"testing"

	"github.com/cristianoliveira/aerospace-marks/internal/format"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOutputEventFormatter(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{"valid text", "text", false},
		{"valid json", "json", false},
		{"valid csv", "csv", false},
		{"case insensitive", "JSON", false},
		{"with spaces", "  text  ", false},
		{"invalid format", "invalid", true},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			formatter, err := format.NewOutputEventFormatter(&buf, tt.format)
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, formatter)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, formatter)
			}
		})
	}
}

func TestOutputEventFormatter_FormatText_GetCommand(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := format.NewOutputEventFormatter(&buf, "text")
	require.NoError(t, err)

	event := format.OutputEvent{
		Command:  "get",
		WindowID: 123,
		AppName:  "Test App",
		Message:  "Test Title",
		Result:   "123 | Test App | Test Title",
	}

	err = formatter.Format(event)
	require.NoError(t, err)

	result := strings.TrimSpace(buf.String())
	// Format: <window_id> | <app_name> | <window_title>
	assert.Equal(t, "123 | Test App | Test Title", result)
}

func TestOutputEventFormatter_FormatText_SingleField(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := format.NewOutputEventFormatter(&buf, "text")
	require.NoError(t, err)

	event := format.OutputEvent{
		Command: "get",
		Result:  "123",
	}

	err = formatter.Format(event)
	require.NoError(t, err)

	result := strings.TrimSpace(buf.String())
	// Single field output should just be the value
	assert.Equal(t, "123", result)
}

func TestOutputEventFormatter_FormatText_FocusCommand(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := format.NewOutputEventFormatter(&buf, "text")
	require.NoError(t, err)

	event := format.OutputEvent{
		Command:  "focus",
		Action:   "focus",
		WindowID: 456,
		Message:  "Focus moved to window ID 456",
		Result:   "success",
	}

	err = formatter.Format(event)
	require.NoError(t, err)

	result := strings.TrimSpace(buf.String())
	// Focus command should output just the message
	assert.Equal(t, "Focus moved to window ID 456", result)
}

func TestOutputEventFormatter_FormatJSON(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := format.NewOutputEventFormatter(&buf, "json")
	require.NoError(t, err)

	event := format.OutputEvent{
		Command:   "get",
		WindowID:  789,
		AppName:   "JSON App",
		Message:   "JSON Title",
		Workspace: "json-ws",
		Result:    "789 | JSON App | JSON Title",
	}

	err = formatter.Format(event)
	require.NoError(t, err)

	result := strings.TrimSpace(buf.String())
	var jsonResult format.OutputEvent
	err = json.Unmarshal([]byte(result), &jsonResult)
	require.NoError(t, err)
	assert.Equal(t, "get", jsonResult.Command)
	assert.Equal(t, 789, jsonResult.WindowID)
	assert.Equal(t, "JSON App", jsonResult.AppName)
	assert.Equal(t, "JSON Title", jsonResult.Message)
}

func TestOutputEventFormatter_FormatCSV(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := format.NewOutputEventFormatter(&buf, "csv")
	require.NoError(t, err)

	event := format.OutputEvent{
		Command:   "get",
		WindowID:  999,
		AppName:   "CSV App",
		Message:   "CSV Title",
		Workspace: "csv-ws",
		Result:    "999 | CSV App | CSV Title",
	}

	err = formatter.Format(event)
	require.NoError(t, err)

	result := buf.String()
	reader := csv.NewReader(strings.NewReader(result))
	records, err := reader.ReadAll()
	require.NoError(t, err)
	assert.Len(t, records, 2) // Header + 1 row
	expectedHeader := []string{
		"command",
		"action",
		"window_id",
		"app_name",
		"workspace",
		"target_workspace",
		"result",
		"message",
	}
	assert.Equal(t, expectedHeader, records[0])
	assert.Equal(
		t,
		[]string{
			"get",
			"",
			"999",
			"CSV App",
			"csv-ws",
			"",
			"999 | CSV App | CSV Title",
			"CSV Title",
		},
		records[1],
	)
}
