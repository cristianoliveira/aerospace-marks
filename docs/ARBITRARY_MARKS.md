# Arbitrary Values for Marks - Design Proposal

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

#### 1. Extended Mark Command

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

#### 2. New Data Management Commands

```bash
# Set/update arbitrary data for existing mark
aerospace-marks set-data <identifier> key=value
aerospace-marks set-data <identifier> key1=value1 key2=value2

# Get arbitrary data from a mark
aerospace-marks get-data <identifier> [key]
aerospace-marks get-data <identifier>  # Returns all data as key=value pairs

# Remove data from a mark
aerospace-marks remove-data <identifier> key
aerospace-marks remove-data <identifier> key1 key2
aerospace-marks remove-data <identifier> --all  # Remove all data, keep the mark
```

#### 3. Enhanced Get Command

```bash
# Current functionality (unchanged)
aerospace-marks get <identifier>
aerospace-marks get <identifier> --window-id
aerospace-marks get <identifier> --window-title

# NEW: Get arbitrary data
aerospace-marks get <identifier> --data
aerospace-marks get <identifier> --data key
aerospace-marks get <identifier> --json  # Get all info as JSON
```

#### 4. Enhanced List Command

```bash
# Current functionality (unchanged)
aerospace-marks list

# NEW: Include arbitrary data in output
aerospace-marks list --with-data
aerospace-marks list --json  # JSON output with all data
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

#### JSON Format
```bash
$ aerospace-marks get myproject --json
{
  "mark": "myproject",
  "window_id": 12345,
  "window": {
    "app_name": "Code",
    "window_title": "MyProject.xcodeproj",
    "workspace": "1",
    "app_bundle_id": "com.apple.dt.Xcode"
  },
  "data": {
    "project_path": "/home/user/myproject",
    "build_config": "release",
    "last_commit": "abc123"
  }
}

$ aerospace-marks list --json
[
  {
    "mark": "myproject",
    "window_id": 12345,
    "window": {...},
    "data": {...}
  },
  {
    "mark": "browser",
    "window_id": 67890,
    "window": {...},
    "data": {
      "session": "work"
    }
  }
]
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

### Development Workflow
```bash
# Mark a development window with project metadata
aerospace-marks mark dev --data project=/home/user/myproject --data branch=feature/new-ui

# Later, retrieve project path for scripts
PROJECT_PATH=$(aerospace-marks get-data dev project)
cd "$PROJECT_PATH"

# Update branch info
aerospace-marks set-data dev branch=feature/updated-ui

# Get all development context
aerospace-marks get dev --json | jq .
```

### Session Management
```bash
# Mark browser windows with session info
aerospace-marks mark work-browser --data session=work --data profile=corporate
aerospace-marks mark personal-browser --data session=personal --data profile=default

# Script to restore work session
if aerospace-marks get work-browser --data session | grep -q "work"; then
    aerospace-marks focus work-browser
    echo "Restored work browser session"
fi
```

### Build System Integration
```bash
# Mark editor with build configuration
aerospace-marks mark editor --data build_config=debug --data target=ios

# Build script reads configuration
BUILD_CONFIG=$(aerospace-marks get-data editor build_config)
TARGET=$(aerospace-marks get-data editor target)
make build CONFIG=$BUILD_CONFIG TARGET=$TARGET
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