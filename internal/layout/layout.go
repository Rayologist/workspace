package layout

import (
	"os"
	"path/filepath"
)

type Layout struct {
	root          string
	workspacesDir string
	configFile    string
}

type Option func(*Layout) error

func WithConfigPath(path string) Option {
	return func(l *Layout) error {
		dir := filepath.Dir(path)
		base := filepath.Base(path)

		root, err := filepath.Abs(filepath.Join(dir, ".."))
		if err != nil {
			return err
		}

		l.root = root
		l.workspacesDir = dir
		l.configFile = base

		return nil
	}
}

func New(opts ...Option) (*Layout, error) {
	l := &Layout{
		root:          "",
		workspacesDir: "workspaces",
		configFile:    "ws.yaml",
	}

	for _, opt := range opts {
		err := opt(l)
		if err != nil {
			return nil, err
		}
	}

	if l.root == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		l.root = cwd
	}

	return l, nil
}

func (l *Layout) Root() string {
	return l.root
}

func (l *Layout) WorkspacesDir() string {
	return filepath.Join(l.Root(), l.workspacesDir)
}

func (l *Layout) ConfigPath() string {
	return filepath.Join(l.WorkspacesDir(), l.configFile)
}
