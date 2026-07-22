package whttp

import "go.uber.org/fx"

const groupController = `group:"controller"`

var serverFx = fx.Options(
	fx.Provide(
		newServerSettings,
		// createHTTPServer consumes the grouped controllers.
		fx.Annotate(
			createHTTPServer,
			fx.ParamTags(groupController, ""),
		),
	),
	fx.Invoke(startHTTPServer),
)

// DefaultServer wires a basic echo HTTP server into an fx application.
//
//	fx.New(
//		whttp.DefaultServer(
//			whttp.WithControllers(NewFooController),
//		),
//	).Run()
func DefaultServer(opts ...Options) fx.Option {
	fxOpts := []fx.Option{serverFx}
	for _, opt := range opts {
		fxOpts = append(fxOpts, opt.apply())
	}
	return fx.Options(fxOpts...)
}

func registerControllers(controllers ...interface{}) fx.Option {
	a := make([]interface{}, len(controllers))
	for i, c := range controllers {
		a[i] = fx.Annotate(c, fx.ResultTags(groupController))
	}
	return fx.Provide(a...)
}
