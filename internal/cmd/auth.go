package cmd

import (
	"github.com/spf13/cobra"
)

func newAuthCmd() *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Authentication commands",
	}

	authCmd.AddCommand(newAuthCheckCmd())
	return authCmd
}

func newAuthCheckCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Validate API token",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				return err
			}
			f := getFormatter()

			if err := c.AuthCheck(); err != nil {
				f.Hint("Token is invalid")
				return err
			}

			f.PrintMessage("Token is valid")
			return nil
		},
	}
}
