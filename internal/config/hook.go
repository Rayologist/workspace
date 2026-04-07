package config

import (
	"strings"

	"workspace/internal/set"
)

type HooksConfig struct {
	Setup []string `yaml:"setup,omitempty"`
}

func (h *HooksConfig) AppendSetupHooks(hooks []string) {
	hooksSet := set.FromSlice(h.Setup)

	for _, hook := range hooks {
		trimmed := strings.TrimSpace(hook)
		if !hooksSet.Contains(trimmed) {
			h.Setup = append(h.Setup, trimmed)
			hooksSet.Add(hook)
		}
	}
}
