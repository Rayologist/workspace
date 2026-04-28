package config

import "fmt"

type WorkspaceRepoConfig struct {
	Branch string `yaml:"branch"`
}

type WorkspaceRepoConfigs map[string]*WorkspaceRepoConfig

type WorkspaceConfig struct {
	Repos WorkspaceRepoConfigs `yaml:"repos"`
}

type WorkspaceConfigs map[string]*WorkspaceConfig

func NewWorkspaceConfigs() WorkspaceConfigs {
	return make(WorkspaceConfigs)
}

func (c *Config) Workspace(name string) (*WorkspaceConfig, error) {
	w, exists := c.Workspaces[name]
	if !exists {
		return nil, fmt.Errorf("workspace '%s' not exist in the config", name)
	}
	return w, nil
}

func (c *Config) AddWorkspaceRepo(workspaceName, sourceAlias string, config *WorkspaceRepoConfig) error {
	if _, err := c.SourceByAlias(sourceAlias); err != nil {
		return err
	}

	w, err := c.Workspace(workspaceName)

	if err == nil {
		w.Repos[sourceAlias] = config
		return nil
	}

	c.Workspaces[workspaceName] = &WorkspaceConfig{
		Repos: WorkspaceRepoConfigs{
			sourceAlias: config,
		},
	}

	return nil
}

func (c *Config) UpdateWorkspaceRepoBranch(workspaceName, sourceAlias, repoBranch string) error {
	w, err := c.Workspace(workspaceName)
	if err != nil {
		return err
	}

	if _, err := c.SourceByAlias(sourceAlias); err != nil {
		return err
	}

	w.Repos[sourceAlias] = &WorkspaceRepoConfig{
		Branch: repoBranch,
	}
	return nil
}
