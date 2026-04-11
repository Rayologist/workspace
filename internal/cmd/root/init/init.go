package init

import (
	"fmt"
	"os"

	"workspace/internal/cli"
	"workspace/internal/config"
	"workspace/internal/layout"

	"github.com/spf13/cobra"
)

type InitOptions struct {
	Layout    func() *layout.Layout
	IOStreams cli.IOStreams
}

func New(r *cli.Runtime) *cobra.Command {
	opts := &InitOptions{
		Layout:    r.Layout,
		IOStreams: r.IOStreams,
	}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new workspace in the current directory",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit(opts)
		},
	}

	return cmd
}

func runInit(opts *InitOptions) error {
	path := opts.Layout().ConfigPath()
	f, err := os.Stat(path)

	if f != nil {
		fmt.Fprintf(opts.IOStreams.Out, "Workspace already initialized.\n")
		return nil
	}

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	c := config.New(path)
	err = c.Save()
	if err != nil {
		return err
	}

	fmt.Printf("Workspace initialized successfully.\n")

	return nil
}
