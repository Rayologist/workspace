package update

import (
	"workspace/internal/config"

	"github.com/spf13/cobra"
)

type UpdateOptions struct {
	Alias string
	config.RepoConfig
}

func New() *cobra.Command {
	opts := &UpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <alias>",
		Args:  cobra.ExactArgs(1),
		Short: "Update a registered source repository in the workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Alias = args[0]
			return update(opts, cmd)
		},
	}

	cmd.Flags().StringVarP(&opts.Path, "path", "p", "", "Path to the repository")
	cmd.Flags().StringVarP(&opts.Branch, "branch", "b", "", "Default branch to use for this repository")
	cmd.Flags().StringArrayVarP(&opts.Hooks.Setup, "setup", "s", nil, "Commands to run after cloning the repository")

	return cmd
}
