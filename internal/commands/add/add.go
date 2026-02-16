package add

import (
	"fmt"
	"os"
	"strings"

	"github.com/darkliquid/tag/internal/commands"
	"github.com/spf13/cobra"
)

func NewAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add [files]",
		Short:   "Add tags to a file",
		Aliases: []string{"a"},
		Args:    cobra.MinimumNArgs(2),
		RunE:    runAdd,
	}
	cmd.Flags().StringSlice("tags", nil, "comma separated tags")

	return cmd
}

func runAdd(cmd *cobra.Command, args []string) error {
	files := args[:len(args)-1]
	tags, err := cmd.Flags().GetStringSlice("tags")
	if err != nil {
		return err
	}

	for _, path := range files {
		if err := commands.SetTags(path, commands.GetTags(path).Add(tags...)); err != nil {
			fmt.Fprintf(os.Stderr, "error setting tags for %q: %v\n", path, err)
			continue
		}

		fmt.Printf("Added tags to %s: %s\n", path, strings.Join(tags, ", "))
	}

	return nil
}
