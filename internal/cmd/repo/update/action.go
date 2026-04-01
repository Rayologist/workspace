package update

import (
	"fmt"

	"workspace/internal/config"
	"workspace/internal/git"
	"workspace/internal/set"

	"github.com/spf13/cobra"
)

func update(opts *UpdateOptions, cmd *cobra.Command) error {
	c, err := config.Load()
	if err != nil {
		return err
	}

	repo, exists := c.Repos[opts.Alias]
	if !exists {
		return fmt.Errorf("repository '%s' not found", opts.Alias)
	}

	if cmd.Flags().Changed("path") {
		if err := git.ValidateRepo(opts.Path); err != nil {
			return err
		}

		repo.Path = opts.Path

		if err := git.ValidateBranch(repo.Path, repo.Branch); err != nil {
			return err
		}
	}

	if cmd.Flags().Changed("branch") {
		if err := git.ValidateBranch(repo.Path, opts.Branch); err != nil {
			return err
		}

		repo.Branch = opts.Branch
	}

	if cmd.Flags().Changed("setup") {
		configSet := set.FromSlice(repo.Hooks.Setup)
		for _, hook := range opts.Hooks.Setup {
			if !configSet.Contains(hook) {
				repo.Hooks.Setup = append(repo.Hooks.Setup, hook)
			}
		}
	}

	aliasUpdated := false

	// alias update should be at the end to avoid conflicts when renaming the repository
	if cmd.Flags().Changed("alias") {
		if _, exists := c.Repos[opts.NewAlias]; exists {
			return fmt.Errorf("repository '%s' already exists", opts.NewAlias)
		}

		delete(c.Repos, opts.Alias)
		c.Repos[opts.NewAlias] = repo

		aliasUpdated = true
	}

	if err := c.Save(); err != nil {
		return err
	}

	if aliasUpdated {
		fmt.Printf("Repository '%s' updated successfully (new alias: '%s').\n", opts.Alias, opts.NewAlias)
	} else {
		fmt.Printf("Repository '%s' updated successfully.\n", opts.Alias)
	}

	return nil
}
