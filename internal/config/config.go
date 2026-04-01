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

type (
	Repos      map[string]*RepoConfig
	Workspaces map[string]*WorkspaceConfig
)

type Config struct {
	Repos      Repos      `yaml:"repos,omitempty"`
	Workspaces Workspaces `yaml:"workspaces,omitempty"`
}

type RepoConfig struct {
	Path   string      `yaml:"path"`
	Branch string      `yaml:"branch,omitempty"`
	Hooks  HooksConfig `yaml:"hooks,omitempty"`
}

type HooksConfig struct {
	Setup []string `yaml:"setup,omitempty"`
}

type WorkspaceConfig struct {
	Repos map[string]*WorkspaceRepoConfig `yaml:"repos"`
}

type WorkspaceRepoConfig struct {
	Branch string `yaml:"branch"`
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

func New() *Config {
	return &Config{
		Repos:      make(Repos),
		Workspaces: make(Workspaces),
	}
}

func Load() (*Config, error) {
	p, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %w (run 'ws init' first)", p, err)
	}

	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	if c.Repos == nil {
		c.Repos = make(Repos)
	}

	if c.Workspaces == nil {
		c.Workspaces = make(Workspaces)
	}

	return &c, nil
}

func (c *Config) Save() error {
	p, err := ConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}

	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)

	if err := encoder.Encode(c); err != nil {
		return err
	}

	return os.WriteFile(p, buf.Bytes(), 0o644)
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
