package vhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v5"
	"go.uber.org/fx"
)

const defaultPort = 8080

type serverSettings struct {
	Port int
}

// newServerSettings reads the HTTP port from the PORT env var, falling back to defaultPort.
func newServerSettings() *serverSettings {
	port := defaultPort
	if v := os.Getenv("PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			port = p
		}
	}
	return &serverSettings{Port: port}
}

func createHTTPServer(controllers []Controller, settings *serverSettings) *http.Server {
	e := echo.New()

	for _, controller := range controllers {
		controller.Register(e)
	}

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Port),
		Handler: e,
	}
}

func startHTTPServer(lc fx.Lifecycle, srv *http.Server) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
}
