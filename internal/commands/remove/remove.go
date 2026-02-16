package remove

import (
	"fmt"
	"os"
	"strings"

	"github.com/darkliquid/tag/internal/commands"
	"github.com/spf13/cobra"
)

func NewRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove [files] --tags [tags]",
		Short:   "Remove tags from a file",
		Aliases: []string{"delete", "del", "rm", "r", "d"},
		Args:    cobra.MinimumNArgs(1),
		RunE:    runRemove,
	}
	cmd.Flags().StringSlice("tags", nil, "comma separated tags")
	cmd.MarkFlagRequired("tags")

	return cmd
}

func runRemove(cmd *cobra.Command, args []string) error {
	tags, err := cmd.Flags().GetStringSlice("tags")
	if err != nil {
		return err
	}

	for _, path := range args {
		if err := commands.SetTags(path, commands.GetTags(path).Remove(tags...)); err != nil {
			fmt.Fprintf(os.Stderr, "error setting tags for %q: %v\n", path, err)
			continue
		}
		fmt.Printf("Removed tags from %s: %s\n", path, strings.Join(tags, ", "))
	}

	return nil
}
