package add

import (
	"fmt"

	"workspace/internal/config"
	"workspace/internal/git"
	"workspace/internal/set"
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

	// add command should be idempotent, so if the alias already exists, we should update it instead of returning an error
	if repo, exists := c.Repos[opts.Alias]; exists {

		if repo.Path != opts.Path {
			repo.Path = opts.Path
		}

		if repo.Branch != opts.Branch {
			repo.Branch = opts.Branch
		}

		if opts.Hooks.Setup != nil {
			configSet := set.FromSlice(repo.Hooks.Setup)
			for _, hook := range opts.Hooks.Setup {
				if !configSet.Contains(hook) {
					repo.Hooks.Setup = append(repo.Hooks.Setup, hook)
				}
			}
		}

		err := c.Save()
		if err != nil {
			return err
		}

		fmt.Printf("Repository '%s' updated successfully.\n", opts.Alias)

		return nil
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
