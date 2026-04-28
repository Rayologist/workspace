package shared

import (
	"fmt"
	"path/filepath"

	"workspace/internal/config"
	"workspace/internal/git"
)

type SourceBuilder struct {
	config   *config.Config
	source   *config.SourceConfig
	alias    string
	newAlias string
	isUpdate bool
	err      error
}

func NewAddSourceBuilder(c *config.Config, alias string) *SourceBuilder {
	b := &SourceBuilder{
		config: c,
		alias:  alias,
		source: &config.SourceConfig{},
	}

	return b
}

func NewUpdateSourceBuilder(c *config.Config, alias string) *SourceBuilder {
	b := &SourceBuilder{
		config:   c,
		alias:    alias,
		isUpdate: true,
	}

	source, err := c.SourceByAlias(alias)
	if err != nil {
		b.err = err
		return b
	}

	copy := *source
	b.source = &copy

	return b
}

func (b *SourceBuilder) Path(path string) *SourceBuilder {
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

	b.source.Path = absPath
	return b
}

func (b *SourceBuilder) Branch(branch string) *SourceBuilder {
	if b.err != nil {
		return b
	}

	if b.source.Path == "" {
		b.err = fmt.Errorf("source path must be set before validating branch")
		return b
	}

	if err := git.ValidateBranch(b.source.Path, branch); err != nil {
		b.err = err
		return b
	}

	b.source.Branch = branch
	return b
}

func (b *SourceBuilder) SetupHookAppend(setups []string) *SourceBuilder {
	if b.err != nil {
		return b
	}

	b.source.Hooks.AppendSetupHooks(setups)
	return b
}

func (b *SourceBuilder) AliasUpdate(newAlias string) *SourceBuilder {
	if b.err != nil {
		return b
	}

	b.newAlias = newAlias
	return b
}

func (b *SourceBuilder) Commit() error {
	if b.err != nil {
		return b.err
	}

	if b.isUpdate {
		if err := git.ValidateBranch(b.source.Path, b.source.Branch); err != nil {
			return err
		}

		return b.config.UpdateSource(b.alias, b.newAlias, b.source)
	}

	return b.config.AddSource(b.alias, b.source)
}
