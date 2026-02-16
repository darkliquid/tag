package commands

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"iter"
	"maps"
	"os"
	"path/filepath"
	"slices"

	"github.com/pkg/xattr"
)

// TagSet is an unordered set of tags.
type TagSet map[string]struct{}

// Remove tags from the tagset.
func (ts TagSet) Remove(tags ...string) TagSet {
	for _, tag := range tags {
		delete(ts, tag)
	}
	return ts
}

// Add tags to the tagset.
func (ts TagSet) Add(tags ...string) TagSet {
	for _, tag := range tags {
		ts[tag] = struct{}{}
	}
	return ts
}

// Iter returns an iterator over the tagset.
func (ts TagSet) Iter() iter.Seq[string] {
	return maps.Keys(ts)
}

// GetTags returns a tagset for the file at path.
func GetTags(path string) TagSet {
	current, err := xattr.Get(path, "user.xdg.tags")
	if err != nil {
		if errors.Is(err, xattr.ENOATTR) {
			return make(TagSet)
		}

		fmt.Fprintf(os.Stderr, "error reading tags for %q: %v\n", path, err)
		return nil
	}

	tags, err := csv.NewReader(bytes.NewReader(current)).Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error decoding tags for %q: %v\n", path, err)
		return nil
	}

	return make(TagSet).Add(tags...)
}

// SetTags sets the tagset for the file at paths.
func SetTags(path string, tagset TagSet) error {
	buf := &bytes.Buffer{}
	w := csv.NewWriter(buf)
	if err := w.Write(slices.Collect(tagset.Iter())); err != nil {
		return err
	}
	w.Flush()

	return xattr.Set(path, "user.xdg.tags", bytes.TrimSpace(buf.Bytes()))
}

// ExpandedPaths returns all the paths after glob expansion.
func ExpandedPaths(paths ...string) ([]string, error) {
	expanded := make(map[string]struct{})
	for _, path := range paths {
		matches, err := filepath.Glob(path)
		if err != nil {
			return nil, err
		}

		for _, match := range matches {
			expanded[match] = struct{}{}
		}
	}

	return slices.Collect(maps.Keys(expanded)), nil
}
