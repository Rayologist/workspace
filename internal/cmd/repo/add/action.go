package add

import (
	"fmt"

	"workspace/internal/config"
	"workspace/internal/git"
)

func add(opts *AddOptions) error {
	c, err := config.Load()
	if err != nil {
		return err
	}

	if err := git.ValidateRepo(opts.Path); err != nil {
		return err
	}

	if err := git.ValidateBranch(opts.Path, opts.Branch); err != nil {
		return err
	}

	if _, exists := c.Repos[opts.Alias]; exists {
		return fmt.Errorf("repository '%s' already exists (use 'repo update' to modify it)", opts.Alias)
	}

	repo := &config.RepoConfig{}

	repo.Path = opts.Path
	repo.Branch = opts.Branch
	repo.Hooks.Setup = opts.Hooks.Setup

	c.Repos[opts.Alias] = repo

	err = c.Save()
	if err != nil {
		return err
	}

	fmt.Printf("Repository '%s' added successfully.\n", opts.Alias)

	return nil
}
