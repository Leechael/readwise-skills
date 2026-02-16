package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/leechael/readwise-skills/internal/client"
	"github.com/leechael/readwise-skills/internal/model"
)

func newHighlightCmd() *cobra.Command {
	highlightCmd := &cobra.Command{
		Use:     "highlight",
		Aliases: []string{"hl"},
		Short:   "Manage highlights",
	}

	highlightCmd.AddCommand(newHighlightListCmd())
	highlightCmd.AddCommand(newHighlightGetCmd())
	highlightCmd.AddCommand(newHighlightCreateCmd())
	highlightCmd.AddCommand(newHighlightUpdateCmd())
	highlightCmd.AddCommand(newHighlightDeleteCmd())
	highlightCmd.AddCommand(newHighlightTagCmd())

	return highlightCmd
}

func newHighlightListCmd() *cobra.Command {
	var params client.HighlightListParams

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List highlights",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				return err
			}
			result, err := c.HighlightList(params)
			if err != nil {
				return err
			}
			return getFormatter().Print(result)
		},
	}

	cmd.Flags().IntVar(&params.PageSize, "page-size", 0, "Number of results per page")
	cmd.Flags().IntVar(&params.Page, "page", 0, "Page number")
	cmd.Flags().IntVar(&params.BookID, "book-id", 0, "Filter by book ID")
	cmd.Flags().StringVar(&params.UpdatedLT, "updated-before", "", "Filter by updated before (ISO 8601)")
	cmd.Flags().StringVar(&params.UpdatedGT, "updated-after", "", "Filter by updated after (ISO 8601)")
	cmd.Flags().StringVar(&params.HighlightedAtLT, "highlighted-before", "", "Filter by highlighted before (ISO 8601)")
	cmd.Flags().StringVar(&params.HighlightedAtGT, "highlighted-after", "", "Filter by highlighted after (ISO 8601)")

	return cmd
}

func newHighlightGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a highlight by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid highlight ID: %s", args[0])
			}
			c, err := getClient()
			if err != nil {
				return err
			}
			result, err := c.HighlightGet(id)
			if err != nil {
				return err
			}
			return getFormatter().Print(result)
		},
	}
}

func newHighlightCreateCmd() *cobra.Command {
	var (
		text          string
		title         string
		author        string
		imageURL      string
		sourceURL     string
		sourceType    string
		category      string
		note          string
		location      int
		locationType  string
		highlightedAt string
		highlightURL  string
		fromStdin     bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create highlights",
		Long:  "Create highlights from flags or from JSON on stdin (--stdin).",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				return err
			}

			var req model.HighlightCreateRequest

			if fromStdin {
				if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
					return fmt.Errorf("reading stdin: %w", err)
				}
			} else {
				if text == "" {
					return fmt.Errorf("--text is required (or use --stdin)")
				}
				item := model.HighlightCreateItem{
					Text:          text,
					Title:         title,
					Author:        author,
					ImageURL:      imageURL,
					SourceURL:     sourceURL,
					SourceType:    sourceType,
					Category:      category,
					Note:          note,
					LocationType:  locationType,
					HighlightedAt: highlightedAt,
					HighlightURL:  highlightURL,
				}
				if cmd.Flags().Changed("location") {
					item.Location = &location
				}
				req = model.HighlightCreateRequest{
					Highlights: []model.HighlightCreateItem{item},
				}
			}

			result, err := c.HighlightCreate(req)
			if err != nil {
				return err
			}
			return getFormatter().Print(result)
		},
	}

	cmd.Flags().StringVar(&text, "text", "", "Highlight text")
	cmd.Flags().StringVar(&title, "title", "", "Source title")
	cmd.Flags().StringVar(&author, "author", "", "Source author")
	cmd.Flags().StringVar(&imageURL, "image-url", "", "Cover image URL")
	cmd.Flags().StringVar(&sourceURL, "source-url", "", "Source URL")
	cmd.Flags().StringVar(&sourceType, "source-type", "", "Source type")
	cmd.Flags().StringVar(&category, "category", "", "Category (books|articles|tweets|podcasts)")
	cmd.Flags().StringVar(&note, "note", "", "Note for the highlight")
	cmd.Flags().IntVar(&location, "location", 0, "Position in source")
	cmd.Flags().StringVar(&locationType, "location-type", "", "Location type (page|location|none|order|offset|time_offset)")
	cmd.Flags().StringVar(&highlightedAt, "highlighted-at", "", "Highlight time (ISO 8601)")
	cmd.Flags().StringVar(&highlightURL, "highlight-url", "", "Unique highlight URL")
	cmd.Flags().BoolVar(&fromStdin, "stdin", false, "Read highlight JSON from stdin")

	return cmd
}

