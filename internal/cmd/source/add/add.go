package add

import (
	"fmt"

	"workspace/internal/cli"
	"workspace/internal/cmd/source/shared"
	"workspace/internal/config"

	"github.com/spf13/cobra"
)

type AddOptions struct {
	Config func() (*config.Config, error)

	Alias string
	config.SourceConfig
}

func New(r *cli.Runtime) *cobra.Command {
	opts := &AddOptions{
		Config: r.Config,
	}

	cmd := &cobra.Command{
		Use:   "add <alias> [path]",
		Args:  cobra.ExactArgs(2),
		Short: "Register a source repository in the workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Alias = args[0]
			opts.Path = args[1]

			return runAdd(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Branch, "branch", "b", "main", "Default branch to use for this repository")
	cmd.Flags().StringArrayVarP(&opts.Hooks.Setup, "setup", "s", []string{}, "Commands to run after cloning the repository")

	return cmd
}

func runAdd(opts *AddOptions) error {
	c, err := opts.Config()
	if err != nil {
		return err
	}

	builder := shared.NewAddSourceBuilder(c, opts.Alias).
		Path(opts.Path).
		Branch(opts.Branch).
		SetupHookAppend(opts.Hooks.Setup)
	if err := builder.Commit(); err != nil {
		return err
	}

	if err := c.Save(); err != nil {
		return err
	}

	fmt.Printf("Repository '%s' added successfully.\n", opts.Alias)

	return nil
}
