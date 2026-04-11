package update

import (
	"fmt"

	"workspace/internal/cli"
	"workspace/internal/cmd/repo/shared"
	"workspace/internal/config"

	"github.com/spf13/cobra"
)

type UpdateOptions struct {
	Config func() (*config.Config, error)

	Alias    string
	NewAlias string
	config.RepoConfig
}

func New(r *cli.Runtime) *cobra.Command {
	opts := &UpdateOptions{
		Config: r.Config,
	}

	cmd := &cobra.Command{
		Use:   "update <alias>",
		Args:  cobra.ExactArgs(1),
		Short: "Update a registered source repository in the workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Alias = args[0]
			return runUpdate(opts, cmd)
		},
	}

	cmd.Flags().StringVarP(&opts.NewAlias, "alias", "a", "", "New alias for the repository")
	cmd.Flags().StringVarP(&opts.Path, "path", "p", "", "Path to the repository")
	cmd.Flags().StringVarP(&opts.Branch, "branch", "b", "", "Default branch to use for this repository")
	cmd.Flags().StringArrayVarP(&opts.Hooks.Setup, "setup", "s", nil, "Commands to run after cloning the repository")

	return cmd
}

func runUpdate(opts *UpdateOptions, cmd *cobra.Command) error {
	c, err := opts.Config()
	if err != nil {
		return err
	}

	builder := shared.NewUpdateRepoBuilder(c, opts.Alias)

	if cmd.Flags().Changed("alias") {
		builder.AliasUpdate(opts.NewAlias)
	}

	if cmd.Flags().Changed("path") {
		builder.Path(opts.Path)
	}

	if cmd.Flags().Changed("branch") {
		builder.Branch(opts.Branch)
	}

	if cmd.Flags().Changed("setup") {
		builder.SetupHookAppend(opts.Hooks.Setup)
	}

	err = builder.Commit()
	if err != nil {
		return err
	}

	if err := c.Save(); err != nil {
		return err
	}

	if cmd.Flags().Changed("alias") {
		fmt.Printf("Repository '%s' updated successfully (new alias: '%s').\n", opts.Alias, opts.NewAlias)
	} else {
		fmt.Printf("Repository '%s' updated successfully.\n", opts.Alias)
	}

	return nil
}
