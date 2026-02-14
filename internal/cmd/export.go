package cmd

import (
	"github.com/spf13/cobra"

	"github.com/leechael/readwise-skills/internal/client"
)

func newExportCmd() *cobra.Command {
	var params client.ExportParams

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export highlights with full book metadata",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				return err
			}
			result, err := c.Export(params)
			if err != nil {
				return err
			}
			return getFormatter().Print(result)
		},
	}

	cmd.Flags().StringVar(&params.UpdatedAfter, "updated-after", "", "Only return items updated after this date (ISO 8601)")
	cmd.Flags().StringVar(&params.IDs, "ids", "", "Comma-separated list of book IDs")
	cmd.Flags().BoolVar(&params.IncludeDeleted, "include-deleted", false, "Include deleted highlights")
	cmd.Flags().StringVar(&params.PageCursor, "cursor", "", "Pagination cursor")

	return cmd
}
