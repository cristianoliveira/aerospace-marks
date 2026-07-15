# internal/

Reusable packages with narrow, stable APIs. Keep packages cohesive and
covered by tests.

## Dependency direction (one-way, no cycles)

```
cmd ──► aerospace, format, storage, logger, stdout, cli, constants
aerospace ──► logger
format     ──► logger
storage    ──► logger, constants, storage/db/queries
storage/db/queries  (sqlc-generated leaf)
logger / stdout / cli / constants  (leaf utilities)
mocks, testutils    (test-only)
```

`main.go` is the only place that constructs concrete clients (composition root).
Everything else programs against interfaces.

## Package index

- `aerospace/` — adapter over the external `aerospace-ipc` package. See its AGENTS.md.
- `format/` — pure text/json/csv output formatters. See its AGENTS.md.
- `storage/` — SQLite persistence + migrations + sqlc queries. See its AGENTS.md.
- `logger/` — leveled logger configured via env (`internal/constants`).
- `stdout/` — `ErrorAndExit` prints to stderr and exits. `ShouldExit` flag is
  for tests; set it `false` to assert error output without killing the process.
- `cli/` — cobra argument validators (e.g. `ValidateArgIsNotEmpty`).
- `constants/` — canonical names for every environment variable. Add new env
  vars here, never inline string literals.
- `mocks/` — gomock-generated clients + fixture loaders. Test-only. See its AGENTS.md.
- `testutils/` — CLI test harness + snapshot builder. See its AGENTS.md.

## Conventions

- Inject dependencies via interfaces defined where they are consumed (e.g.
  `storage.MarkStorage`, `aerospace.AerosSpaceMarkWindows`).
- Return `error`; do not panic in packages (only the aerospace client asserts
  on misuse of an uninitialised client).
- Prefer early returns and flat logic.
