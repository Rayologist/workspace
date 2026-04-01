package repo

import (
	addCmd "workspace/internal/cmd/repo/add"

	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repo",
		Short: "Manage source repository registrations",
	}

	cmd.AddCommand(addCmd.New())

	return cmd
}
