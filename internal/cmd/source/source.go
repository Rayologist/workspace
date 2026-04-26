package source

import (
	"workspace/internal/cli"
	addCmd "workspace/internal/cmd/source/add"
	listCmd "workspace/internal/cmd/source/list"
	removeCmd "workspace/internal/cmd/source/remove"

	updateCmd "workspace/internal/cmd/source/update"

	"github.com/spf13/cobra"
)

func New(r *cli.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "source",
		Short: "Manage source repository registrations",
	}

	cmd.AddCommand(listCmd.New(r))
	cmd.AddCommand(addCmd.New(r))
	cmd.AddCommand(updateCmd.New(r))
	cmd.AddCommand(removeCmd.New(r))

	return cmd
}
