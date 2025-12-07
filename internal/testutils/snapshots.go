package testutils

import (
	"encoding/json"
	"strings"

	"github.com/kr/pretty"
	"gopkg.in/yaml.v3"
)

// SnapshotBuilder renders a structured snapshot that highlights
// the context, command, and outcome of a CLI test.
type SnapshotBuilder struct {
	title    string
	command  string
	contexts []snapshotSection
	stdout   string
	stderr   string
}

type snapshotSection struct {
	title   string
	content string
}

// NewSnapshotBuilder creates a new builder with the provided command text.
func NewSnapshotBuilder(command string) *SnapshotBuilder {
	return (&SnapshotBuilder{}).WithCommand(command)
}

// WithTitle sets an optional snapshot title.
func (s *SnapshotBuilder) WithTitle(title string) *SnapshotBuilder {
	s.title = strings.TrimSpace(title)
	return s
}

// WithCommand sets or overrides the command text.
func (s *SnapshotBuilder) WithCommand(command string) *SnapshotBuilder {
	s.command = strings.TrimSpace(command)
	return s
}

// WithContext appends a named context section (e.g. AeroSpace windows, stored marks).
func (s *SnapshotBuilder) WithContext(title string, value any) *SnapshotBuilder {
	s.contexts = append(s.contexts, snapshotSection{
		title:   strings.TrimSpace(title),
		content: formatSectionValue(value),
	})
	return s
}

// WithDetail appends a named detail section (e.g. computed result blocks).
func (s *SnapshotBuilder) WithDetail(title string, value any) *SnapshotBuilder {
	return s.WithContext(title, value)
}

// WithOutput sets the primary command output (stdout).
func (s *SnapshotBuilder) WithOutput(out string) *SnapshotBuilder {
	s.stdout = strings.TrimSpace(out)
	return s
}

// WithError sets the error text and marks the outcome as error.
func (s *SnapshotBuilder) WithError(err error) *SnapshotBuilder {
	if err != nil {
		s.stderr = strings.TrimSpace(err.Error())
	}
	return s
}

// WithErrorText sets the error text and marks the outcome as error.
func (s *SnapshotBuilder) WithErrorText(err string) *SnapshotBuilder {
	s.stderr = strings.TrimSpace(err)
	return s
}

// WithStdout sets the stdout text.
func (s *SnapshotBuilder) WithStdout(out string) *SnapshotBuilder {
	s.stdout = strings.TrimRight(out, "\n")
	return s
}

// WithStderr sets the stderr text.
func (s *SnapshotBuilder) WithStderr(err string) *SnapshotBuilder {
	s.stderr = strings.TrimRight(err, "\n")
	return s
}

// WithResult sets stdout and stderr together.
func (s *SnapshotBuilder) WithResult(stdout, stderr string) *SnapshotBuilder {
	s.stdout = strings.TrimRight(stdout, "\n")
	s.stderr = strings.TrimRight(stderr, "\n")
	return s
}

// String builds the final snapshot string.
func (s *SnapshotBuilder) String() string {
	var builder strings.Builder

	builder.WriteString("Context:\n")
	if len(s.contexts) == 0 {
		builder.WriteString("  (none)\n\n")
	} else {
		builder.WriteString(renderContextSections(s.contexts))
		builder.WriteString("\n")
	}

	builder.WriteString("Command:\n")
	builder.WriteString("  $ ")
	builder.WriteString(s.command)
	builder.WriteString("\n\n")

	builder.WriteString("Result:\n")
	builder.WriteString(renderResultBlock("stdout", s.stdout))
	builder.WriteString(renderResultBlock("stderr", s.stderr))

	return strings.TrimRight(builder.String(), "\n")
}

// CommandString renders a consistent CLI command prefix for snapshots.
func CommandString(args ...string) string {
	return strings.TrimSpace("aerospace-marks " + strings.Join(args, " "))
}

// SnapshotContext represents a named context section for the snapshot.
type SnapshotContext struct {
	Title string
	Value any
}

// SnapshotSpec captures all fields required to render a snapshot in a single struct.
type SnapshotSpec struct {
	Command  string
	Stdout   string
	Stderr   string
	Contexts []SnapshotContext
}

// RenderSnapshotSpec builds a snapshot from a SnapshotSpec struct.
func RenderSnapshotSpec(spec SnapshotSpec) string {
	builder := NewSnapshotBuilder(spec.Command).WithResult(spec.Stdout, spec.Stderr)
	for _, ctx := range spec.Contexts {
		builder.WithContext(ctx.Title, ctx.Value)
	}
	return builder.String()
}

// Context creates a SnapshotContext value.
func Context(title string, value any) SnapshotContext {
	return SnapshotContext{Title: strings.TrimSpace(title), Value: value}
}

func renderContextSections(sections []snapshotSection) string {
	var builder strings.Builder
	for _, ctx := range sections {
		builder.WriteString("  ")
		builder.WriteString(ctx.title)
		builder.WriteString(":\n")
		builder.WriteString(indent(ctx.content, "    "))
		builder.WriteString("\n")
	}
	return builder.String()
}

func indent(value, prefix string) string {
	if value == "" {
		return prefix
	}

	lines := strings.Split(strings.TrimRight(value, "\n"), "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}

func formatSectionValue(value any) string {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	default:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return strings.TrimSpace(pretty.Sprint(value))
		}

		yamlBytes, err := formatJSONAsYAML(jsonBytes)
		if err != nil {
			return strings.TrimSpace(string(jsonBytes))
		}

		return strings.TrimSpace(string(yamlBytes))
	}
}

func formatJSONAsYAML(jsonBytes []byte) ([]byte, error) {
	var data any
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		return nil, err
	}
	return yaml.Marshal(data)
}

func renderResultBlock(title, value string) string {
	var builder strings.Builder
	if value == "" {
		builder.WriteString("  ")
		builder.WriteString(title)
		builder.WriteString(": \"\"")
		builder.WriteString("\n")
		return builder.String()
	}

	builder.WriteString("  ")
	builder.WriteString(title)
	builder.WriteString(":\n")
	builder.WriteString(indent(value, "    "))
	builder.WriteString("\n")
	return builder.String()
}
