package list

import (
	"fmt"
	"slices"
	"text/tabwriter"

	"workspace/internal/cli"
	"workspace/internal/config"

	"github.com/spf13/cobra"
)

type ListOptions struct {
	Config    func() (*config.Config, error)
	IOStreams cli.IOStreams
}

func New(r *cli.Runtime) *cobra.Command {
	opts := &ListOptions{
		Config:    r.Config,
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

	fmt.Fprintf(w, "WORKSPACE\tREPO\tBRANCH\n")

	workspaces := make([]string, 0, len(c.Workspaces))
	for workspace := range c.Workspaces {
		workspaces = append(workspaces, workspace)
	}

	slices.Sort(workspaces)

	for _, feat := range workspaces {

		r := c.Workspaces[feat].Repos

		aliases := make([]string, 0, len(r))
		for alias := range r {
			aliases = append(aliases, alias)
		}
		slices.Sort(aliases)

		isFirst := true

		for _, alias := range aliases {
			repo := r[alias]
			branch := repo.Branch
			if isFirst {
				fmt.Fprintf(w, "%s\t%s\t%s\n", feat, alias, branch)
				isFirst = false
			} else {
				fmt.Fprintf(w, "\t%s\t%s\n", alias, branch)
			}
		}

	}

	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
}
