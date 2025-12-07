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

func TestNewListOutputFormatter(t *testing.T) {
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
			formatter, err := format.NewListOutputFormatter(&buf, tt.format)
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

func TestListOutputFormatter_FormatText(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := format.NewListOutputFormatter(&buf, "text")
	require.NoError(t, err)

	windows := []format.MarkedWindow{
		{
			Mark:        "mark1",
			WindowID:    1,
			AppName:     "app1",
			WindowTitle: "title1",
			Workspace:   "workspace1",
			AppBundleID: "bundle1",
		},
		{
			Mark:        "mark2",
			WindowID:    22,
			AppName:     "app2 super long",
			WindowTitle: "title2",
			Workspace:   "",
			AppBundleID: "",
		},
	}

	err = formatter.Format(windows)
	require.NoError(t, err)

	result := buf.String()
	lines := strings.Split(strings.TrimSpace(result), "\n")
	assert.Len(t, lines, 2)
	assert.Contains(t, lines[0], "mark1")
	assert.Contains(t, lines[0], "|")
	assert.Contains(t, lines[1], "mark2")
	// Check that empty values are replaced with underscore
	assert.Contains(t, lines[1], "_")
}

func TestListOutputFormatter_FormatText_Empty(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := format.NewListOutputFormatter(&buf, "text")
	require.NoError(t, err)

	err = formatter.Format([]format.MarkedWindow{})
	require.NoError(t, err)
	// Empty should produce no output for text format
	assert.Empty(t, strings.TrimSpace(buf.String()))
}

func TestListOutputFormatter_FormatJSON(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := format.NewListOutputFormatter(&buf, "json")
	require.NoError(t, err)

	windows := []format.MarkedWindow{
		{
			Mark:        "mark1",
			WindowID:    1,
			AppName:     "app1",
			WindowTitle: "title1",
			Workspace:   "workspace1",
			AppBundleID: "bundle1",
		},
		{
			Mark:        "mark2",
			WindowID:    2,
			AppName:     "app2",
			WindowTitle: "",
			Workspace:   "",
			AppBundleID: "",
		},
	}

	err = formatter.Format(windows)
	require.NoError(t, err)

	result := strings.TrimSpace(buf.String())
	var jsonResult []format.MarkedWindow
	err = json.Unmarshal([]byte(result), &jsonResult)
	require.NoError(t, err)
	assert.Len(t, jsonResult, 2)
	assert.Equal(t, "mark1", jsonResult[0].Mark)
	assert.Equal(t, 1, jsonResult[0].WindowID)
	assert.Empty(t, jsonResult[1].WindowTitle) // Empty strings preserved
}

func TestListOutputFormatter_FormatJSON_Empty(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := format.NewListOutputFormatter(&buf, "json")
	require.NoError(t, err)

	err = formatter.Format([]format.MarkedWindow{})
	require.NoError(t, err)

	result := strings.TrimSpace(buf.String())
	var jsonResult []format.MarkedWindow
	err = json.Unmarshal([]byte(result), &jsonResult)
	require.NoError(t, err)
	assert.Empty(t, jsonResult)
	assert.Equal(t, "[]", result)
}

func TestListOutputFormatter_FormatCSV(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := format.NewListOutputFormatter(&buf, "csv")
	require.NoError(t, err)

	windows := []format.MarkedWindow{
		{
			Mark:        "mark1",
			WindowID:    1,
			AppName:     "app1",
			WindowTitle: "title1",
			Workspace:   "workspace1",
			AppBundleID: "bundle1",
		},
		{
			Mark:        "mark2",
			WindowID:    2,
			AppName:     "app2",
			WindowTitle: "",
			Workspace:   "",
			AppBundleID: "",
		},
	}

	err = formatter.Format(windows)
	require.NoError(t, err)

	result := buf.String()
	reader := csv.NewReader(strings.NewReader(result))
	records, err := reader.ReadAll()
	require.NoError(t, err)
	assert.Len(t, records, 3) // Header + 2 rows
	assert.Equal(
		t,
		[]string{"mark", "window_id", "app_name", "window_title", "workspace", "app_bundle_id"},
		records[0],
	)
	assert.Equal(t, []string{"mark1", "1", "app1", "title1", "workspace1", "bundle1"}, records[1])
	assert.Equal(t, []string{"mark2", "2", "app2", "", "", ""}, records[2])
}

func TestListOutputFormatter_FormatCSV_Empty(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := format.NewListOutputFormatter(&buf, "csv")
	require.NoError(t, err)

	err = formatter.Format([]format.MarkedWindow{})
	require.NoError(t, err)

	result := buf.String()
	reader := csv.NewReader(strings.NewReader(result))
	records, err := reader.ReadAll()
	require.NoError(t, err)
	assert.Len(t, records, 1) // Header only
	assert.Equal(
		t,
		[]string{"mark", "window_id", "app_name", "window_title", "workspace", "app_bundle_id"},
		records[0],
	)
}

func TestListOutputFormatter_EmptyToUnderscore(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := format.NewListOutputFormatter(&buf, "text")
	require.NoError(t, err)

	windows := []format.MarkedWindow{
		{
			Mark:        "mark1",
			WindowID:    1,
			AppName:     "",
			WindowTitle: "",
			Workspace:   "",
			AppBundleID: "",
		},
	}

	err = formatter.Format(windows)
	require.NoError(t, err)

	result := buf.String()
	// All empty fields should be replaced with underscore
	assert.Contains(t, result, "_")
	// Count underscores (should be 4 for empty fields)
	underscoreCount := strings.Count(result, "_")
	assert.GreaterOrEqual(t, underscoreCount, 4)
}

func TestListOutputFormatter_FormatText_ColumnAlignment(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := format.NewListOutputFormatter(&buf, "text")
	require.NoError(t, err)

	windows := []format.MarkedWindow{
		{
			Mark:        "a",
			WindowID:    1,
			AppName:     "short",
			WindowTitle: "title",
			Workspace:   "ws",
			AppBundleID: "bundle",
		},
		{
			Mark:        "very-long-mark-name",
			WindowID:    99999,
			AppName:     "very long app name here",
			WindowTitle: "very long window title",
			Workspace:   "very-long-workspace-name",
			AppBundleID: "very.long.bundle.id",
		},
	}

	err = formatter.Format(windows)
	require.NoError(t, err)

	result := buf.String()
	lines := strings.Split(strings.TrimSpace(result), "\n")
	assert.Len(t, lines, 2)

	// Verify columns are aligned (all lines should have same structure)
	fields1 := strings.Split(lines[0], "|")
	fields2 := strings.Split(lines[1], "|")
	assert.Len(t, fields1, 6)
	assert.Len(t, fields2, 6)

	// Verify pipe separators exist
	assert.Contains(t, lines[0], "|")
	assert.Contains(t, lines[1], "|")
}

func TestListOutputFormatter_FormatText_SpecialCharacters(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := format.NewListOutputFormatter(&buf, "text")
	require.NoError(t, err)

	windows := []format.MarkedWindow{
		{
			Mark:        "mark|with|pipes",
			WindowID:    1,
			AppName:     "app with spaces",
			WindowTitle: "title\nwith\nnewlines",
			Workspace:   "workspace",
			AppBundleID: "bundle.id",
		},
	}

	err = formatter.Format(windows)
	require.NoError(t, err)

	result := buf.String()
	// Should still format correctly even with special characters
	assert.Contains(t, result, "mark|with|pipes")
	assert.Contains(t, result, "app with spaces")
}

func TestListOutputFormatter_FormatJSON_SpecialCharacters(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := format.NewListOutputFormatter(&buf, "json")
	require.NoError(t, err)

	windows := []format.MarkedWindow{
		{
			Mark:        "mark\"with\"quotes",
			WindowID:    1,
			AppName:     "app\nwith\nnewlines",
			WindowTitle: "title\twith\ttabs",
			Workspace:   "workspace",
			AppBundleID: "bundle.id",
		},
	}

	err = formatter.Format(windows)
	require.NoError(t, err)

	result := strings.TrimSpace(buf.String())
	// Should be valid JSON even with special characters
	var jsonResult []format.MarkedWindow
	err = json.Unmarshal([]byte(result), &jsonResult)
	require.NoError(t, err)
	assert.Len(t, jsonResult, 1)
	assert.Equal(t, "mark\"with\"quotes", jsonResult[0].Mark)
	assert.Equal(t, "app\nwith\nnewlines", jsonResult[0].AppName)
}

func TestListOutputFormatter_FormatJSON_ExactFormat(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := format.NewListOutputFormatter(&buf, "json")
	require.NoError(t, err)

	windows := []format.MarkedWindow{
		{
			Mark:        "test-mark",
			WindowID:    123,
			AppName:     "Test App",
			WindowTitle: "Test Title",
			Workspace:   "workspace1",
			AppBundleID: "com.test.app",
		},
	}

	err = formatter.Format(windows)
	require.NoError(t, err)

	result := strings.TrimSpace(buf.String())
	// Should be properly indented JSON
	assert.True(t, strings.HasPrefix(result, "["))
	assert.True(t, strings.HasSuffix(result, "]"))
	assert.Contains(t, result, "  ") // Should have indentation

	// Verify structure
	var jsonResult []format.MarkedWindow
	err = json.Unmarshal([]byte(result), &jsonResult)
	require.NoError(t, err)
	assert.Len(t, jsonResult, 1)
	assert.Equal(t, "test-mark", jsonResult[0].Mark)
	assert.Equal(t, 123, jsonResult[0].WindowID)
}

func TestListOutputFormatter_FormatCSV_SpecialCharacters(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := format.NewListOutputFormatter(&buf, "csv")
	require.NoError(t, err)

	windows := []format.MarkedWindow{
		{
			Mark:        "mark,with,commas",
			WindowID:    1,
			AppName:     "app\"with\"quotes",
			WindowTitle: "title\nwith\nnewlines",
			Workspace:   "workspace",
			AppBundleID: "bundle.id",
		},
	}

	err = formatter.Format(windows)
	require.NoError(t, err)

	result := buf.String()
	reader := csv.NewReader(strings.NewReader(result))
	records, err := reader.ReadAll()
	require.NoError(t, err)
	assert.Len(t, records, 2) // Header + 1 row
	assert.Equal(t, "mark,with,commas", records[1][0])
	assert.Equal(t, "app\"with\"quotes", records[1][2])
}

func TestListOutputFormatter_FormatCSV_ExactFormat(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := format.NewListOutputFormatter(&buf, "csv")
	require.NoError(t, err)

	windows := []format.MarkedWindow{
		{
			Mark:        "mark1",
			WindowID:    1,
			AppName:     "App1",
			WindowTitle: "Title1",
			Workspace:   "Workspace1",
			AppBundleID: "com.app1",
		},
		{
			Mark:        "mark2",
			WindowID:    2,
			AppName:     "",
			WindowTitle: "",
			Workspace:   "",
			AppBundleID: "",
		},
	}

	err = formatter.Format(windows)
	require.NoError(t, err)

	result := buf.String()
	reader := csv.NewReader(strings.NewReader(result))
	records, err := reader.ReadAll()
	require.NoError(t, err)
	assert.Len(t, records, 3) // Header + 2 rows

	// Verify header
	expectedHeader := []string{
		"mark",
		"window_id",
		"app_name",
		"window_title",
		"workspace",
		"app_bundle_id",
	}
	assert.Equal(t, expectedHeader, records[0])

	// Verify first row
	assert.Equal(t, []string{"mark1", "1", "App1", "Title1", "Workspace1", "com.app1"}, records[1])

	// Verify second row with empty values
	assert.Equal(t, []string{"mark2", "2", "", "", "", ""}, records[2])
}

func TestListOutputFormatter_FormatText_ExactFormat(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := format.NewListOutputFormatter(&buf, "text")
	require.NoError(t, err)

	windows := []format.MarkedWindow{
		{
			Mark:        "mark1",
			WindowID:    1,
			AppName:     "App1",
			WindowTitle: "Title1",
			Workspace:   "Workspace1",
			AppBundleID: "com.app1",
		},
		{
			Mark:        "mark2",
			WindowID:    22,
			AppName:     "Very Long App Name",
			WindowTitle: "Short",
			Workspace:   "",
			AppBundleID: "",
		},
	}

	err = formatter.Format(windows)
	require.NoError(t, err)

	result := strings.TrimSpace(buf.String())
	lines := strings.Split(result, "\n")
	assert.Len(t, lines, 2)

	// Verify structure: mark | window_id | app_name | window_title | workspace | app_bundle_id
	for _, line := range lines {
		fields := strings.Split(line, "|")
		assert.Len(t, fields, 6, "Each line should have 6 fields separated by |")
		// Trim spaces from fields
		for i, field := range fields {
			fields[i] = strings.TrimSpace(field)
		}
		// Verify we have the mark
		assert.NotEmpty(t, fields[0])
		// Verify we have window_id
		assert.NotEmpty(t, fields[1])
	}

	// Verify first line content
	assert.Contains(t, lines[0], "mark1")
	assert.Contains(t, lines[0], "1")
	assert.Contains(t, lines[0], "App1")

	// Verify second line has underscores for empty values
	assert.Contains(t, lines[1], "mark2")
	assert.Contains(t, lines[1], "_")
}

func TestListOutputFormatter_FormatText_SingleWindow(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := format.NewListOutputFormatter(&buf, "text")
	require.NoError(t, err)

	windows := []format.MarkedWindow{
		{
			Mark:        "single",
			WindowID:    42,
			AppName:     "SingleApp",
			WindowTitle: "SingleTitle",
			Workspace:   "SingleWorkspace",
			AppBundleID: "com.single",
		},
	}

	err = formatter.Format(windows)
	require.NoError(t, err)

	result := strings.TrimSpace(buf.String())
	// Single window should not have trailing newline in the trimmed result
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "single")
	assert.Contains(t, result, "42")
	// Should have pipe separators
	assert.Contains(t, result, "|")
}

