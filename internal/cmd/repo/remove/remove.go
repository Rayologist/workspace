package remove

import (
	"github.com/spf13/cobra"
)

type RemoveOptions struct {
	Alias string
}

func New() *cobra.Command {
	opts := &RemoveOptions{}

	cmd := &cobra.Command{
		Use:   "remove <alias>",
		Args:  cobra.ExactArgs(1),
		Short: "Remove a registered source repository from the workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Alias = args[0]
			return remove(opts)
		},
	}

	return cmd
}
