package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Check readwise-cli credential and API connectivity status",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				return fmt.Errorf("Readwise credentials are not configured. Set READWISE_API_TOKEN, then try again.\n%s", err)
			}

			f := getFormatter()

			if err := c.AuthCheck(); err != nil {
				f.Hint("Readwise API check failed. Verify token and network connectivity.")
				return err
			}

			f.PrintMessage("OK: readwise-cli is configured and Readwise API is reachable.")
			return nil
		},
	}
}
