# aerospace-ext mark

USAGE: `aerospacex mark [--add|--replace] [--toggle] <identifier>`

## Description

This extension is heavily base in `sway mark` see in `man sway(1)` for more information.

```text
mark --add|--replace [--toggle] <identifier>

Marks are arbitrary labels that can be used to identify certain windows and then jump to them at a later time. Each identifier can only be set on a single window at a time since they act as a unique identifier. By default, mark sets identifier as the only mark on a window. --add will instead add identifier to the list of current marks for that window. If --toggle is specified mark will remove identifier if it is already marked.
```

## Options

- `--add` - Add a mark to the window.

## Examples

```bash
aerospacex mark --add mymark  # Add a mark to the current focused window
```

## Implemantation details

 - The command will call `aerospace list-windows --focused` and collect the window id (first column).
 - It will store the given mark among with the window id in memory at first.
