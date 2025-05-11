# aerospace-ext mark

USAGE: `aerospacex mark [--add|--replace] [--toggle] <identifier>`

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
- `--get-id <mark>` - Retrieve the window ID associated with the specified mark. This option allows you to find which window currently has a particular mark assigned to it. If the mark is not found, no window ID will be returned.

## Examples

```bash
aerospacex mark --add mymark  # Add a mark to the current focused window
```

## Implemantation details

 - The command will call `aerospace list-windows --focused` and collect the window id (first column).
 - It will store the given mark among with the window id in memory at first.

### Storage

 - The marks are stored using sqlite3 in the `~/.local/state/aerospacex/storage.db` file.
 - Each window may have one or more marks. (list of strings)
 - The table is called `marks` and has the following columns:
    - `window_id` - The id of the window.
    - `mark` - The mark of the window.
   
 - The sqlite3 database is created if it does not exist.