func TestListOutputFormatter_FormatJSON_MultipleWindows(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := format.NewListOutputFormatter(&buf, "json")
	require.NoError(t, err)

	windows := []format.MarkedWindow{
		{
			Mark:        "mark1",
			WindowID:    1,
			AppName:     "App1",
			WindowTitle: "Title1",
			Workspace:   "WS1",
			AppBundleID: "bundle1",
		},
		{
			Mark:        "mark2",
			WindowID:    2,
			AppName:     "App2",
			WindowTitle: "Title2",
			Workspace:   "WS2",
			AppBundleID: "bundle2",
		},
		{
			Mark:        "mark3",
			WindowID:    3,
			AppName:     "App3",
			WindowTitle: "Title3",
			Workspace:   "WS3",
			AppBundleID: "bundle3",
		},
	}

	err = formatter.Format(windows)
	require.NoError(t, err)

	result := strings.TrimSpace(buf.String())
	var jsonResult []format.MarkedWindow
	err = json.Unmarshal([]byte(result), &jsonResult)
	require.NoError(t, err)
	assert.Len(t, jsonResult, 3)
	assert.Equal(t, "mark1", jsonResult[0].Mark)
	assert.Equal(t, "mark2", jsonResult[1].Mark)
	assert.Equal(t, "mark3", jsonResult[2].Mark)
}

