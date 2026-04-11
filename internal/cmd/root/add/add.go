package add

import (
	"workspace/internal/cli"
	"workspace/internal/config"

	"github.com/spf13/cobra"
)

type AddOptions struct {
	Config func() (*config.Config, error)

	ProjectName string
	RepoConfigs []string
}

func New(r *cli.Runtime) *cobra.Command {
	opts := &AddOptions{
		Config: r.Config,
	}

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new project to the workspace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProjectName = args[0]
			return runAdd(opts)
		},
	}

	cmd.Flags().StringArrayVarP(&opts.RepoConfigs, "repo", "r", []string{}, "Add a repository to the project (branch defaults to 'main', if not specified) - format: <alias>[:<branch>]")
	cmd.MarkFlagRequired("repo")

	return cmd
}
