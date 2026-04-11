package cli

import (
	"workspace/internal/config"
	"workspace/internal/layout"
)

type Runtime struct {
	config    *config.Config
	layout    *layout.Layout
	IOStreams IOStreams
}

type Option func(*Runtime) error

func WithLayout(opts ...layout.Option) Option {
	return func(r *Runtime) error {
		l, err := layout.New(opts...)
		if err != nil {
			return err
		}
		r.layout = l
		return nil
	}
}

func NewRuntime() *Runtime {
	return &Runtime{
		IOStreams: NewSystemIOStreams(),
	}
}

func (r *Runtime) Init(opts ...Option) error {
	for _, opt := range opts {
		err := opt(r)
		if err != nil {
			return err
		}
	}

	if r.layout == nil {
		l, err := layout.New()
		if err != nil {
			return err
		}
		r.layout = l
	}

	return nil
}

func (r *Runtime) Config() (*config.Config, error) {
	if r.config != nil {
		return r.config, nil
	}

	cfg, err := config.Load(r.layout.ConfigPath())
	if err != nil {
		return nil, err
	}

	r.config = cfg
	return r.config, nil
}

func (r *Runtime) Layout() *layout.Layout {
	return r.layout
}
