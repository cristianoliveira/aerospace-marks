# internal/format/

Output formatting. Pure translation from in-memory structures to text, JSON,
or CSV. No storage or aerospace dependencies.

## Responsibility

- `OutputFormat` (text|json|csv) — the values accepted by the `--output` flag.
- `MarkedWindow` — the row shape for list output.
- `ListOutputFormatter` — formats `[]MarkedWindow` (aligned pipe-separated text,
  JSON array, CSV with header). `FormatEmpty` handles empty results per format.
- `OutputEventFormatter` + `OutputEvent` — formats a single command result
  (used by focus/summon/get). `get` has special single-field text behavior.

Formatters take an `io.Writer`, so they are easy to test against a buffer.

## What belongs here

- Any new output representation or column/formatting rule.

## What does NOT belong here

- Fetching windows/marks — receive already-resolved data.
- cobra command wiring — that is `cmd`.

## Conventions

- Empty fields render as `_` in text format; JSON/CSV use real values.
- Text list output is pipe-separated (` | `) with per-column width alignment
  (`textFormatColumnCount = 6`).
- Add a new format by extending the `switch` in the constructor and each `Format`.

## Testing

Direct unit tests (`output_test.go`, `output_event_test.go`, `list_test.go`)
drive each format with fixed inputs and assert writer output. Keep tests
deterministic — no time, randomness, or external state.
