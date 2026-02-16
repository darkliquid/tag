package main

import (
	"context"
	"fmt"
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
		Use:   "tag [command]",
		Short: "Tag your files",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.LocalFlags().Lookup("tags") == nil {
				return nil
			}

			tags, err := cmd.LocalFlags().GetStringSlice("tags")
			if err != nil {
				return err
			}

			if len(tags) == 0 {
				return fmt.Errorf("no tags provided")
			}

			return nil
		},
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
