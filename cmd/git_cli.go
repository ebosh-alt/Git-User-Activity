package cmd

import (
	"github.com/spf13/cobra"
)

func newGitCLICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "git-cli",
		Short: "GitHub утилиты",
	}

	cmd.AddCommand(newActivityCmd())

	return cmd
}
