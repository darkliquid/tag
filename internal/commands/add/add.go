package add

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/xattr"
	"github.com/spf13/cobra"
)

func NewAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add [files] [tags]",
		Short:   "Add tags to a file",
		Aliases: []string{"a"},
		Args:    cobra.MinimumNArgs(2),
		RunE:    runAdd,
	}

	return cmd
}

func runAdd(cmd *cobra.Command, args []string) error {
	files := args[:len(args)-1]
	tagsStr := args[len(args)-1]
	newTags := parseTags(tagsStr)

	for _, path := range files {
		current, err := xattr.Get(path, "user.xdg.tags")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading tags for %q: %v\n", path, err)
			continue
		}

		currentTags := parseTags(string(current))
		merged := mergeTags(currentTags, newTags)

		if err := xattr.Set(path, "user.xdg.tags", []byte(strings.Join(merged, ","))); err != nil {
			fmt.Fprintf(os.Stderr, "error setting tags for %q: %v\n", path, err)
			continue
		}
		fmt.Printf("Added tags to %s: %s\n", path, strings.Join(newTags, ", "))
	}

	return nil
}

func parseTags(s string) []string {
	if s == "" {
		return []string{}
	}
	tags := strings.Split(s, ",")
	for i, t := range tags {
		tags[i] = strings.TrimSpace(t)
	}
	return tags
}

func mergeTags(existing, new []string) []string {
	existingMap := make(map[string]bool)
	for _, t := range existing {
		if t != "" {
			existingMap[t] = true
		}
	}
	for _, t := range new {
		if t != "" {
			existingMap[t] = true
		}
	}

	result := make([]string, 0, len(existingMap))
	for t := range existingMap {
		result = append(result, t)
	}
	return result
}
