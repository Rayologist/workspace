package root

import (
	"workspace/internal/cli"
	sourceCmd "workspace/internal/cmd/source"
	"workspace/internal/layout"

	addCmd "workspace/internal/cmd/root/add"
	doctorCmd "workspace/internal/cmd/root/doctor"
	initCmd "workspace/internal/cmd/root/init"
	listCmd "workspace/internal/cmd/root/list"

	"github.com/spf13/cobra"
)

type RootOptions struct {
	ConfigPath string
}

func New(r *cli.Runtime) *cobra.Command {
	opts := &RootOptions{}
	cmd := &cobra.Command{
		Use:   "ws",
		Short: "Workspace manager — multi-repo worktree orchestration",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var o []cli.Option
			if opts.ConfigPath != "" {
				o = append(
					o,
					cli.WithLayout(
						layout.WithConfigPath(opts.ConfigPath),
					),
				)
			}

			return r.Init(o...)
		},
	}

	cmd.PersistentFlags().StringVarP(&opts.ConfigPath, "config", "c", "", "Path to workspace configuration file")

	cmd.AddCommand(initCmd.New(r))
	cmd.AddCommand(doctorCmd.New(r))
	cmd.AddCommand(addCmd.New(r))
	cmd.AddCommand(listCmd.New(r))

	// subcommands
	cmd.AddCommand(sourceCmd.New(r))

	cmd.SilenceUsage = true

	return cmd
}
