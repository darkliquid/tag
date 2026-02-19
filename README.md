# Tag

Tag is a simple utility for tagging files and searching for files by tag.

## Usage

- `tag add [files...] --tags tag1` - adds tag1 to the given files
- `tag remove [files...] --tags tag1` - removes tag1 from the given files
- `tag set [files...] --tags tag1` - sets the tags for the given files to tag1
- `tag unset [files...]` - removes all tags for the given files
- `tag list [files...]` - lists all tags of the given files

## Tag Restrictions

Tags can use any characters, including spaces or commas. However, they **must**
be quoted if using spaces, commas or quotes.

Basically just treat the tag list as a CSV record.

For example:

- `tag name` - bad!
- `"tag name"` - good!
- `tag-name` - good!
- `has-a-"-character` - bad!
- `'has-a-"-character'` - good!
- `one,two,a b c` - bad!
- `one,two,"a b c"` - good!

## Implementation

Tagging is implemented using filesystem extended attributes. Naturally this means
that this only works for files on filesystems that _support_ extended attributes.

Specifically, it uses the attribute `user.xdg.tags` and stores tags within that
as a comma separated list.

For indexing and search, a simple database is used.

## Limitations

Tags must be quoted if using characters with special meanings.
