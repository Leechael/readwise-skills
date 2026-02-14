package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/leechael/readwise-skills/internal/client"
)

func newBookCmd() *cobra.Command {
	bookCmd := &cobra.Command{
		Use:   "book",
		Short: "Manage books/sources",
	}

	bookCmd.AddCommand(newBookListCmd())
	bookCmd.AddCommand(newBookGetCmd())
	bookCmd.AddCommand(newBookTagCmd())

	return bookCmd
}

func newBookListCmd() *cobra.Command {
	var params client.BookListParams

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List books/sources",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				return err
			}
			result, err := c.BookList(params)
			if err != nil {
				return err
			}
			return getFormatter().Print(result)
		},
	}

	cmd.Flags().IntVar(&params.PageSize, "page-size", 0, "Number of results per page")
	cmd.Flags().IntVar(&params.Page, "page", 0, "Page number")
	cmd.Flags().StringVar(&params.Category, "category", "", "Filter by category (books|articles|tweets|podcasts)")
	cmd.Flags().StringVar(&params.Source, "source", "", "Filter by source")
	cmd.Flags().StringVar(&params.UpdatedLT, "updated-before", "", "Filter by updated before (ISO 8601)")
	cmd.Flags().StringVar(&params.UpdatedGT, "updated-after", "", "Filter by updated after (ISO 8601)")
	cmd.Flags().StringVar(&params.LastHighlightAtLT, "last-highlight-before", "", "Filter by last highlight before (ISO 8601)")
	cmd.Flags().StringVar(&params.LastHighlightAtGT, "last-highlight-after", "", "Filter by last highlight after (ISO 8601)")

	return cmd
}

func newBookGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a book by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid book ID: %s", args[0])
			}
			c, err := getClient()
			if err != nil {
				return err
			}
			result, err := c.BookGet(id)
			if err != nil {
				return err
			}
			return getFormatter().Print(result)
		},
	}
}

func newBookTagCmd() *cobra.Command {
	tagCmd := &cobra.Command{
		Use:   "tag",
		Short: "Manage book tags",
	}

	tagCmd.AddCommand(newBookTagListCmd())
	tagCmd.AddCommand(newBookTagAddCmd())
	tagCmd.AddCommand(newBookTagUpdateCmd())
	tagCmd.AddCommand(newBookTagDeleteCmd())

	return tagCmd
}

func newBookTagListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list <book-id>",
		Short: "List tags for a book",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid book ID: %s", args[0])
			}
			c, err := getClient()
			if err != nil {
				return err
			}
			result, err := c.BookTagList(id)
			if err != nil {
				return err
			}
			return getFormatter().Print(result)
		},
	}
}

func newBookTagAddCmd() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "add <book-id>",
		Short: "Add a tag to a book",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid book ID: %s", args[0])
			}
			c, err := getClient()
			if err != nil {
				return err
			}
			result, err := c.BookTagAdd(id, name)
			if err != nil {
				return err
			}
			return getFormatter().Print(result)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Tag name")
	cmd.MarkFlagRequired("name")

	return cmd
}

func newBookTagUpdateCmd() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "update <book-id> <tag-id>",
		Short: "Update a tag on a book",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			bookID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid book ID: %s", args[0])
			}
			tagID, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid tag ID: %s", args[1])
			}
			c, err := getClient()
			if err != nil {
				return err
			}
			result, err := c.BookTagUpdate(bookID, tagID, name)
			if err != nil {
				return err
			}
			return getFormatter().Print(result)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "New tag name")
	cmd.MarkFlagRequired("name")

	return cmd
}

func newBookTagDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <book-id> <tag-id>",
		Short: "Delete a tag from a book",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			bookID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid book ID: %s", args[0])
			}
			tagID, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid tag ID: %s", args[1])
			}
			c, err := getClient()
			if err != nil {
				return err
			}
			if err := c.BookTagDelete(bookID, tagID); err != nil {
				return err
			}
			getFormatter().PrintMessage("Tag deleted")
			return nil
		},
	}
}
