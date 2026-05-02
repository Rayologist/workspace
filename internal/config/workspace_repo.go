package config

type WorkspaceRepoConfig struct {
	Branch string `yaml:"branch"`
}

type WorkspaceConfigs map[string]*WorkspaceConfig

func (w *WorkspaceConfig) Repo(sourceAlias string) (*WorkspaceRepoConfig, error) {
	repo, exists := w.Repos[sourceAlias]
	if !exists {
		return nil, &WorkspaceRepoNotFoundError{
			SourceAlias: sourceAlias,
		}
	}

	return repo, nil
}

func (w *WorkspaceConfig) AddRepo(sourceAlias string, config *WorkspaceRepoConfig) error {
	w.Repos[sourceAlias] = config
	return nil
}

func (w *WorkspaceConfig) UpdateRepoBranch(sourceAlias, repoBranch string) error {
	repo, err := w.Repo(sourceAlias)
	if err != nil {
		return err
	}

	repo.Branch = repoBranch

	return nil
}

func (w *WorkspaceConfig) RemoveRepo(sourceAlias string) error {
	_, err := w.Repo(sourceAlias)
	if err != nil {
		return err
	}

	delete(w.Repos, sourceAlias)
	return nil
}
