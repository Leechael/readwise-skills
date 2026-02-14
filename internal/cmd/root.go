package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/leechael/readwise-skills/internal/client"
	"github.com/leechael/readwise-skills/internal/output"
)

const (
	ExitSuccess  = 0
	ExitError    = 1
	ExitAuth     = 2
	ExitNotFound = 3
)

var (
	token     string
	jsonMode  bool
	plainMode bool
	jqFilter  string
)

// NewRootCmd creates the root cobra command.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "readwise",
		Short: "CLI for Readwise v2 and Reader v3 APIs",
		Long:  "A unified CLI tool for interacting with both Readwise (highlights/books) and Reader (documents) APIs.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if jqFilter != "" && !jsonMode {
				return fmt.Errorf("--jq requires --json")
			}
			return nil
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.PersistentFlags().StringVar(&token, "token", "", "Readwise API token (or set READWISE_TOKEN)")
	rootCmd.PersistentFlags().BoolVar(&jsonMode, "json", false, "Output as JSON")
	rootCmd.PersistentFlags().BoolVar(&plainMode, "plain", false, "Output as plain tab-separated text")
	rootCmd.PersistentFlags().StringVar(&jqFilter, "jq", "", "JQ filter (requires --json)")

	rootCmd.AddCommand(newAuthCmd())
	rootCmd.AddCommand(newHighlightCmd())
	rootCmd.AddCommand(newBookCmd())
	rootCmd.AddCommand(newExportCmd())
	rootCmd.AddCommand(newReviewCmd())
	rootCmd.AddCommand(newReaderCmd())

	return rootCmd
}

func getClient() (*client.Client, error) {
	t := token
	if t == "" {
		t = os.Getenv("READWISE_TOKEN")
	}
	if t == "" {
		return nil, fmt.Errorf("no API token provided; set READWISE_TOKEN or use --token")
	}
	return client.New(t), nil
}

func getFormatter() *output.Formatter {
	return output.NewFormatter(jsonMode, plainMode, jqFilter)
}

// ExitCode returns an appropriate exit code for the given error.
func ExitCode(err error) int {
	if err == nil {
		return ExitSuccess
	}
	if apiErr, ok := err.(*client.APIError); ok {
		switch {
		case apiErr.StatusCode == 401 || apiErr.StatusCode == 403:
			return ExitAuth
		case apiErr.StatusCode == 404:
			return ExitNotFound
		}
	}
	return ExitError
}
