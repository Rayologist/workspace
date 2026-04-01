package root

import (
	repoCmd "workspace/internal/cmd/repo"
	initCmd "workspace/internal/cmd/root/init"

	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ws",
		Short: "Workspace manager — multi-repo worktree orchestration",
	}

	cmd.AddCommand(initCmd.New())

	// subcommands
	cmd.AddCommand(repoCmd.New())

	return cmd
}
