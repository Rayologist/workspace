package repo

import (
	addCmd "workspace/internal/cmd/repo/add"
	removeCmd "workspace/internal/cmd/repo/remove"
	updateCmd "workspace/internal/cmd/repo/update"

	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repo",
		Short: "Manage source repository registrations",
	}

	cmd.AddCommand(addCmd.New())
	cmd.AddCommand(updateCmd.New())
	cmd.AddCommand(removeCmd.New())

	return cmd
}
