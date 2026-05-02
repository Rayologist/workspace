package doctor

import (
	"fmt"
	"sync"
	"text/tabwriter"

	"workspace/internal/cli"
	"workspace/internal/config"
	"workspace/internal/git"

	"github.com/spf13/cobra"
)

type DoctorOptions struct {
	Config    func() (*config.Config, error)
	IOStreams cli.IOStreams

	ShouldFix bool
}

func New(r *cli.Runtime) *cobra.Command {
	opts := &DoctorOptions{
		Config:    r.Config,
		IOStreams: r.IOStreams,
	}

	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check the health of the workspace and diagnose common issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDoctor(opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.ShouldFix, "fix", "f", false, "Automatically fix detected issues")

	return cmd
}

type Message struct {
	Alias string
	Error []string
}

func runDoctor(opts *DoctorOptions) error {
	c, err := opts.Config()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	size := len(c.Sources)
	result := make([]Message, size)

	keys := make([]string, 0, size)
	for k := range c.Sources {
		keys = append(keys, k)
	}

	for i, key := range keys {
		wg.Go(
			func() {
				source := c.Sources[key]
				if err := git.ValidateRepo(source.Path); err != nil {
					result[i] = Message{
						Alias: key,
						Error: []string{err.Error()},
					}
					return
				}

				if err := git.ValidateBranch(source.Path, source.Branch); err != nil {
					result[i] = Message{
						Alias: key,
						Error: []string{err.Error()},
					}
					return
				}
			})
	}

	wg.Wait()

	maps := make(map[string][]string)
	for _, msg := range result {
		if len(msg.Error) > 0 {
			maps[msg.Alias] = msg.Error
		}
	}

	if len(maps) == 0 {
		fmt.Fprintln(opts.IOStreams.Out, "No issues found in the workspace.")
		return nil
	}

	w := tabwriter.NewWriter(opts.IOStreams.Out, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "ALIAS\tERROR\n")

	for k, v := range maps {
		for i, e := range v {
			if i == 0 {
				fmt.Fprintf(w, "%s\t%d. %s\n", k, i+1, e)
			} else {
				fmt.Fprintf(w, "\t%d. %s\n", i+1, e)
			}
		}
	}

	if err := w.Flush(); err != nil {
		return err
	}

	return nil
}
