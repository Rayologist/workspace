package add

import (
	"workspace/internal/config"

	"github.com/spf13/cobra"
)

type AddOptions struct {
	Alias string
	config.RepoConfig
}

func New() *cobra.Command {
	opts := &AddOptions{}

	cmd := &cobra.Command{
		Use:   "add <alias> [path]",
		Args:  cobra.ExactArgs(2),
		Short: "Register a source repository in the workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Alias = args[0]
			opts.Path = args[1]

			return add(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Branch, "branch", "b", "main", "Default branch to use for this repository")
	cmd.Flags().StringArrayVarP(&opts.Hooks.Setup, "setup", "s", []string{}, "Commands to run after cloning the repository")

	return cmd
}
