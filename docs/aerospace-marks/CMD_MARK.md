# Command: `mark`

USAGE: `aerospace-marks mark [--add|--replace] [--toggle] <identifier>`

## Description

This extension is heavily base in `sway mark` see in `man sway(1)` for more information.

```text
mark --add|--replace [--toggle] <identifier>

Marks are arbitrary labels that can be used to identify certain windows and then jump to them at a later time. Each identifier can only be set on a single window at a time since they act as a unique identifier. By default, mark sets identifier as the only mark on a window. --add will instead add identifier to the list of current marks for that window. If --toggle is specified mark will remove identifier if it is already marked.
```

- Each window may have one or more marks. (list of strings)
- When no flag is specified, it behaves like `--replace` and replaces all marks on the window with the new one.

## Options

- `--add` - Add a mark to the window. 
- `--replace` - Replace the current mark with the new one.
- `--toggle` - Toggle the mark on the window. If the mark is already set, it will be removed.
- `--window-id` - The id of the window to mark. If not specified, it will use the current focused window.

## Examples

```bash
aerospace-marks mark --add mymark  # Add a mark to the current focused window
```
