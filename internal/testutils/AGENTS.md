# internal/testutils/

Test harness for CLI command tests. Test-only.

## Responsibility

- `clis.go` — execute cobra commands under test:
  - `CmdExecute(cmd, args...)` and `CmdExecuteWithStdin(...)` — set args and
    run, returning captured stdout (or an error if stderr is non-empty).
  - `CaptureStdOut(f)` — redirect os.Stdout/os.Stderr into pipes for capture.
  - `MockEmptyAerspaceMarkWindows` — minimal no-op aerospace double.
- `snapshots.go` — build deterministic, readable golden snapshots:
  - `SnapshotSpec` + `RenderSnapshotSpec` — the preferred entry point: pass
    `Command`, `Stdout`, `Stderr`, and named `Contexts` (rendered as YAML).
  - `NewSnapshotBuilder` — fluent builder for the same output.
  - `CommandString(args...)` — consistent `aerospace-marks ...` prefix.

Snapshot layout:
```
Context:
  <name>:
    <yaml>
Command:
  $ aerospace-marks ...
Result:
  stdout: ...
  stderr: ...
```

## Conventions

- Snapshot context values are rendered as YAML (JSON-marshalable); strings are
  printed verbatim.
- Use these helpers in `cmd/*_test.go` so golden files stay uniform.

## What belongs here

- Cross-command test plumbing and snapshot rendering helpers.

## What does NOT belong here

- Per-command test cases (live in `cmd/*_test.go`).
- Mocks of domain interfaces (live in `internal/mocks`).
