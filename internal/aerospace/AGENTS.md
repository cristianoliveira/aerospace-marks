# internal/aerospace/

Adapter over the external AeroSpace window manager. Isolates the rest of the
codebase from `aerospace-ipc` details.

## Responsibility

- `AerosSpaceMarkWindows` interface (`windows.go`) — the contract commands use:
  `GetWindowByID(windowID int)` and `Client()` (raw `*aerospacecli.AeroSpaceWM`
  for when the ipc package does not expose what we need).
- `DefaultAeroSpaceWindows` — production implementation.

The external `aerospace-ipc` package exposes some helpers but not all. When you
need a command it does not wrap, call `Client()` and use
`ipc.Connection().SendCommand(command, args)` to send raw commands to AeroSpace.

## Raw command reference

- AeroSpaceWM source/docs: `~/other/AeroSpace/**` and `.tmp/docs/aerospace-docs`.
- aerospace-ipc docs/code: `.tmp/docs/aerospace-ipc`.
- `aerospace --help` lists available WM commands.

`GetWindowByID` fetches all windows and matches by id.

## What belongs here

- Window/WM lookups and any new thin wrappers over `aerospace-ipc`.

## What does NOT belong here

- Marks persistence (`internal/storage`).
- Output formatting (`internal/format`).

## Testing

Mocked in command tests via `mocks.MockAerospaceConnection`, which swaps the
default connector so `NewAeroSpaceClient()` returns a client backed by a mock
socket. Expect `SendCommand` calls with the exact command/args your wrapper
emits.
