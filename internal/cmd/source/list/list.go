package list

import (
	"fmt"
	"path/filepath"
	"text/tabwriter"

	"workspace/internal/cli"
	"workspace/internal/config"
	"workspace/internal/layout"

	"github.com/spf13/cobra"
)

type ListOptions struct {
	Config    func() (*config.Config, error)
	Layout    func() *layout.Layout
	IOStreams cli.IOStreams
}

func New(r *cli.Runtime) *cobra.Command {
	opts := &ListOptions{
		Config:    r.Config,
		Layout:    r.Layout,
		IOStreams: r.IOStreams,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List projects in the workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(opts)
		},
	}

	return cmd
}

func runList(opts *ListOptions) error {
	c, err := opts.Config()
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(opts.IOStreams.Out, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "ALIAS\tREPO\tBRANCH\tSETUPS\n")

	layout := opts.Layout()
	for k, v := range c.Sources {
		relPath, err := filepath.Rel(layout.Root(), v.Path)
		if err != nil {
			relPath = v.Path
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%v\n", k, relPath, v.Branch, len(v.Hooks.Setup))
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
}
