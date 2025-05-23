# Command: `mark`

## Description

This extension is heavily base in `sway mark` see in `man sway(1)` for more information.

```text
mark --add|--replace [--toggle] <identifier>

Marks are arbitrary labels that can be used to identify certain windows and then jump to them at a later time. Each identifier can only be set on a single window at a time since they act as a unique identifier. By default, mark sets identifier as the only mark on a window. --add will instead add identifier to the list of current marks for that window. If --toggle is specified mark will remove identifier if it is already marked.
```

- Each window may have one or more marks. (list of strings)
- When no flag is specified, it behaves like `--replace` and replaces all marks on the window with the new one.

## Usage
```bash
aerospace-marks mark [--add|--replace] [--toggle] <identifier>

aerospace-marks mark foo
# Will mark the current focused window with the identifier "foo"

aerospace-marks mark --add foo
# Will add the identifier "foo" to the current focused window

aerospace-marks mark bar
# Will replace the current mark with "bar" on the current focused window
# Same as `--replace` flag

aerospace-marks mark --toggle foo
# Will toggle the mark "foo" on the current focused window
# If the mark "foo" exists, it will be removed
# If the mark "foo" does not exist, it will be added
```

## Options

- `--add` - Add a mark to the window. 
- `--replace` - Replace the current mark with the new one.
- `--toggle` - Toggle the mark on the window. If the mark is already set, it will be removed.
- `--window-id` - The id of the window to mark. If not specified, it will use the current focused window.
- `--silent` - Suppress output messages. This is useful for scripting or when you don't want to pipe the output.

## Examples

```bash
aerospace-marks mark --add mymark  # Add a mark to the current focused window
```
