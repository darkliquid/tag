package tags

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"iter"
	"maps"
	"slices"

	"github.com/pkg/xattr"
)

// ErrNoAttr is returned when a file has no extended attribute.
var ErrNoAttr = errors.New("no such attribute")

// Set is an unordered set of tags.
type Set map[string]struct{}

// Remove tags from the tagset.
func (ts Set) Remove(tags ...string) Set {
	for _, tag := range tags {
		delete(ts, tag)
	}
	return ts
}

// Add tags to the tagset.
func (ts Set) Add(tags ...string) Set {
	for _, tag := range tags {
		ts[tag] = struct{}{}
	}
	return ts
}

// Iter returns an iterator over the tagset.
func (ts Set) Iter() iter.Seq[string] {
	return maps.Keys(ts)
}

// ToBytes converts the tagset to CSV bytes.
func (ts Set) ToBytes() ([]byte, error) {
	if len(ts) == 0 {
		return nil, nil
	}
	buf := &bytes.Buffer{}
	w := csv.NewWriter(buf)
	if err := w.Write(slices.Collect(ts.Iter())); err != nil {
		return nil, err
	}
	w.Flush()
	return bytes.TrimSpace(buf.Bytes()), nil
}

// Read returns a tagset for the file at path.
func Read(path string) (Set, error) {
	current, err := xattr.Get(path, "user.xdg.tags")
	if err != nil {
		if errors.Is(err, xattr.ENOATTR) {
			return make(Set), nil
		}

		return nil, fmt.Errorf("xattr.Get: %w", err)
	}

	if string(current) == "" {
		return make(Set), nil
	}

	tags, err := csv.NewReader(bytes.NewReader(current)).Read()
	if err != nil {
		return nil, fmt.Errorf("csv.Read: %q %w", current, err)
	}

	return make(Set).Add(tags...), nil
}

// Write sets the tagset for the file at paths.
func Write(path string, tagset Set) error {
	buf := &bytes.Buffer{}
	w := csv.NewWriter(buf)
	if err := w.Write(slices.Collect(tagset.Iter())); err != nil {
		return fmt.Errorf("csv.Write: %w", err)
	}
	w.Flush()

	if err := xattr.Set(path, "user.xdg.tags", bytes.TrimSpace(buf.Bytes())); err != nil {
		return fmt.Errorf("xattr.Set: %w", err)
	}

	return nil
}
