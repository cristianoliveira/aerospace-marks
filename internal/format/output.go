package format

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// OutputFormat represents the output format type.
type OutputFormat string

const (
	// OutputFormatText is the default pipe-separated text format.
	OutputFormatText OutputFormat = "text"
	// OutputFormatJSON outputs data as JSON array.
	OutputFormatJSON OutputFormat = "json"
	// OutputFormatCSV outputs data as comma-separated values.
	OutputFormatCSV OutputFormat = "csv"
)

const (
	// textFormatColumnCount is the number of columns in the text output format.
	textFormatColumnCount = 6
)

// MarkedWindow represents a window with its mark.
type MarkedWindow struct {
	Mark        string `json:"mark"`
	WindowID    int    `json:"window_id"`
	AppName     string `json:"app_name"`
	WindowTitle string `json:"window_title"`
	Workspace   string `json:"workspace"`
	AppBundleID string `json:"app_bundle_id"`
}

// ListOutputFormatter formats a list of marked windows.
type ListOutputFormatter struct {
	format OutputFormat
	writer io.Writer
}

// NewListOutputFormatter creates a new ListOutputFormatter.
func NewListOutputFormatter(w io.Writer, format string) (*ListOutputFormatter, error) {
	normalized := strings.ToLower(strings.TrimSpace(format))
	switch normalized {
	case string(OutputFormatText):
		return &ListOutputFormatter{format: OutputFormatText, writer: w}, nil
	case string(OutputFormatJSON):
		return &ListOutputFormatter{format: OutputFormatJSON, writer: w}, nil
	case string(OutputFormatCSV):
		return &ListOutputFormatter{format: OutputFormatCSV, writer: w}, nil
	default:
		return nil, fmt.Errorf(
			"unsupported output format: %s (valid formats: text, json, csv)",
			format,
		)
	}
}

// Format formats and writes the list of marked windows.
func (f *ListOutputFormatter) Format(windows []MarkedWindow) error {
	switch f.format {
	case OutputFormatJSON:
		return f.formatJSON(windows)
	case OutputFormatCSV:
		return f.formatCSV(windows)
	case OutputFormatText:
		return f.formatText(windows)
	default:
		return fmt.Errorf("unsupported output format: %s", f.format)
	}
}

// FormatEmpty formats and writes empty results with an optional message for text format.
// For JSON, outputs "[]". For CSV, outputs header only. For text, outputs the message.
func (f *ListOutputFormatter) FormatEmpty(message string) error {
	switch f.format {
	case OutputFormatJSON:
		_, err := fmt.Fprintln(f.writer, "[]")
		return err
	case OutputFormatCSV:
		return f.formatCSV([]MarkedWindow{})
	case OutputFormatText:
		if message != "" {
			_, err := fmt.Fprintln(f.writer, message)
			return err
		}
		return nil
	default:
		return fmt.Errorf("unsupported output format: %s", f.format)
	}
}

// formatText formats windows as pipe-separated aligned columns.
func (f *ListOutputFormatter) formatText(windows []MarkedWindow) error {
	if len(windows) == 0 {
		return nil
	}

	// Convert to string rows for alignment
	rows := make([][]string, len(windows))
	for i, w := range windows {
		rows[i] = []string{
			w.Mark,
			strconv.Itoa(w.WindowID),
			f.emptyToUnderscore(w.AppName),
			f.emptyToUnderscore(w.WindowTitle),
			f.emptyToUnderscore(w.Workspace),
			f.emptyToUnderscore(w.AppBundleID),
		}
	}

	// Calculate column widths
	colWidths := make([]int, textFormatColumnCount)
	for _, row := range rows {
		for j, field := range row {
			if len(field) > colWidths[j] {
				colWidths[j] = len(field)
			}
		}
	}

	// Format and write
	var b strings.Builder
	for i, row := range rows {
		for j, field := range row {
			b.WriteString(fmt.Sprintf("%-*s", colWidths[j], field))
			if j < len(row)-1 {
				b.WriteString(" | ")
			}
		}
		if i < len(rows)-1 {
			b.WriteByte('\n')
		}
	}

	_, err := fmt.Fprintln(f.writer, b.String())
	return err
}

