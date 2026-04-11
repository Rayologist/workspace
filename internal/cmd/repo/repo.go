package repo

import (
	"workspace/internal/cli"
	addCmd "workspace/internal/cmd/repo/add"
	listCmd "workspace/internal/cmd/repo/list"
	removeCmd "workspace/internal/cmd/repo/remove"

	updateCmd "workspace/internal/cmd/repo/update"

	"github.com/spf13/cobra"
)

func New(r *cli.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repo",
		Short: "Manage source repository registrations",
	}

	cmd.AddCommand(listCmd.New(r))
	cmd.AddCommand(addCmd.New(r))
	cmd.AddCommand(updateCmd.New(r))
	cmd.AddCommand(removeCmd.New(r))

	return cmd
}
