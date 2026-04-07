package shared

import (
	"fmt"
	"path/filepath"

	"workspace/internal/config"
	"workspace/internal/git"
)

type RepoBuilder struct {
	config   *config.Config
	repo     *config.RepoConfig
	alias    string
	newAlias string
	isUpdate bool
	err      error
}

func NewAddRepoBuilder(c *config.Config, alias string) *RepoBuilder {
	b := &RepoBuilder{
		config: c,
		alias:  alias,
		repo:   &config.RepoConfig{},
	}

	return b
}

func NewUpdateRepoBuilder(c *config.Config, alias string) *RepoBuilder {
	b := &RepoBuilder{
		config:   c,
		alias:    alias,
		isUpdate: true,
	}

	repo, err := c.RepoByAlias(alias)
	if err != nil {
		b.err = err
		return b
	}

	copy := *repo
	b.repo = &copy

	return b
}

func (b *RepoBuilder) Path(path string) *RepoBuilder {
	if b.err != nil {
		return b
	}

	if err := git.ValidateRepo(path); err != nil {
		b.err = err
		return b
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		b.err = fmt.Errorf("failed to get absolute path: %w", err)
		return b
	}

	b.repo.Path = absPath
	return b
}

func (b *RepoBuilder) Branch(branch string) *RepoBuilder {
	if b.err != nil {
		return b
	}

	if b.repo.Path == "" {
		b.err = fmt.Errorf("repository path must be set before validating branch")
		return b
	}

	if err := git.ValidateBranch(b.repo.Path, branch); err != nil {
		b.err = err
		return b
	}

	b.repo.Branch = branch
	return b
}

func (b *RepoBuilder) SetupHookAppend(setups []string) *RepoBuilder {
	if b.err != nil {
		return b
	}

	b.repo.Hooks.AppendSetupHooks(setups)
	return b
}

func (b *RepoBuilder) AliasUpdate(newAlias string) *RepoBuilder {
	if b.err != nil {
		return b
	}

	b.newAlias = newAlias
	return b
}

func (b *RepoBuilder) Commit() error {
	if b.err != nil {
		return b.err
	}

	if b.isUpdate {
		if err := git.ValidateBranch(b.repo.Path, b.repo.Branch); err != nil {
			return err
		}

		return b.config.UpdateRepo(b.alias, b.newAlias, b.repo)
	}

	return b.config.AddRepo(b.alias, b.repo)
}
