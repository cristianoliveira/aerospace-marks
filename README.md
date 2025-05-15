# aerospace-marks

**Beta**: I use this daily, but it's still a work in progress. Please report any issues or ideas in the [issues](https://github.com/cristianoliveira/aerospace-marks/issues) section.

This is a CLI tool for [AeroSpace WM](https://github.com/nikitabobko/AeroSpace) that helps manage marks on windows.

It’s heavily inspired by [sway marks](https://man.archlinux.org/man/sway.5.en), but follows the `aerospace` style of commands:
```text
mark --add|--replace [--toggle] <identifier>

Marks are arbitrary labels that can be used to identify certain
windows and then jump to them at a later time. Each identifier
can only be set on a single window at a time since they act as
a unique identifier. By default, mark sets identifier as the only
mark on a window. --add will instead add identifier to the list
of current marks for that window. If --toggle is specified mark will
remove identifier if it is already marked.
```
### Why not use named workspaces?

> TL;DR: Dynamic mappings to specific windows, not workspaces.

I wanted something more granular than workspaces — I want to jump to a specific window, given a context.

Sometimes my "video" context means YouTube, sometimes it’s a video player, sometimes a browser. I want to be able to jump to the right window, regardless of which workspace it's on.

## Usage

Allow one to mark a window with the given identifier. 
```text
aerospace-marks mark [--add|--replace] [--toggle] <identifier>
```
And then focus to a window with the given mark.
```text
aerospace-marks focus <identifier>
```

### Config Usage

```toml
# ~/.config/aerospace/config.toml
[mode.main.binding] 
# Mark 
cmd-shift-h = ["exec-and-forget aerospace-marks mark h 2> /tmp/marks.log", "mode main"]
cmd-shift-j = ["exec-and-forget aerospace-marks mark j 2> /tmp/marks.log", "mode main"]
cmd-shift-k = ["exec-and-forget aerospace-marks mark k 2> /tmp/marks.log", "mode main"]
cmd-shift-l = ["exec-and-forget aerospace-marks mark l 2> /tmp/marks.log", "mode main"]
cmd-shift-b = ["exec-and-forget aerospace-marks mark browser 2> /tmp/marks.log", "mode main"]

# Focus
cmd-ctrl-h = ["exec-and-forget aerospace-marks focus h 2> /tmp/marks.log", "mode main"]
cmd-ctrl-j = ["exec-and-forget aerospace-marks focus j 2> /tmp/marks.log", "mode main"]
cmd-ctrl-k = ["exec-and-forget aerospace-marks focus k 2> /tmp/marks.log", "mode main"]
cmd-ctrl-l = ["exec-and-forget aerospace-marks focus l 2> /tmp/marks.log", "mode main"]
cmd-ctrl-b = ["exec-and-forget aerospace-marks focus browser 2> /tmp/marks.log", "mode main"]
```

See more in [documentation](docs/aerospace-marks)

## Installation

### Using Nix

If you have Nix installed, you can build and install `aerospace-marks` using the following command:

```bash
nix profile install github:cristianoliveira/aerospace-marks
```

You can also run without installing it by using:

```bash
nix run github:cristianoliveira/aerospace-marks
```

This will build the default package defined in `flake.nix`.

### Installing with Go

If you have Go installed, you can install `aerospace-marks` directly using:

```bash
go install github.com/cristianoliveira/aerospace-marks@latest
```

This will download and install the latest version of `aerospace-marks` to your `$GOPATH/bin`.

### Building from Source

To build from source, ensure you have Go installed. Then, clone the repository and run:

```bash
git clone https://github.com/cristianoliveira/aerospace-marks.git
cd aerospace-marks
go build -o aerospace-marks
```

This will create an executable named `aerospace-marks` in the current directory.

### Packages

- AeroSpace Socket Client - [aerospacecli](pkgs/aerospacecli)
