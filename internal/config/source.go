package config

import "fmt"

type SourceConfig struct {
	Path   string      `yaml:"path"`
	Branch string      `yaml:"branch,omitempty"`
	Hooks  HooksConfig `yaml:"hooks,omitempty"`
}

type SourceConfigs map[string]*SourceConfig

func NewSourceConfigs() SourceConfigs {
	return make(SourceConfigs)
}

func (c *Config) SourceByAlias(alias string) (*SourceConfig, error) {
	source, exists := c.Sources[alias]
	if !exists {
		return nil, fmt.Errorf("source '%s' not exist in the config", alias)
	}
	return source, nil
}

func (c *Config) RemoveSource(alias string) error {
	if _, err := c.SourceByAlias(alias); err != nil {
		return err
	}

	delete(c.Sources, alias)

	return nil
}

func (c *Config) AddSource(alias string, source *SourceConfig) error {
	if _, err := c.SourceByAlias(alias); err == nil {
		return fmt.Errorf("source '%s' already exists (use 'source update' to modify it)", alias)
	}

	c.Sources[alias] = source
	return nil
}

func (c *Config) UpdateSource(alias, newAlias string, source *SourceConfig) error {
	if _, err := c.SourceByAlias(alias); err != nil {
		return err
	}

	targetAlias := alias
	if newAlias != "" {
		if _, err := c.SourceByAlias(newAlias); err == nil {
			return fmt.Errorf("source alias '%s' already exists", newAlias)
		}

		delete(c.Sources, alias)
		targetAlias = newAlias
	}

	c.Sources[targetAlias] = source
	return nil
}
