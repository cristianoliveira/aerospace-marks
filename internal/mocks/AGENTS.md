# internal/mocks/

Test doubles for storage and aerospace, plus JSON fixture loaders. Test-only.

## Responsibility

- `mocks.go` — entry helpers used by command tests:
  - `MockStorageDBClient(ctrl)` → (`MockStorageDBClient`, `MockMarkStorage`).
  - `MockAerospaceConnection(ctrl)` → (`MockAeroSpaceConnection`, client).
    Swaps the aerospace-ipc default connector so `NewAeroSpaceClient()` is
    backed by a mock socket.
  - `LoadMarksFixture` / `LoadAeroWindowsFixture` / `LoadAeroWindowsFixtureRaw`
    read JSON fixtures into typed slices.
- `storage/` — gomock mocks for `StorageDBClient` and `MarkStorage`.
- `aerospacecli/` — gomock mocks for the ipc connection and connector.
- `fixtures/` — JSON data (`aerospace/list-windows-all.json`,
  `storage/list-marked-windows.json`).

## Generated code — DO NOT EDIT

`*_mock.go` files are generated from interfaces via `scripts/mock-generator.sh`.
Regenerate whenever `storage.MarkStorage`, `StorageDBClient`, or the aerospace
ipc interfaces change.

## What belongs here

- New mock helpers or fixture files for repeatable test scenarios.

## What does NOT belong here

- Production logic. If a mock grows real behavior, the abstraction is wrong.
- Editing generated files by hand.
