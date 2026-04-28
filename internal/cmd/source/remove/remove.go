package remove

import (
	"fmt"

	"workspace/internal/cli"
	"workspace/internal/config"

	"github.com/spf13/cobra"
)

type RemoveOptions struct {
	Config func() (*config.Config, error)

	Alias string
}

func New(r *cli.Runtime) *cobra.Command {
	opts := &RemoveOptions{
		Config: r.Config,
	}

	cmd := &cobra.Command{
		Use:   "remove <alias>",
		Args:  cobra.ExactArgs(1),
		Short: "Remove a registered source repository from the workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Alias = args[0]
			return runRemove(opts)
		},
	}

	return cmd
}

func runRemove(opts *RemoveOptions) error {
	c, err := opts.Config()
	if err != nil {
		return err
	}

	err = c.RemoveSource(opts.Alias)
	if err != nil {
		return err
	}

	if err := c.Save(); err != nil {
		return err
	}

	fmt.Printf("Repository '%s' removed successfully.\n", opts.Alias)

	return nil
}
