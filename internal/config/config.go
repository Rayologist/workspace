package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v3"
)

const (
	WorkspacesDir = "workspaces"
	ConfigFile    = "ws.yaml"
)

type Config struct {
	path       string
	Repos      RepoConfigs      `yaml:"repos,omitempty"`
	Workspaces WorkspaceConfigs `yaml:"workspaces,omitempty"`
}

func WorkspacesDirPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, WorkspacesDir), nil
}

func ConfigPath() (string, error) {
	wsDir, err := WorkspacesDirPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(wsDir, ConfigFile), nil
}

func New(path string) *Config {
	return &Config{
		path:       path,
		Repos:      NewRepoConfigs(),
		Workspaces: NewWorkspaceConfigs(),
	}
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("cannot read %s (run 'ws init' first)", path)
	}

	if err != nil {
		return nil, err
	}

	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	c.path = path
	if c.Repos == nil {
		c.Repos = NewRepoConfigs()
	}

	if c.Workspaces == nil {
		c.Workspaces = NewWorkspaceConfigs()
	}

	return &c, nil
}

func (c *Config) Save() error {
	if c.path == "" {
		return fmt.Errorf("config path is not set")
	}

	if err := os.MkdirAll(filepath.Dir(c.path), 0o755); err != nil {
		return err
	}

	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)

	if err := encoder.Encode(c); err != nil {
		return err
	}

	return os.WriteFile(c.path, buf.Bytes(), 0o644)
}

func (c *Config) RepoAbsPath(alias string) (string, error) {
	repo, ok := c.Repos[alias]
	if !ok {
		return "", fmt.Errorf("repo '%s' not found in config", alias)
	}

	if filepath.IsAbs(repo.Path) {
		return repo.Path, nil
	}

	wsDir, err := WorkspacesDirPath()
	if err != nil {
		return "", err
	}

	abs, err := filepath.Abs(filepath.Join(wsDir, repo.Path))
	if err != nil {
		return "", err
	}

	return abs, nil
}