func TestListOutputFormatter_FormatCSV_MultipleWindows(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := format.NewListOutputFormatter(&buf, "csv")
	require.NoError(t, err)

	windows := []format.MarkedWindow{
		{
			Mark:        "mark1",
			WindowID:    1,
			AppName:     "App1",
			WindowTitle: "Title1",
			Workspace:   "WS1",
			AppBundleID: "bundle1",
		},
		{
			Mark:        "mark2",
			WindowID:    2,
			AppName:     "App2",
			WindowTitle: "Title2",
			Workspace:   "WS2",
			AppBundleID: "bundle2",
		},
	}

	err = formatter.Format(windows)
	require.NoError(t, err)

	result := buf.String()
	reader := csv.NewReader(strings.NewReader(result))
	records, err := reader.ReadAll()
	require.NoError(t, err)
	assert.Len(t, records, 3) // Header + 2 rows
	assert.Equal(t, "mark1", records[1][0])
	assert.Equal(t, "mark2", records[2][0])
}

func TestListOutputFormatter_FormatText_AllEmptyFields(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := format.NewListOutputFormatter(&buf, "text")
	require.NoError(t, err)

	windows := []format.MarkedWindow{
		{
			Mark:        "mark1",
			WindowID:    1,
			AppName:     "",
			WindowTitle: "",
			Workspace:   "",
			AppBundleID: "",
		},
	}

	err = formatter.Format(windows)
	require.NoError(t, err)

	result := buf.String()
	// Should have mark and window_id, but underscores for empty fields
	assert.Contains(t, result, "mark1")
	assert.Contains(t, result, "1")
	// Should have 4 underscores (one for each empty field)
	underscoreCount := strings.Count(result, "_")
	assert.GreaterOrEqual(t, underscoreCount, 4)
}

