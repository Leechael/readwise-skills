package cmd

import (
	"github.com/spf13/cobra"
)

func newReviewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "review",
		Short: "Get daily review highlights",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				return err
			}
			result, err := c.DailyReview()
			if err != nil {
				return err
			}
			return getFormatter().Print(result)
		},
	}
}
