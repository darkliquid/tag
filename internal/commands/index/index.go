package index

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/darkliquid/tag/internal/commands"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/xattr"
	"github.com/spf13/cobra"
)

func NewIndexCommand() *cobra.Command {
	var dbPath string

	cmd := &cobra.Command{
		Use:     "index [paths]",
		Short:   "Index your tagged files for searching",
		Aliases: []string{"idx", "i"},
		RunE:    runIndex,
	}

	cmd.Flags().StringVarP(&dbPath, "db", "d", commands.DBDir, "Path to the database file")

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		// Expand ~ in db path
		if strings.HasPrefix(dbPath, "~") {
			home, _ := os.UserHomeDir()
			dbPath = filepath.Join(home, dbPath[2:])
		}

		// Ensure directory exists
		if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
			return fmt.Errorf("failed to create database directory: %w", err)
		}
		return nil
	}

	return cmd
}

func runIndex(cmd *cobra.Command, args []string) error {
	dbPath := cmd.Flag("db").Value.String()

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Create tables if they don't exist
	if err := createTables(db); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	paths := args
	if len(paths) == 0 {
		paths = []string{"."}
	}

	var indexed int
	for _, path := range paths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error getting absolute path for %q: %v\n", path, err)
			continue
		}

		if err := filepath.Walk(absPath, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Fprintf(os.Stderr, "error accessing %q: %v\n", filePath, err)
				return nil
			}

			if info.IsDir() {
				return nil
			}

			if err := processFile(db, filePath); err != nil {
				fmt.Fprintf(os.Stderr, "error indexing %q: %v\n", filePath, err)
				return nil
			}

			indexed++
			return nil
		}); err != nil {
			return fmt.Errorf("error walking path %q: %w", path, err)
		}
	}

	fmt.Printf("Indexed %d files\n", indexed)
	return nil
}

func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		path TEXT UNIQUE NOT NULL,
		mtime INTEGER NOT NULL
	);

	CREATE TABLE IF NOT EXISTS tags (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		file_id INTEGER NOT NULL,
		tag TEXT NOT NULL,
		FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_tags_file ON tags(file_id);
	CREATE INDEX IF NOT EXISTS idx_tags_tag ON tags(tag);
	`

	_, err := db.Exec(query)
	return err
}

func processFile(db *sql.DB, path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	// Get current tags
	tagsBytes, err := xattr.Get(path, "user.xdg.tags")
	if err != nil {
		// Skip files without tags
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	tagsStr := strings.TrimSpace(string(tagsBytes))
	if tagsStr == "" {
		// No tags to index
		return nil
	}

	tagList := strings.Split(tagsStr, ",")

	// Check if file already exists
	var existingID int
	err = db.QueryRow("SELECT id FROM files WHERE path = ?", path).Scan(&existingID)

	if err == sql.ErrNoRows {
		// Insert new file
		result, err := db.Exec("INSERT INTO files (path, mtime) VALUES (?, ?)", path, info.ModTime().Unix())
		if err != nil {
			return err
		}
		fileID, _ := result.LastInsertId()

		// Insert tags
		for _, tag := range tagList {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				db.Exec("INSERT INTO tags (file_id, tag) VALUES (?, ?)", fileID, tag)
			}
		}
	} else if err == nil {
		// Update existing file
		db.Exec("UPDATE files SET mtime = ? WHERE id = ?", info.ModTime().Unix(), existingID)

		// Clear and re-insert tags
		db.Exec("DELETE FROM tags WHERE file_id = ?", existingID)
		for _, tag := range tagList {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				db.Exec("INSERT INTO tags (file_id, tag) VALUES (?, ?)", existingID, tag)
			}
		}
	}

	return nil
}
