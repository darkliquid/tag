package find

import (
	"database/sql"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/darkliquid/tag/internal/commands"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

func NewFindCommand() *cobra.Command {
	var dbPath string

	cmd := &cobra.Command{
		Use:     "find [tags]",
		Short:   "Find your tagged files by tag",
		Aliases: []string{"search", "f"},
		Args:    cobra.MinimumNArgs(1),
		RunE:    runFind,
	}

	cmd.Flags().StringVarP(&dbPath, "db", "d", commands.DBDir, "Path to the database file")

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		dbPath := cmd.Flag("db").Value.String()
		if strings.HasPrefix(dbPath, "~") {
			home, _ := os.UserHomeDir()
			dbPath = strings.Replace(dbPath, "~", home, 1)
		}

		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			return fmt.Errorf("database not found at %s - run 'tag index' first", dbPath)
		}
		// Update the flag value with the expanded path
		cmd.Flags().Set("db", dbPath)
		return nil
	}

	return cmd
}

func runFind(cmd *cobra.Command, args []string) error {
	dbPath := cmd.Flag("db").Value.String()

	db, err := openDB(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Parse all search groups (separated by spaces)
	groups := args
	results := make(map[string]bool)

	for _, group := range groups {
		groupResults, err := executeGroupQuery(db, group)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error executing query for %q: %v\n", group, err)
			continue
		}
		for _, path := range groupResults {
			results[path] = true
		}
	}

	// Output sorted results
	var paths []string
	for path := range results {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	for _, path := range paths {
		fmt.Println(path)
	}

	return nil
}

func openDB(path string) (*sql.DB, error) {
	return sql.Open("sqlite3", path)
}

func executeGroupQuery(db *sql.DB, queryStr string) ([]string, error) {
	// Split by comma for AND conditions, space would be handled at higher level
	tokens := strings.Split(queryStr, ",")

	var selectParts []string
	var args []interface{}

	for _, token := range tokens {
		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}

		var condition string
		var param string

		if strings.HasPrefix(token, "-") {
			// NOT condition
			param = strings.TrimPrefix(token, "-")
			condition = "NOT EXISTS (SELECT 1 FROM tags t2 WHERE t2.file_id = f.id AND t2.tag = ?)"
		} else {
			// Positive condition
			param = token
			condition = "EXISTS (SELECT 1 FROM tags t1 WHERE t1.file_id = f.id AND t1.tag = ?)"
		}

		selectParts = append(selectParts, condition)
		args = append(args, param)
	}

	if len(selectParts) == 0 {
		return nil, nil
	}

	sqlQuery := fmt.Sprintf(`
		SELECT DISTINCT f.path FROM files f
		WHERE %s
	`, strings.Join(selectParts, " AND "))

	rows, err := db.Query(sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return nil, err
		}
		results = append(results, path)
	}

	return results, rows.Err()
}
