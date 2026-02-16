package unset

import (
	"fmt"
	"os"

	"github.com/pkg/xattr"
	"github.com/spf13/cobra"
)

func NewUnsetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "unset [files]",
		Short:   "Clear tags on a file",
		Aliases: []string{"clear", "u", "c"},
		Args:    cobra.MinimumNArgs(1),
		RunE:    runUnset,
	}

	return cmd
}

func runUnset(cmd *cobra.Command, args []string) error {
	files := args

	for _, path := range files {
		if err := xattr.Set(path, "user.xdg.tags", nil); err != nil {
			fmt.Fprintf(os.Stderr, "error clearing tags for %q: %v\n", path, err)
			continue
		}
		fmt.Printf("Cleared tags on %s\n", path)
	}

	return nil
}
