package config

import "fmt"

type RepoConfig struct {
	Path   string      `yaml:"path"`
	Branch string      `yaml:"branch,omitempty"`
	Hooks  HooksConfig `yaml:"hooks,omitempty"`
}

type RepoConfigs map[string]*RepoConfig

func NewRepoConfigs() RepoConfigs {
	return make(RepoConfigs)
}

func (c *Config) RepoByAlias(alias string) (*RepoConfig, error) {
	repo, exists := c.Repos[alias]
	if !exists {
		return nil, fmt.Errorf("repo '%s' not exist in the config", alias)
	}
	return repo, nil
}

func (c *Config) RemoveRepo(alias string) error {
	if _, err := c.RepoByAlias(alias); err != nil {
		return err
	}

	delete(c.Repos, alias)

	return nil
}

func (c *Config) AddRepo(alias string, repo *RepoConfig) error {
	if _, err := c.RepoByAlias(alias); err == nil {
		return fmt.Errorf("repository '%s' already exists (use 'repo update' to modify it)", alias)
	}

	c.Repos[alias] = repo
	return nil
}

func (c *Config) UpdateRepo(alias, newAlias string, repo *RepoConfig) error {
	if _, err := c.RepoByAlias(alias); err != nil {
		return err
	}

	targetAlias := alias
	if newAlias != "" {
		if _, err := c.RepoByAlias(newAlias); err == nil {
			return fmt.Errorf("repository alias '%s' already exists", newAlias)
		}

		delete(c.Repos, alias)
		targetAlias = newAlias
	}

	c.Repos[targetAlias] = repo
	return nil
}
