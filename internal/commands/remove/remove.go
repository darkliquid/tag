package remove

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/xattr"
	"github.com/spf13/cobra"
)

func NewRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove [files] [tags]",
		Short:   "Remove tags from a file",
		Aliases: []string{"delete", "del", "rm", "r", "d"},
		Args:    cobra.MinimumNArgs(2),
		RunE:    runRemove,
	}

	return cmd
}

func runRemove(cmd *cobra.Command, args []string) error {
	files := args[:len(args)-1]
	tagsStr := args[len(args)-1]
	tagsToRemove := parseTags(tagsStr)

	for _, path := range files {
		current, err := xattr.Get(path, "user.xdg.tags")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading tags for %q: %v\n", path, err)
			continue
		}

		currentTags := parseTags(string(current))
		remaining := removeTags(currentTags, tagsToRemove)

		if err := xattr.Set(path, "user.xdg.tags", []byte(strings.Join(remaining, ","))); err != nil {
			fmt.Fprintf(os.Stderr, "error setting tags for %q: %v\n", path, err)
			continue
		}
		fmt.Printf("Removed tags from %s: %s\n", path, strings.Join(tagsToRemove, ", "))
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

func removeTags(existing, toRemove []string) []string {
	removeMap := make(map[string]bool)
	for _, t := range toRemove {
		if t != "" {
			removeMap[t] = true
		}
	}

	result := make([]string, 0)
	for _, t := range existing {
		if t != "" && !removeMap[t] {
			result = append(result, t)
		}
	}
	return result
}
