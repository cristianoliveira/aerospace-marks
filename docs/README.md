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

Mark the current focused window with the given identifier. 
You may specify the window with `--window-id <id>` option.

USAGE: `aerospace-marks mark [--add|--replace] [--toggle] <identifier>`

[read more](/docs/CMD_MARK.md)

## Command: `focus`

Focus to a window with the given mark.

USAGE: `aerospace-marks focus <identifier>`

## Command: `list`

List all marks.

USAGE: `aerospace-marks list [--output <format>]`

### Output Formats

The `list` command supports multiple output formats via the `--output` (or `-o`) flag:

- **`text`** (default): Pipe-separated values with aligned columns
  ```
  mark-1 | 1 | Alacritty     | Alacritty      | _ | _
  mark-3 | 3 | Brave Browser | GitHub - Brave | _ | _
  ```

- **`json`**: JSON array of objects, easy to parse with `jq`
  ```json
  [
    {
      "mark": "mark-1",
      "window_id": 1,
      "app_name": "Alacritty",
      "window_title": "Alacritty",
      "workspace": "",
      "app_bundle_id": ""
    }
  ]
  ```

- **`csv`**: Comma-separated values with headers, compatible with csvkit
  ```csv
  mark,window_id,app_name,window_title,workspace,app_bundle_id
  mark-1,1,Alacritty,Alacritty,,
  ```

### Usage Examples

#### Text Format (default)
```bash
# Default behavior
aerospace-marks list

# Explicit text format
aerospace-marks list -o text

# Extract marks using awk
aerospace-marks list | awk -F'|' '{print $1}' | tr -d ' '
```

#### JSON Format
```bash
# Get all marks as JSON
aerospace-marks list -o json

# Extract all mark names
aerospace-marks list -o json | jq -r '.[].mark'

# Filter windows by app name
aerospace-marks list -o json | jq '.[] | select(.app_name == "Brave Browser")'

# Get window ID for a specific mark
aerospace-marks list -o json | jq -r '.[] | select(.mark == "mark-1") | .window_id'

# Count marked windows
aerospace-marks list -o json | jq 'length'
```

#### CSV Format
```bash
# Get all marks as CSV
aerospace-marks list -o csv

# Extract marks column using csvkit
aerospace-marks list -o csv | csvcut -c mark

# Filter using csvkit
aerospace-marks list -o csv | csvgrep -c app_name -m "Brave Browser"

# Process with awk (skip header)
aerospace-marks list -o csv | awk -F',' 'NR>1 {print $1}'
```

## Command: `unmark`

unmark will remove identifier from the list of current marks on a window. If identifier is omitted , all marks are removed.

USAGE: `aerospace-marks unmark [<identifier>]`

[read more](/docs/CMD_UNMARK.md)

## Command: `summon`

summon will bring the marked window to the current workspace.

USAGE: `aerospace-marks summon [<identifier>]`

## Command: `get`

Get a window by its mark (identifier) and prints the details in the following format:

USAGE: `aerospace-marks get <identifier> [flags]`

## Command: `info`

Show the current configurations and other info related

----

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
