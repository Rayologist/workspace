package update

import (
	"fmt"

	"workspace/internal/cmd/repo/shared"
	"workspace/internal/config"

	"github.com/spf13/cobra"
)

func update(opts *UpdateOptions, cmd *cobra.Command) error {
	c, err := config.Load()
	if err != nil {
		return err
	}

	builder := shared.NewUpdateRepoBuilder(c, opts.Alias)

	if cmd.Flags().Changed("alias") {
		builder.AliasUpdate(opts.NewAlias)
	}

	if cmd.Flags().Changed("path") {
		builder.Path(opts.Path)
	}

	if cmd.Flags().Changed("branch") {
		builder.Branch(opts.Branch)
	}

	if cmd.Flags().Changed("setup") {
		builder.SetupHookAppend(opts.Hooks.Setup)
	}

	err = builder.Commit()
	if err != nil {
		return err
	}

	if err := c.Save(); err != nil {
		return err
	}

	if cmd.Flags().Changed("alias") {
		fmt.Printf("Repository '%s' updated successfully (new alias: '%s').\n", opts.Alias, opts.NewAlias)
	} else {
		fmt.Printf("Repository '%s' updated successfully.\n", opts.Alias)
	}

	return nil
}