func newHighlightUpdateCmd() *cobra.Command {
	var (
		text     string
		note     string
		location string
		url      string
		color    string
	)

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a highlight",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid highlight ID: %s", args[0])
			}
			c, err := getClient()
			if err != nil {
				return err
			}

			req := model.HighlightUpdateRequest{}
			if cmd.Flags().Changed("text") {
				req.Text = &text
			}
			if cmd.Flags().Changed("note") {
				req.Note = &note
			}
			if cmd.Flags().Changed("location") {
				req.Location = &location
			}
			if cmd.Flags().Changed("url") {
				req.URL = &url
			}
			if cmd.Flags().Changed("color") {
				req.Color = &color
			}

			result, err := c.HighlightUpdate(id, req)
			if err != nil {
				return err
			}
			return getFormatter().Print(result)
		},
	}

	cmd.Flags().StringVar(&text, "text", "", "New highlight text")
	cmd.Flags().StringVar(&note, "note", "", "New note")
	cmd.Flags().StringVar(&location, "location", "", "New location")
	cmd.Flags().StringVar(&url, "url", "", "New URL")
	cmd.Flags().StringVar(&color, "color", "", "New color (yellow|blue|pink|orange|green|purple)")

	return cmd
}

func newHighlightDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a highlight",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid highlight ID: %s", args[0])
			}
			c, err := getClient()
			if err != nil {
				return err
			}
			if err := c.HighlightDelete(id); err != nil {
				return err
			}
			getFormatter().PrintMessage("Highlight deleted")
			return nil
		},
	}
}

func newHighlightTagCmd() *cobra.Command {
	tagCmd := &cobra.Command{
		Use:   "tag",
		Short: "Manage highlight tags",
	}

	tagCmd.AddCommand(newHighlightTagListCmd())
	tagCmd.AddCommand(newHighlightTagAddCmd())
	tagCmd.AddCommand(newHighlightTagUpdateCmd())
	tagCmd.AddCommand(newHighlightTagDeleteCmd())

	return tagCmd
}

func newHighlightTagListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list <highlight-id>",
		Short: "List tags for a highlight",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid highlight ID: %s", args[0])
			}
			c, err := getClient()
			if err != nil {
				return err
			}
			result, err := c.HighlightTagList(id)
			if err != nil {
				return err
			}
			return getFormatter().Print(result)
		},
	}
}

func newHighlightTagAddCmd() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "add <highlight-id>",
		Short: "Add a tag to a highlight",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid highlight ID: %s", args[0])
			}
			c, err := getClient()
			if err != nil {
				return err
			}
			result, err := c.HighlightTagAdd(id, name)
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

func newHighlightTagUpdateCmd() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "update <highlight-id> <tag-id>",
		Short: "Update a tag on a highlight",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			highlightID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid highlight ID: %s", args[0])
			}
			tagID, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid tag ID: %s", args[1])
			}
			c, err := getClient()
			if err != nil {
				return err
			}
			result, err := c.HighlightTagUpdate(highlightID, tagID, name)
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

func newHighlightTagDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <highlight-id> <tag-id>",
		Short: "Delete a tag from a highlight",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			highlightID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid highlight ID: %s", args[0])
			}
			tagID, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid tag ID: %s", args[1])
			}
			c, err := getClient()
			if err != nil {
				return err
			}
			if err := c.HighlightTagDelete(highlightID, tagID); err != nil {
				return err
			}
			getFormatter().PrintMessage("Tag deleted")
			return nil
		},
	}
}
