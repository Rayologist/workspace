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
	Sources    SourceConfigs    `yaml:"sources,omitempty"`
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
		Sources:    NewSourceConfigs(),
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
	if c.Sources == nil {
		c.Sources = NewSourceConfigs()
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

func (c *Config) SourceAbsPath(alias string) (string, error) {
	source, ok := c.Sources[alias]
	if !ok {
		return "", fmt.Errorf("source '%s' not found in config", alias)
	}

	if filepath.IsAbs(source.Path) {
		return source.Path, nil
	}

	wsDir, err := WorkspacesDirPath()
	if err != nil {
		return "", err
	}

	abs, err := filepath.Abs(filepath.Join(wsDir, source.Path))
	if err != nil {
		return "", err
	}

	return abs, nil
}