func TestListOutputFormatter_FormatJSON_AllEmptyFields(t *testing.T) {
	var buf bytes.Buffer
	formatter, err := format.NewListOutputFormatter(&buf, "json")
	require.NoError(t, err)

	windows := []format.MarkedWindow{
		{
			Mark:        "mark1",
			WindowID:    1,
			AppName:     "",
			WindowTitle: "",
			Workspace:   "",
			AppBundleID: "",
		},
	}

	err = formatter.Format(windows)
	require.NoError(t, err)

	result := strings.TrimSpace(buf.String())
	var jsonResult []format.MarkedWindow
	err = json.Unmarshal([]byte(result), &jsonResult)
	require.NoError(t, err)
	assert.Len(t, jsonResult, 1)
	assert.Empty(t, jsonResult[0].AppName)
	assert.Empty(t, jsonResult[0].WindowTitle)
	assert.Empty(t, jsonResult[0].Workspace)
	assert.Empty(t, jsonResult[0].AppBundleID)
}

func TestListOutputFormatter_FormatEmpty(t *testing.T) {
	t.Run("JSON format outputs empty array", func(t *testing.T) {
		var buf bytes.Buffer
		formatter, err := format.NewListOutputFormatter(&buf, "json")
		require.NoError(t, err)

		err = formatter.FormatEmpty("")
		require.NoError(t, err)

		result := strings.TrimSpace(buf.String())
		assert.Equal(t, "[]", result)
	})

	t.Run("CSV format outputs header only", func(t *testing.T) {
		var buf bytes.Buffer
		formatter, err := format.NewListOutputFormatter(&buf, "csv")
		require.NoError(t, err)

		err = formatter.FormatEmpty("")
		require.NoError(t, err)

		result := buf.String()
		reader := csv.NewReader(strings.NewReader(result))
		records, err := reader.ReadAll()
		require.NoError(t, err)
		assert.Len(t, records, 1) // Header only
		expectedHeader := []string{
			"mark",
			"window_id",
			"app_name",
			"window_title",
			"workspace",
			"app_bundle_id",
		}
		assert.Equal(t, expectedHeader, records[0])
	})

	t.Run("Text format outputs message", func(t *testing.T) {
		var buf bytes.Buffer
		formatter, err := format.NewListOutputFormatter(&buf, "text")
		require.NoError(t, err)

		err = formatter.FormatEmpty("No marks found")
		require.NoError(t, err)

		result := strings.TrimSpace(buf.String())
		assert.Equal(t, "No marks found", result)
	})

	t.Run("Text format with empty message outputs nothing", func(t *testing.T) {
		var buf bytes.Buffer
		formatter, err := format.NewListOutputFormatter(&buf, "text")
		require.NoError(t, err)

		err = formatter.FormatEmpty("")
		require.NoError(t, err)

		result := strings.TrimSpace(buf.String())
		assert.Empty(t, result)
	})
}
