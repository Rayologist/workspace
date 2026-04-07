package add

import (
	"fmt"

	"workspace/internal/cmd/repo/shared"
	"workspace/internal/config"
)

func add(opts *AddOptions) error {
	c, err := config.Load()
	if err != nil {
		return err
	}

	err = shared.NewAddRepoBuilder(c, opts.Alias).
		Path(opts.Path).
		Branch(opts.Branch).
		SetupHookAppend(opts.Hooks.Setup).
		Commit()
	if err != nil {
		return err
	}

	err = c.Save()
	if err != nil {
		return err
	}

	fmt.Printf("Repository '%s' added successfully.\n", opts.Alias)

	return nil
}
