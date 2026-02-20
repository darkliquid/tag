# Tag

Tag is a simple utility for tagging files using filesystem extended attributes.

## Usage

```
tag [files...] [flags]
```

### Flags

- `-a, --add strings` - add comma separated tags to files
- `-d, --del strings` - delete comma separated tags from files
- `-s, --set strings` - set files tags to comma separated tags
- `-c, --clear` - clear tags on the files

### Examples

List tags on files:
```
tag file1.txt file2.txt
```

Add tags:
```
tag file.txt --add important,work
```

Remove tags:
```
tag file.txt --del old-tag
```

Replace all tags:
```
tag file.txt --set new-tag-1,new-tag-2
```

Clear all tags:
```
tag file.txt --clear
```

## Tag Syntax

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
as a comma separated list (CSV format).

## Limitations

Tags must be quoted if using characters with special meanings.