// formatJSON formats windows as JSON array.
func (f *ListOutputFormatter) formatJSON(windows []MarkedWindow) error {
	data, err := json.MarshalIndent(windows, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	_, err = fmt.Fprintln(f.writer, string(data))
	return err
}

// formatCSV formats windows as CSV with headers.
func (f *ListOutputFormatter) formatCSV(windows []MarkedWindow) error {
	writer := csv.NewWriter(f.writer)
	defer writer.Flush()

	headers := []string{
		"mark",
		"window_id",
		"app_name",
		"window_title",
		"workspace",
		"app_bundle_id",
	}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	for _, w := range windows {
		row := []string{
			w.Mark,
			strconv.Itoa(w.WindowID),
			w.AppName,
			w.WindowTitle,
			w.Workspace,
			w.AppBundleID,
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return writer.Error()
}

// emptyToUnderscore converts empty strings to "_" for text format.
func (f *ListOutputFormatter) emptyToUnderscore(s string) string {
	if s == "" {
		return "_"
	}
	return s
}

// OutputEvent describes a single command result in a structured way.
type OutputEvent struct {
	Command         string `json:"command"`
	Action          string `json:"action"`
	WindowID        int    `json:"window_id"`
	AppName         string `json:"app_name"`
	Workspace       string `json:"workspace"`
	TargetWorkspace string `json:"target_workspace"`
	Result          string `json:"result"`
	Message         string `json:"message"`
}

// OutputEventFormatter formats a single command result event.
type OutputEventFormatter struct {
	format OutputFormat
	writer io.Writer
}

// NewOutputEventFormatter creates a new OutputEventFormatter.
func NewOutputEventFormatter(w io.Writer, format string) (*OutputEventFormatter, error) {
	normalized := strings.ToLower(strings.TrimSpace(format))
	switch normalized {
	case string(OutputFormatText):
		return &OutputEventFormatter{format: OutputFormatText, writer: w}, nil
	case string(OutputFormatJSON):
		return &OutputEventFormatter{format: OutputFormatJSON, writer: w}, nil
	case string(OutputFormatCSV):
		return &OutputEventFormatter{format: OutputFormatCSV, writer: w}, nil
	default:
		return nil, fmt.Errorf(
			"unsupported output format: %s (valid formats: text, json, csv)",
			format,
		)
	}
}

// Format formats and writes a single output event.
func (f *OutputEventFormatter) Format(event OutputEvent) error {
	switch f.format {
	case OutputFormatJSON:
		return f.formatJSON(event)
	case OutputFormatCSV:
		return f.formatCSV(event)
	case OutputFormatText:
		return f.formatText(event)
	default:
		return fmt.Errorf("unsupported output format: %s", f.format)
	}
}

// formatText formats event as pipe-separated values.
// Format: <window_id> | <app_name> | <message> (for backward compatibility with get command).
func (f *OutputEventFormatter) formatText(event OutputEvent) error {
	// For single field outputs (when only Result is set and it's a simple value), output just the value
	if event.Command == "get" && event.Action == "" && event.AppName == "" && event.Message == "" {
		// Single field output (window-id, window-title, app-name, app-bundle-id)
		_, err := fmt.Fprintln(f.writer, event.Result)
		return err
	}

	appName := event.AppName
	if appName == "" {
		appName = "_"
	}
	message := event.Message
	if message == "" {
		message = event.Result
	}
	if message == "" {
		message = "_"
	}

	// For get command full output, maintain backward compatibility: <window_id> | <app_name> | <window_title>
	if event.Command == "get" && event.Action == "" && event.AppName != "" {
		output := fmt.Sprintf("%d | %s | %s",
			event.WindowID,
			appName,
			message,
		)
		_, err := fmt.Fprintln(f.writer, output)
		return err
	}

	// Default format: just the message/result
	_, err := fmt.Fprintln(f.writer, message)
	return err
}

// formatJSON formats event as JSON object.
func (f *OutputEventFormatter) formatJSON(event OutputEvent) error {
	data, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	_, err = fmt.Fprintln(f.writer, string(data))
	return err
}

// formatCSV formats event as CSV with header.
func (f *OutputEventFormatter) formatCSV(event OutputEvent) error {
	writer := csv.NewWriter(f.writer)
	defer writer.Flush()

	headers := []string{
		"command",
		"action",
		"window_id",
		"app_name",
		"workspace",
		"target_workspace",
		"result",
		"message",
	}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	row := []string{
		event.Command,
		event.Action,
		strconv.Itoa(event.WindowID),
		event.AppName,
		event.Workspace,
		event.TargetWorkspace,
		event.Result,
		event.Message,
	}
	if err := writer.Write(row); err != nil {
		return fmt.Errorf("failed to write CSV row: %w", err)
	}

	return writer.Error()
}
