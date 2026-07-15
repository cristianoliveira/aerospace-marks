# cmd/

CLI command layer. Cobra command definitions and the application's
composition point for wiring commands to their dependencies.

## Responsibility

Each command is a factory function that accepts injected dependencies and
returns a `*cobra.Command`:

- `mark` / `unmark` — manage marks on windows (storage only).
- `focus` / `summon` — act on a window identified by mark (storage + aerospace).
- `list` / `get` — read marks/windows (storage + aerospace).
- `info` — diagnostics (storage + aerospace).

Commands NEVER construct their own `MarkStorage` or aerospace client.
They receive `storage.MarkStorage` and `aerospace.AerosSpaceMarkWindows` as
arguments. `NewRootCmd(storage, aerospaceClient)` assembles the tree; the
`--output/-o` flag (text|json|csv) is attached via `enableOutputFlag` to the
read/action commands (focus, list, summon, get).

`VERSION` in `root.go` is generated at build time — do not edit by hand.

## What belongs here

- Cobra command definitions, flags, arg validation (use `internal/cli`).
- Translating user input into calls on injected interfaces.
- Formatting output via `internal/format` formatters.

## What does NOT belong here

- Business logic or persistence — delegate to `internal/storage`.
- Direct IPC with AeroSpace — go through `internal/aerospace`.
- New formatting logic — extend `internal/format`.

## Testing

Golden-snapshot tests drive every command. See the pattern in `*_test.go`:

1. Build gomock controllers; mock storage via `mocks.MockStorageDBClient`
   and the aerospace socket via `mocks.MockAerospaceConnection`.
2. Assemble with `cmd.NewRootCmd(strg, aerospaceClient)`.
3. Execute via `testutils.CmdExecute(cmd, args...)`.
4. Render with `testutils.RenderSnapshotSpec(...)` and assert with
   `snaps.MatchSnapshot(t, snapshot)`.

Golden files live in `__snapshots__/` (`*_test.snap`). Cover both happy and
unhappy paths (missing args, empty identifier, not-found windows).

- Run: `make test`
- Regenerate snapshots after intentional output changes: `make update-snap-all`
