# Command: `unmark`

Unmark one or more windows by identifier.

unmark [<identifier>]

unmark will remove identifier from the list of current marks on a window. If identifier is omitted , all marks are removed.

## Usage

```bash
aerospace-marks unmark [<identifier>]

aerospace-marks unmark foo
# Will unmark the current focused window with the identifier "foo"

aerospace-marks unmark
# Will unmark all windows 
```
