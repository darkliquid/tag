# Tag

Tag is a simple utility for tagging files and searching for files by tag.

## Usage

 - `tag [files...]` - lists current tags on specified files
 - `tag [files...] -- tag1,tag2,tag3` - sets (replaces) the tags on the specified files
 - `tag [files...] -- +tag1,-tag2` - adds tag1 and removes tag2 from the specified files
 - `tag index [path]` - updates a tag-based search index by scanning path
 - `tag search tag1` - searches for files with tag1
 - `tag search tag1,tag2,-tag3` - searches for files with tag1 and tag2 but not with tag3
 - `tag search tag1,tag2 tag1,tag3 tag2,-tag1` - searches for files with tag1 and tag2, or files with tag1 and tag3, or files with tag2 but not tag1

## Naming & Searching

Tags can use any characters, including spaces or commas. However, they **must** be
quoted if using spaces, commas, quotes or if they begin with a `-` or `+` character.

For example:

 - `tag name` - bad!
 - `"tag name" - good!
 - `+ or -` - bad!
 - `"+ or -` - good!
 - `tag-name` - good!
 - `has-a-"-character` - bad!
 - `'has-a-"-character'` - good!

Searching uses a simple syntax. Spaces indicate OR matching, `,` indicate AND matching.
Parentheses can also be used to group matches together.

## Implementation

Tagging is implemented using filesystem extended attributes. Naturally this means
that this only works for files on filesystems that _support_ extended attributes.

Specifically, it uses the attribute `user.xdg.tags` and stores tags within that as
a comma separated list.

For indexing and search, a simple database is used.

## Limitations

Tags must be quoted if using characters with special meanings.

Moving files can desync the search index (the index will not be updated with the
new file name/location). There are _some_ ways to try and mitigate this, such as
creating inotify subscriptions to detect when files are moved/renamed/deleted but
that adds overhead and also isn't fully reliable (such as cases when moving files
between different filesystems, or moving one tagged file over the top of another
tagged file). For now, reindexing is easier, though sub-optimal.
