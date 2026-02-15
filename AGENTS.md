# Agent Guide - Tag Project

## Project Overview

This is a Go CLI utility for tagging files and searching by tag. The project is currently in early/initial scaffolding stage.

## Current State

- **Project type**: Go CLI application
- **Module**: `github.com/darkliquid/tag`
- **Go version**: 1.25.0
- **Primary framework**: Cobra (with fang for execution)
- **Source files**: Only `cmd/tag/main.go` exists (scaffolding only)

## Essential Commands

```bash
# Build
go build -o tag ./cmd/tag

# Run
go run ./cmd/tag

# Test
# No tests currently exist

# Lint/typecheck
# No linting configuration currently exists
```

## Code Organization

```
cmd/
└── tag/
    └── main.go              # CLI entry point with command wiring

internal/
└── commands/
    ├── commands.go          # Common interfaces and types
    ├── list/
    │   └── list.go          # List tags on files
    ├── add/
    │   └── add.go           # Add tags to files
    ├── remove/
    │   └── remove.go        # Remove tags from files
    ├── set/
    │   └── set.go           # Set/replace tags on files
    ├── unset/
    │   └── unset.go         # Clear all tags on files
    ├── index/
    │   └── index.go         # Build search index
    └── find/
        └── find.go          # Search for tagged files
```

## Command Structure

### Tagging Commands (Group: "tagging")
| Command | Aliases | Description |
|---------|---------|-------------|
| `tag list [files]` | `ls`, `l` | List tags on files |
| `tag add [files] [tags]` | `a` | Add tags to files |
| `tag remove [files] [tags]` | `delete`, `del`, `rm`, `r`, `d` | Remove tags from files |
| `tag set [files] [tags]` | `replace`, `s` | Replace tags on files |
| `tag unset [files] [tags]` | `clear`, `u`, `c` | Clear all tags on files |

### Search Commands (Group: "search")
| Command | Aliases | Description |
|---------|---------|-------------|
| `tag index [paths]` | `idx`, `i` | Index tagged files for searching |
| `tag find [tags]` | `search`, `f` | Find files by tag |

## Implementation Notes

### Tagging Backend
- Uses filesystem extended attributes
- Attribute name: `user.xdg.tags`
- Tags stored as comma-separated values
- Use `github.com/pkg/xattr` package for xattr operations

### Indexing & Search
- Uses SQLite database for search index
- Database path: `~/.local/share/tag/index.db`
- Tables: `files` (id, path, mtime) and `tags` (id, file_id, tag)
- Indexes: `idx_tags_file` and `idx_tags_tag`
- Files moved/deleted can desync the index (reindexing required)

## Dependencies

Key dependencies (from `go.mod`):
- `github.com/spf13/cobra` - CLI framework
- `github.com/charmbracelet/fang` - Command execution
- `github.com/pkg/xattr` - Extended attribute access
- `github.com/mattn/go-sqlite3` - SQLite database driver
- `charm.land/lipgloss/v2` - Terminal styling

## Gotchas & Conventions

- **Tag naming**: Tags can contain any characters including spaces/commas, but **must be quoted** if they contain spaces, commas, quotes, or start with `+` or `-`
- **Search syntax**: Spaces indicate OR matching, commas indicate AND matching. Parentheses can be used for grouping.
  - Example: `tag find "tag1,tag2 tag3"` finds files with (tag1 AND tag2) OR (tag3)
  - Example: `tag find "tag1,-tag2"` finds files with tag1 but NOT tag2
- **Filesystem requirement**: Extended attributes only work on filesystems that support them (ext4, xfs, btrfs, etc.)
- **Index limitations**: Moving files can break search index (no automatic sync). Run `tag index` to rebuild.
- **Testing**: /tmp often doesn't support xattrs. Use a proper ext4/xfs/btrfs filesystem for testing.

## Database Schema

```sql
CREATE TABLE files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    path TEXT UNIQUE NOT NULL,
    mtime INTEGER NOT NULL
);

CREATE TABLE tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    file_id INTEGER NOT NULL,
    tag TEXT NOT NULL,
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
);

CREATE INDEX idx_tags_file ON tags(file_id);
CREATE INDEX idx_tags_tag ON tags(tag);
```

## Testing Tips

1. Create a test directory on a local ext4/xfs/btrfs filesystem
2. Use `setfattr` and `getfattr` to manually set xattrs for testing
3. Example:
   ```bash
   mkdir -p /tmp/testtag
   echo "content" > /tmp/testtag/file.txt
   setfattr -n user.xdg.tags -v "important,work" /tmp/testtag/file.txt
   tag list /tmp/testtag/file.txt
   ```

## Pending Tasks

The following functionality needs implementation:
1. Indexing functionality ✅
2. Search functionality ✅
3. Database implementation for search index ✅
4. Extended attribute handling utilities ✅

## Command Patterns

Each command follows this structure:

```go
package <commandname>

import "github.com/spf13/cobra"

func New<CommandName>Command() *cobra.Command {
    cmd := &cobra.Command{
        Use:     "<command> [args]",
        Short:   "Description",
        Aliases: []string{"a", "b"},
        Args:    cobra.MinimumNArgs(1),  // or cobra.ExactArgs(2), etc.
        RunE:    runCommand,
    }
    // Add flags if needed
    cmd.Flags().StringVarP(&opt, "flag", "f", "default", "help")
    return cmd
}

func runCommand(cmd *cobra.Command, args []string) error {
    // Command logic here
    return nil
}
```
