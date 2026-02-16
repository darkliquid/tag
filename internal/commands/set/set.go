package set

import (
	"fmt"
	"os"
	"strings"

	"github.com/darkliquid/tag/internal/commands"
	"github.com/spf13/cobra"
)

func NewSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set [files] --tags [tags]",
		Short:   "Replace tags on a file",
		Aliases: []string{"replace", "s"},
		Args:    cobra.MinimumNArgs(1),
		RunE:    runSet,
	}
	cmd.Flags().StringSlice("tags", nil, "comma separated tags")
	cmd.MarkFlagRequired("tags")

	return cmd
}

func runSet(cmd *cobra.Command, args []string) error {
	tags, err := cmd.Flags().GetStringSlice("tags")
	if err != nil {
		return err
	}

	ts := make(commands.TagSet).Add(tags...)

	for _, path := range args {
		if err := commands.SetTags(path, ts); err != nil {
			fmt.Fprintf(os.Stderr, "error setting tags for %q: %v\n", path, err)
			continue
		}
		fmt.Printf("Set tags on %s: %s\n", path, strings.Join(tags, ", "))
	}

	return nil
}
