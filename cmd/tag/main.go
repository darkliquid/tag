package main

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "tag",
		Short: "Tag your files",
	}
	cmd.AddGroup(
		&cobra.Group{
			ID:    "tagging",
			Title: "Tagging",
		},
		&cobra.Group{
			ID:    "search",
			Title: "Searching",
		},
	)
	cmd.AddCommand(
		&cobra.Command{
			Use:     "list [files]",
			Short:   "List tags on files",
			Aliases: []string{"ls", "l"},
			GroupID: "tagging",
		},
		&cobra.Command{
			Use:     "add [files] [tags]",
			Short:   "Add tags to a file",
			Aliases: []string{"a"},
			GroupID: "tagging",
		},
		&cobra.Command{
			Use:     "remove [files] [tags]",
			Short:   "Remove tags from a file",
			Aliases: []string{"delete", "del", "rm", "r", "d"},
			GroupID: "tagging",
		},
		&cobra.Command{
			Use:     "set [files] [tags]",
			Short:   "Replace tags on a file",
			Aliases: []string{"replace", "s"},
			GroupID: "tagging",
		},
		&cobra.Command{
			Use:     "unset [files] [tags]",
			Short:   "Clear tags on a file",
			Aliases: []string{"clear", "u", "c"},
			GroupID: "tagging",
		},
		&cobra.Command{
			Use:     "index [paths]",
			Short:   "Index your tagged files for searching",
			Aliases: []string{"idx", "i"},
			GroupID: "search",
		},
		&cobra.Command{
			Use:     "find [tags]",
			Short:   "Find your tagged files by tag",
			Aliases: []string{"search", "f"},
			GroupID: "search",
		},
	)

	if err := fang.Execute(context.Background(), cmd); err != nil {
		os.Exit(1)
	}
}
