package root

import (
	repoCmd "workspace/internal/cmd/repo"
	addCmd "workspace/internal/cmd/root/add"
	doctorCmd "workspace/internal/cmd/root/doctor"
	initCmd "workspace/internal/cmd/root/init"

	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ws",
		Short: "Workspace manager — multi-repo worktree orchestration",
	}

	cmd.AddCommand(initCmd.New())
	cmd.AddCommand(doctorCmd.New())
	cmd.AddCommand(addCmd.New())

	// subcommands
	cmd.AddCommand(repoCmd.New())

	cmd.SilenceUsage = true

	return cmd
}
