package main

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/darkliquid/tag/internal/commands/add"
	"github.com/darkliquid/tag/internal/commands/list"
	"github.com/darkliquid/tag/internal/commands/remove"
	"github.com/darkliquid/tag/internal/commands/set"
	"github.com/darkliquid/tag/internal/commands/unset"
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
		list.NewListCommand(),
		add.NewAddCommand(),
		remove.NewRemoveCommand(),
		set.NewSetCommand(),
		unset.NewUnsetCommand(),
	)

	if err := fang.Execute(context.Background(), cmd); err != nil {
		os.Exit(1)
	}
}
