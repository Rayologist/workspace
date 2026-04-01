package list

import "github.com/spf13/cobra"

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List projects in the workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			return list()
		},
	}

	return cmd
}
