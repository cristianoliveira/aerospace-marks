# Arbitrary Values for Marks - Design Proposal

## Ideas

### Allow add extra data to an already existing mark by window id

```bash
aerospace-marks mark --data key=value
# This will add key=value to the existing mark of the current window
```

## Context

Currently, aerospace-marks stores only a `window_id` associated with each mark identifier. The marks are stored in a simple SQLite table structure:

```sql
CREATE TABLE IF NOT EXISTS marks (
    window_id INTEGER NOT NULL,
    mark TEXT NOT NULL UNIQUE
);
```

This feature proposal aims to extend marks functionality to allow storing arbitrary key-value data associated with each mark, enabling advanced scripting and automation use cases while maintaining backward compatibility.

## Motivation

The arbitrary marks feature would enable users to:

1. **Enhanced Scripting**: Store metadata like project paths, configurations, or state information with window marks
2. **Context-Aware Automation**: Scripts can retrieve and act upon stored metadata
3. **Extended Window Management**: Associate custom attributes beyond just the window reference
4. **Workflow Enhancement**: Create more sophisticated window management workflows

## Proposed Interface

### Core Commands

#### Extended Mark Command

```bash
# Current functionality (unchanged)
aerospace-marks mark <identifier>
aerospace-marks mark --add <identifier>
aerospace-marks mark --toggle <identifier>

# NEW: Store arbitrary data with marks
aerospace-marks mark <identifier> --data key=value
aerospace-marks mark <identifier> --data key1=value1 --data key2=value2
aerospace-marks mark <identifier> --data key=value --add
```

**Flags:**
- `--data key=value`: Store key-value pairs with the mark
- Multiple `--data` flags can be specified in a single command
- Works with existing `--add`, `--toggle`, `--replace` flags
- `--silent` flag continues to work as expected

**Behavior:**
- When using `--data` with default (replace) mode: replaces all existing data
- When using `--data` with `--add`: merges new data with existing data (overwrites keys if they exist)
- When using `--data` with `--toggle`: if mark exists, removes mark and all its data; if mark doesn't exist, creates mark with provided data

#### Enhanced Get Command

```bash
# Current functionality (unchanged)
aerospace-marks get <identifier>
aerospace-marks get <identifier> --window-id
aerospace-marks get <identifier> --window-title

# NEW: Get arbitrary data
aerospace-marks get <identifier> --data
aerospace-marks get <identifier> --data key
```

#### Enhanced List Command

```bash
# Current functionality (unchanged)
aerospace-marks list

# NEW: Include arbitrary data in output
aerospace-marks list --with-data
```

### Output Formats

#### Text Format (Default)
```bash
$ aerospace-marks get myproject --data
project_path=/home/user/myproject
build_config=release
last_commit=abc123

$ aerospace-marks get myproject --data project_path
/home/user/myproject

$ aerospace-marks list --with-data
myproject|12345|Code|MyProject.xcodeproj|1|com.apple.dt.Xcode|project_path=/home/user/myproject,build_config=release
browser|67890|Safari|GitHub|2|com.apple.Safari|session=work
```

## Database Schema Changes

### New Table Structure

```sql
-- Keep existing marks table for window associations
CREATE TABLE IF NOT EXISTS marks (
    window_id INTEGER NOT NULL,
    mark TEXT NOT NULL UNIQUE
);

-- New table for arbitrary data
CREATE TABLE IF NOT EXISTS mark_data (
    mark TEXT NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    PRIMARY KEY (mark, key),
    FOREIGN KEY (mark) REFERENCES marks(mark) ON DELETE CASCADE
);
```

### Migration Strategy

1. **Backward Compatibility**: Existing marks continue to work without any data
2. **Automatic Migration**: New table created automatically via existing goose migration system
3. **Data Integrity**: Foreign key constraint ensures data is cleaned up when marks are deleted

## Usage Examples

### Window Management with aerospace-scratchpad
```bash
# Mark the chatgpt window with session info (Title contains "questions")
aerospace-marks mark gpt --data session=questions

# Later, with a shortcut, restore the window and set its title
aerospace-scratchpad show "$(aerospace-marks get gpt -a)" \
                  -F window-title="$(aerospace-marks get gpt --data session)"
```

## Implementation Considerations

### Storage Interface Changes

```go
type MarkStorage interface {
    // Existing methods (unchanged)
    AddMark(id int, mark string) error
    GetMarks() ([]queries.Mark, error)
    // ... other existing methods

    // New data methods
    SetMarkData(mark string, key string, value string) error
    GetMarkData(mark string, key string) (string, error)
    GetAllMarkData(mark string) (map[string]string, error)
    DeleteMarkData(mark string, key string) error
    DeleteAllMarkData(mark string) error
}
```

### Command Structure

- New commands follow existing patterns established in the codebase
- Use cobra.Command structure consistent with existing commands
- Maintain existing error handling and output patterns
- Support existing `--silent` flag across all new commands

### Data Validation

- Keys and values are stored as strings (consistent with current mark identifiers)
- Key names must be valid (alphanumeric + underscore, no spaces)
- Value can be any string (including empty)
- No arbitrary limits on data size, but practical limits for CLI usage

### Error Handling

- Setting data on non-existent mark creates the mark automatically
- Getting data from non-existent mark returns appropriate error
- Getting non-existent key returns empty value with success (bash-friendly)
- Proper error messages maintain consistency with existing commands

## Backward Compatibility

✅ **Full Backward Compatibility Maintained**

- All existing commands work unchanged
- Existing marks continue to function without any data
- Existing scripts and configurations require no changes
- Database migrations are automatic and safe

## Future Considerations

1. **Data Types**: Currently all data stored as strings. Future enhancement could add type support
2. **Data Validation**: Could add schema validation for specific use cases
3. **Data Export/Import**: Bulk operations for backing up/restoring mark data
4. **Data Querying**: Advanced filtering and searching capabilities
5. **Integration**: Hooks for external tools to automatically populate mark data

## Summary

This proposal extends aerospace-marks with arbitrary data storage while maintaining full backward compatibility. The interface follows existing patterns and enables powerful scripting and automation workflows. The implementation leverages the existing SQLite foundation and migration system for a robust, maintainable solution.
