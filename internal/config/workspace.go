package config

type WorkspaceRepoConfigs map[string]*WorkspaceRepoConfig

type WorkspaceConfig struct {
	Repos WorkspaceRepoConfigs `yaml:"repos"`
}

func NewWorkspaceConfigs() WorkspaceConfigs {
	return make(WorkspaceConfigs)
}

func NewWorkspaceConfig() *WorkspaceConfig {
	return &WorkspaceConfig{
		Repos: make(WorkspaceRepoConfigs),
	}
}

func (c *Config) Workspace(name string) (*WorkspaceConfig, error) {
	w, exists := c.Workspaces[name]
	if !exists {
		return nil, &WorkspaceNotFoundError{Name: name}
	}
	return w, nil
}

func (c *Config) AddWorkspace(name string) error {
	if _, err := c.Workspace(name); err == nil {
		return &WorkspaceExistsError{Name: name}
	}
	
	c.Workspaces[name] = NewWorkspaceConfig()
	return nil
}

func (c *Config) AddWorkspaceRepo(workspaceName, sourceAlias string, config *WorkspaceRepoConfig) error {
	if _, err := c.SourceByAlias(sourceAlias); err != nil {
		return err
	}

	w := c.ensureWorkspace(workspaceName)

	return w.AddRepo(sourceAlias, config)
}

func (c *Config) RemoveWorkspaceRepo(workspaceName, sourceAlias string) error {
	w, err := c.Workspace(workspaceName)
	if err != nil {
		return err
	}

	return w.RemoveRepo(sourceAlias)
}

func (c *Config) UpdateWorkspaceRepoBranch(workspaceName, sourceAlias, repoBranch string) error {
	w, err := c.Workspace(workspaceName)
	if err != nil {
		return err
	}

	if _, err := c.SourceByAlias(sourceAlias); err != nil {
		return err
	}

	return w.UpdateRepoBranch(sourceAlias, repoBranch)
}

func (c *Config) ensureWorkspace(name string) *WorkspaceConfig {
	if w, exists := c.Workspaces[name]; exists {
		return w
	}

	w := NewWorkspaceConfig()

	c.Workspaces[name] = w

	return w
}
