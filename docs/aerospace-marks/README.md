# aerospace-marks

This is an cli for AeroSpace WM to manage marks on windows. 

Similar to sway marks [sway](https://man.archlinux.org/man/sway.5.en)
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

## Command: `mark`

Mark a window with the given identifier.

USAGE: `aerospace-marks mark [--add|--replace] [--toggle] <identifier>`

[read more](/docs/aerospace-marks/CMD_MARK.md)

## Command: `focus`

Focus to a window with the given mark.

USAGE: `aerospace-marks focus <identifier>`

## Command: `list`

List all marks.

USAGE: `aerospace-marks list`

## Command: `config`

Show the current configurations and other info related

# Implemantation details

 - The command will send to aerospace socket `list-windows --focused` and collect the window id (first column).
 - It will store the given mark among with the window id in memory at first.

### Storage

 - The marks are stored using sqlite3 in the `~/.local/state/aerospace-marks/storage.db` file.
 - Each window may have one or more marks. (list of strings)
 - The table is called `marks` and has the following columns:
    - `window_id` - The id of the window.
    - `mark` - The mark of the window.
   
 - The sqlite3 database is created if it does not exist.
