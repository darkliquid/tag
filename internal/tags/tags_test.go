package tags

import (
	"bytes"
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"
)

func TestSet_Add(t *testing.T) {
	ts := make(Set)
	ts.Add("tag1", "tag2")

	if len(ts) != 2 {
		t.Errorf("expected 2 tags, got %d", len(ts))
	}

	if _, ok := ts["tag1"]; !ok {
		t.Errorf("tag1 not found")
	}

	if _, ok := ts["tag2"]; !ok {
		t.Errorf("tag2 not found")
	}
}

func TestSet_Add_Duplicate(t *testing.T) {
	ts := make(Set)
	ts.Add("tag1", "tag1")

	if len(ts) != 1 {
		t.Errorf("expected 1 tag, got %d", len(ts))
	}
}

func TestSet_Remove(t *testing.T) {
	ts := make(Set)
	ts.Add("tag1", "tag2", "tag3")
	ts.Remove("tag2")

	if len(ts) != 2 {
		t.Errorf("expected 2 tags, got %d", len(ts))
	}

	if _, ok := ts["tag1"]; !ok {
		t.Errorf("tag1 not found")
	}

	if _, ok := ts["tag2"]; ok {
		t.Errorf("tag2 should be removed")
	}

	if _, ok := ts["tag3"]; !ok {
		t.Errorf("tag3 not found")
	}
}

func TestSet_Remove_NonExistent(t *testing.T) {
	ts := make(Set)
	ts.Add("tag1")
	ts.Remove("nonexistent")

	if len(ts) != 1 {
		t.Errorf("expected 1 tag, got %d", len(ts))
	}
}

func TestSet_ToBytes(t *testing.T) {
	ts := make(Set)
	ts.Add("tag1", "tag2", "tag3")

	tagsBytes, err := ts.ToBytes()
	if err != nil {
		t.Fatalf("ToBytes failed: %v", err)
	}

	reader := csv.NewReader(bytes.NewReader(tagsBytes))
	tags, err := reader.Read()
	if err != nil {
		t.Fatalf("Failed to parse CSV: %v", err)
	}

	if len(tags) != 3 {
		t.Errorf("expected 3 tags, got %d", len(tags))
	}

	// Check that all expected tags are present
	expected := map[string]bool{"tag1": true, "tag2": true, "tag3": true}
	for _, tag := range tags {
		if !expected[tag] {
			t.Errorf("unexpected tag: %s", tag)
		}
	}
}

func TestSet_ToBytes_Empty(t *testing.T) {
	ts := make(Set)
	tagsBytes, err := ts.ToBytes()
	if err != nil {
		t.Fatalf("ToBytes failed: %v", err)
	}

	if tagsBytes != nil {
		t.Errorf("expected nil bytes for empty set, got %s", string(tagsBytes))
	}
}

func TestRead_NonExistentFile(t *testing.T) {
	_, err := Read("/nonexistent/file.txt")
	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestWrite_And_Read(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.txt")

	// Create test file
	if err := os.WriteFile(file, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Write tags
	ts := make(Set)
	ts.Add("tag1", "tag2", "tag with spaces")
	if err := Write(file, ts); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Read tags back
	readTags, err := Read(file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if readTags == nil {
		t.Fatalf("Read returned nil")
	}

	if len(readTags) != 3 {
		t.Errorf("expected 3 tags, got %d", len(readTags))
	}

	if _, ok := readTags["tag1"]; !ok {
		t.Errorf("tag1 not found")
	}

	if _, ok := readTags["tag2"]; !ok {
		t.Errorf("tag2 not found")
	}

	if _, ok := readTags["tag with spaces"]; !ok {
		t.Errorf("tag with spaces not found")
	}
}

func TestWrite_Clear(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.txt")

	// Create test file
	if err := os.WriteFile(file, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Write initial tags
	ts := make(Set)
	ts.Add("tag1", "tag2")
	if err := Write(file, ts); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Clear tags
	if err := Write(file, nil); err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	// Read tags back - should be empty
	readTags, err := Read(file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if readTags == nil {
		t.Fatalf("Read returned nil")
	}

	if len(readTags) != 0 {
		t.Errorf("expected 0 tags after clear, got %d", len(readTags))
	}
}
