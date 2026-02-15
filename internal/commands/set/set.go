package set

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/xattr"
	"github.com/spf13/cobra"
)

func NewSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set [files] [tags]",
		Short:   "Replace tags on a file",
		Aliases: []string{"replace", "s"},
		Args:    cobra.MinimumNArgs(2),
		RunE:    runSet,
	}

	return cmd
}

func runSet(cmd *cobra.Command, args []string) error {
	files := args[:len(args)-1]
	tagsStr := args[len(args)-1]
	newTags := parseTags(tagsStr)

	for _, path := range files {
		if err := xattr.Set(path, "user.xdg.tags", []byte(strings.Join(newTags, ","))); err != nil {
			fmt.Fprintf(os.Stderr, "error setting tags for %q: %v\n", path, err)
			continue
		}
		fmt.Printf("Set tags on %s: %s\n", path, strings.Join(newTags, ", "))
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
