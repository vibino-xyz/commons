package vhttp

import "github.com/labstack/echo/v5"

// Controller registers its routes on the shared echo instance.
// Provide implementations through vhttp.WithControllers.
type Controller interface {
	Register(e *echo.Echo)
}
