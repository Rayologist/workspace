package add

import (
	"fmt"
	"os"

	"workspace/internal/cli"
	"workspace/internal/config"
	"workspace/internal/layout"

	"github.com/spf13/cobra"
)

type AddOptions struct {
	Config func() (*config.Config, error)
	Layout func() *layout.Layout

	TargetWorkspace string
}

func New(r *cli.Runtime) *cobra.Command {
	opts := &AddOptions{
		Config: r.Config,
		Layout: r.Layout,
	}

	cmd := &cobra.Command{
		Use:   "add <workspace>",
		Short: "Add a new workspace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.TargetWorkspace = args[0]
			return runAdd(opts)
		},
	}

	return cmd
}

func runAdd(opts *AddOptions) error {
	c, err := opts.Config()
	if err != nil {
		return err
	}

	if err := c.AddWorkspace(opts.TargetWorkspace); err != nil {
		return err
	}

	if err := os.MkdirAll(opts.Layout().WorkspaceDir(opts.TargetWorkspace), 0o755); err != nil {
		return err
	}

	if err := c.Save(); err != nil {
		return err
	}

	fmt.Printf("Workspace '%s' created successfully.\n", opts.TargetWorkspace)

	return nil
}
