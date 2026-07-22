package whttp

import "go.uber.org/fx"

// Options configures the DefaultServer.
type Options interface {
	apply() fx.Option
}

type fnOption func() fx.Option

func (f fnOption) apply() fx.Option { return f() }

// WithControllers registers one or more controller constructors. Each constructor
// must return a whttp.Controller.
func WithControllers(controllers ...interface{}) Options {
	return fnOption(func() fx.Option {
		return registerControllers(controllers...)
	})
}
