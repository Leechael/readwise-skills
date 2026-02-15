package cmd

import (
	"github.com/spf13/cobra"

	"github.com/leechael/readwise-skills/internal/client"
	"github.com/leechael/readwise-skills/internal/model"
)

func newReaderCmd() *cobra.Command {
	readerCmd := &cobra.Command{
		Use:   "reader",
		Short: "Manage Reader documents",
	}

	readerCmd.AddCommand(newReaderListCmd())
	readerCmd.AddCommand(newReaderSaveCmd())
	readerCmd.AddCommand(newReaderUpdateCmd())
	readerCmd.AddCommand(newReaderDeleteCmd())
	readerCmd.AddCommand(newReaderTagCmd())

	return readerCmd
}

func newReaderListCmd() *cobra.Command {
	var params client.ReaderListParams

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Reader documents",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				return err
			}
			result, err := c.ReaderList(params)
			if err != nil {
				return err
			}
			return getFormatter().Print(result)
		},
	}

	cmd.Flags().StringVar(&params.ID, "id", "", "Filter by document ID")
	cmd.Flags().StringVar(&params.Category, "category", "", "Filter by category")
	cmd.Flags().StringVar(&params.Location, "location", "", "Filter by location (new|later|archive|feed)")
	cmd.Flags().StringVar(&params.UpdatedAfter, "updated-after", "", "Filter by updated after (ISO 8601)")
	cmd.Flags().StringVar(&params.PageCursor, "cursor", "", "Pagination cursor")

	return cmd
}

func newReaderSaveCmd() *cobra.Command {
	var (
		url        string
		html       string
		title      string
		author     string
		summary    string
		imageURL   string
		location   string
		category   string
		savedUsing string
		notes      string
		tags       []string
	)

	cmd := &cobra.Command{
		Use:   "save",
		Short: "Save a document to Reader",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				return err
			}

			req := model.DocumentSaveRequest{
				URL:        url,
				HTML:       html,
				Title:      title,
				Author:     author,
				Summary:    summary,
				ImageURL:   imageURL,
				Location:   location,
				Category:   category,
				SavedUsing: savedUsing,
				Notes:      notes,
				Tags:       tags,
			}

			result, err := c.ReaderSave(req)
			if err != nil {
				return err
			}
			return getFormatter().Print(result)
		},
	}

	cmd.Flags().StringVar(&url, "url", "", "URL to save")
	cmd.MarkFlagRequired("url")
	cmd.Flags().StringVar(&html, "html", "", "HTML content (optional)")
	cmd.Flags().StringVar(&title, "title", "", "Document title")
	cmd.Flags().StringVar(&author, "author", "", "Document author")
	cmd.Flags().StringVar(&summary, "summary", "", "Document summary")
	cmd.Flags().StringVar(&imageURL, "image-url", "", "Image URL")
	cmd.Flags().StringVar(&location, "location", "", "Location (new|later|archive|feed)")
	cmd.Flags().StringVar(&category, "category", "", "Category")
	cmd.Flags().StringVar(&savedUsing, "saved-using", "", "Source of save")
	cmd.Flags().StringVar(&notes, "notes", "", "Notes")
	cmd.Flags().StringSliceVar(&tags, "tag", nil, "Tags (can be specified multiple times)")

	return cmd
}

func newReaderUpdateCmd() *cobra.Command {
	var (
		title    string
		author   string
		summary  string
		notes    string
		location string
		category string
	)

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a Reader document",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				return err
			}

			req := model.DocumentUpdateRequest{}
			if cmd.Flags().Changed("title") {
				req.Title = &title
			}
			if cmd.Flags().Changed("author") {
				req.Author = &author
			}
			if cmd.Flags().Changed("summary") {
				req.Summary = &summary
			}
			if cmd.Flags().Changed("notes") {
				req.Notes = &notes
			}
			if cmd.Flags().Changed("location") {
				req.Location = &location
			}
			if cmd.Flags().Changed("category") {
				req.Category = &category
			}

			if err := c.ReaderUpdate(args[0], req); err != nil {
				return err
			}
			getFormatter().PrintMessage("Document updated")
			return nil
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "New title")
	cmd.Flags().StringVar(&author, "author", "", "New author")
	cmd.Flags().StringVar(&summary, "summary", "", "New summary")
	cmd.Flags().StringVar(&notes, "notes", "", "New notes")
	cmd.Flags().StringVar(&location, "location", "", "New location")
	cmd.Flags().StringVar(&category, "category", "", "New category")

	return cmd
}

func newReaderDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a Reader document",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				return err
			}
			if err := c.ReaderDelete(args[0]); err != nil {
				return err
			}
			getFormatter().PrintMessage("Document deleted")
			return nil
		},
	}
}

func newReaderTagCmd() *cobra.Command {
	tagCmd := &cobra.Command{
		Use:   "tag",
		Short: "Manage Reader tags",
	}

	tagCmd.AddCommand(newReaderTagListCmd())
	return tagCmd
}

func newReaderTagListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all Reader tags",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				return err
			}
			result, err := c.ReaderTagList()
			if err != nil {
				return err
			}
			return getFormatter().Print(result)
		},
	}
}
