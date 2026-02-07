# aerospace-marks

[Vim-like marks](https://i3wm.org/docs/userguide.html#vim_like_marks) for [AeroSpace WM](https://github.com/nikitabobko/AeroSpace)

## AeroSpace Compatibility

 - Stable version: v1.0.0
 - v0.20.0 use v0.3.x or higher
 - v0.19.1 use v0.2.x or lower

## Demo

https://github.com/user-attachments/assets/cfd84749-c436-465d-8f66-486eb2303e30

## Description

Allows you to add custom marks to windows and use them to set focus or summon to the current workspace.

It’s heavily inspired by [sway marks](https://man.archlinux.org/man/sway.5.en), but follows the `aerospace` style of commands

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

Also, by pairing with [aerospace-scratchpad](https://github.com/cristianoliveira/aerospace-scratchpad) allows me to send away windows, mark then with a key for later bring then back!

## Basic Usage

Allow one to mark the current focused window with the given identifier.
```text
aerospace-marks mark [--add|--replace] [--toggle] <identifier>
```
And then focus to a window with the given mark.
```text
aerospace-marks focus <identifier>
```
Or summon a marked window to the current workspace.
```text
aerospace-marks summon <identifier>
```

## Advanced Usage

### Output Formats

Multiple commands support structured output formats for scripting via the `--output` (or `-o`) flag:

#### List Command
```bash
# JSON format for jq processing
aerospace-marks list -o json | jq '.[] | select(.app_name == "Brave")'

# CSV format for spreadsheet tools
aerospace-marks list -o csv | csvcut -c mark,window_title

# Text format (default) for awk/sed
aerospace-marks list | awk -F'|' '{print $1}'
```

#### Focus Command
```bash
# JSON format for scripting
aerospace-marks focus mark1 -o json | jq '.window_id'

# CSV format
aerospace-marks focus mark1 -o csv
```

#### Summon Command
```bash
# JSON format with focus flag
aerospace-marks summon mark1 --focus -o json | jq '.action'

# Text format (default)
aerospace-marks summon mark1
```

#### Get Command
```bash
# JSON format for full window info
aerospace-marks get mark1 -o json | jq '.app_name'

# JSON format for single field (window ID)
aerospace-marks get mark1 -i -o json | jq '.result'

# Plain text for single field (backward compatible)
aerospace-marks get mark1 -i
```

See more in [documentation](docs/)

Check [vim-like marks](https://i3wm.org/docs/userguide.html#vim_like_marks) for a more advanced usage.

### Config Usage
```toml
# ~/.config/aerospace/config.toml
[mode.main.binding]

# Vim's like marks, similar to i3-input
# cmd + ctrl + m and <letter> -- mark a window with a given <letter>
cmd-ctrl-m = """
exec-and-forget aerospace-marks mark \
    $(osascript -e 'text returned of (display dialog "mark" default answer "")')
"""
# cmd + ctrl + ' and <letter> -- set focus on the window marked with <letter>
cmd-ctrl-quote = """
exec-and-forget aerospace-marks focus \
    $(osascript -e 'text returned of (display dialog "focus" default answer "")')
"""

# Mark current window with a given identifier
cmd-shift-h = ["exec-and-forget aerospace-marks mark h"]
cmd-shift-j = ["exec-and-forget aerospace-marks mark j"]
cmd-shift-k = ["exec-and-forget aerospace-marks mark k"]
cmd-shift-l = ["exec-and-forget aerospace-marks mark l"]

cmd-shift-b = ["exec-and-forget aerospace-marks mark browser", "mode main"]

# Focus to a window with the given identifier
cmd-ctrl-h = ["exec-and-forget aerospace-marks focus h"]
cmd-ctrl-j = ["exec-and-forget aerospace-marks focus j"]
cmd-ctrl-k = ["exec-and-forget aerospace-marks focus k"]
cmd-ctrl-l = ["exec-and-forget aerospace-marks focus l"]

cmd-ctrl-b = ["exec-and-forget aerospace-marks focus browser", "mode main"]
```

## Installation

### Using Homebrew

If you have Homebrew installed, you can install `aerospace-marks` using the following command:

```bash
brew install cristianoliveira/tap/aerospace-marks
```

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

#### Post installation

After installing, you may need to include aerospace-marks in aerospace context.

Check where the binary is installed, run:
```bash
echo $(which aerospace-marks) | sed 's/\/aerospace-marks//g'
```

And in your config add:
```toml
[exec]
    inherit-env-vars = true

# OR

[exec.env-vars]
    # Replace 'aerospace-marks/install/path' with the actual path from the above command
    PATH = 'aerospace-marks/install/path/bin:${PATH}'
```

### Building from Source

To build from source, ensure you have Go installed. Then, clone the repository and run:

```bash
git clone https://github.com/cristianoliveira/aerospace-marks.git
cd aerospace-marks
go build -o aerospace-marks
```

This will create an executable named `aerospace-marks` in the current directory.

## Troubleshooting

If you encounter issues while using `aerospace-marks`, you can use the following environment variables to help diagnose the problem:

- `AEROSPACE_MARKS_DB_PATH`: This variable sets the path for the AeroSpace marks database. By default, it is set to `$HOME/.local/state/aerospace-marks`. (You can connect with sqlite client)

- `AEROSPACE_MARKS_LOGS_PATH`: This variable specifies the path for the AeroSpace marks logs. The default path is `/tmp/aerospace-marks.log`. (use `tail -f <path>`)

- `AEROSPACE_MARKS_LOGS_LEVEL`: This variable determines the logging level for AeroSpace marks. The default level is `DISABLED`.

These environment variables can be set directly in the AeroSpace configuration file to ensure they are available whenever AeroSpace is running. Add the following to your [AeroSpace config](https://nikitabobko.github.io/AeroSpace/guide#config-location)

```toml
[exec.env-vars]
# Path for the AeroSpace marks database
AEROSPACE_MARKS_DB_PATH = "$HOME/.local/state/aerospace-marks"

# Path for the AeroSpace marks logs
AEROSPACE_MARKS_LOGS_PATH = "/tmp/aerospace-marks.log"

# Logging level for AeroSpace marks
AEROSPACE_MARKS_LOGS_LEVEL = "DEBUG"
```

Replace the paths and values with your desired settings.

### Packages

- AeroSpace Socket IPC - [aerospace-ipc](https://github.com/cristianoliveira/aerospace-ipc)

## License

This project is licensed under the terms of the LICENSE file.
