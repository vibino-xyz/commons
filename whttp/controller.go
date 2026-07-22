package whttp

import "github.com/labstack/echo/v5"

// Controller registers its routes on the shared echo instance.
// Provide implementations through whttp.WithControllers.
type Controller interface {
	Register(e *echo.Echo)
}
