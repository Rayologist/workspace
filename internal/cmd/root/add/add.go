package add

import (
	"github.com/spf13/cobra"
)

type AddOptions struct {
	ProjectName string
	RepoConfigs []string
}

func New() *cobra.Command {
	opts := &AddOptions{}

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new project to the workspace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProjectName = args[0]
			return add(opts)
		},
	}

	cmd.Flags().StringArrayVarP(&opts.RepoConfigs, "repo", "r", []string{}, "Add a repository to the project (branch defaults to 'main', if not specified) - format: <alias>[:<branch>]")
	cmd.MarkFlagRequired("repo")

	return cmd
}
