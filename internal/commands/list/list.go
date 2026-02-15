package list

import (
	"fmt"
	"os"

	"github.com/pkg/xattr"
	"github.com/spf13/cobra"
)

func NewListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list [files]",
		Short:   "List tags on files",
		Aliases: []string{"ls", "l"},
		Args:    cobra.MinimumNArgs(1),
		RunE:    runList,
	}

	return cmd
}

func runList(cmd *cobra.Command, args []string) error {
	for _, path := range args {
		tagsBytes, err := xattr.Get(path, "user.xdg.tags")
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("%s: (no tags)\n", path)
				continue
			}
			fmt.Fprintf(os.Stderr, "error reading tags for %q: %v\n", path, err)
			continue
		}
		tags := string(tagsBytes)
		if tags == "" {
			fmt.Printf("%s: (no tags)\n", path)
		} else {
			fmt.Printf("%s: %s\n", path, tags)
		}
	}
	return nil
}
