package main

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/darkliquid/tag/internal/tags"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:           "tag [files] [--add tags] [--remove tags] [--set tags] [--clear]",
		Short:         "Tag your files",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.ArbitraryArgs,
		RunE:          runTag,
	}
	cmd.Flags().StringSliceP("add", "a", nil, "add comma separated tags to files")
	cmd.Flags().StringSliceP("del", "d", nil, "delete comma separated tags from files")
	cmd.Flags().StringSliceP("set", "s", nil, "set files tags to comma separated tags")
	cmd.Flags().BoolP("clear", "c", false, "clear tags on the files")
	cmd.MarkFlagsMutuallyExclusive("clear", "set")
	cmd.MarkFlagsMutuallyExclusive("clear", "add")
	cmd.MarkFlagsMutuallyExclusive("clear", "del")
	cmd.MarkFlagsMutuallyExclusive("set", "add")
	cmd.MarkFlagsMutuallyExclusive("set", "del")

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

	if err := fang.Execute(context.Background(), cmd); err != nil {
		os.Exit(1)
	}
}

func runTag(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("at least one file must be specified")
	}

	// List tags on files when no flags specified.
	if cmd.Flags().NFlag() == 0 {
		for _, path := range args {
			fileTags, err := tags.Read(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error reading tags for %q: %v\n", path, err)
				continue
			}

			if fileTags == nil {
				continue
			}

			tagsBytes, err := fileTags.ToBytes()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error converting tags for %q: %v\n", path, err)
				continue
			}

			if len(tagsBytes) == 0 {
				fmt.Printf("%s: (no tags)\n", path)
			} else {
				fmt.Printf("%s: %s\n", path, string(tagsBytes))
			}
		}

		return nil
	}

	// Clear all the tags from the given files.
	clear, err := cmd.Flags().GetBool("clear")
	if err != nil {
		return err
	}
	if clear {
		for _, path := range args {
			if err := tags.Write(path, nil); err != nil {
				fmt.Fprintf(os.Stderr, "error clearing tags for %q: %v\n", path, err)
				continue
			}
		}
		return nil
	}

	// Replace the tags on the given files with the given tags.
	replacements, err := cmd.Flags().GetStringSlice("set")
	if err != nil {
		return err
	}
	if len(replacements) > 0 {
		ts := make(tags.Set)
		ts.Add(replacements...)
		for _, path := range args {
			if err := tags.Write(path, ts); err != nil {
				fmt.Fprintf(os.Stderr, "error clearing tags for %q: %v\n", path, err)
				continue
			}
		}
		return nil
	}

	// Add all the tags to the given files.
	additions, err := cmd.Flags().GetStringSlice("add")
	if err != nil {
		return err
	}

	// Remove all the tags from the given files.
	deletions, err := cmd.Flags().GetStringSlice("del")
	if err != nil {
		return err
	}

	for _, path := range args {
		ts, err := tags.Read(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading tags for %q: %v\n", path, err)
			continue
		}

		ts.Add(additions...)
		ts.Remove(deletions...)
		if err := tags.Write(path, ts); err != nil {
			fmt.Fprintf(os.Stderr, "error updating tags for %q: %v\n", path, err)
			continue
		}
	}

	return nil
}
