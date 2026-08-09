package vhttp

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

// Convenience constructors for common HTTP error responses so controllers don't
// repeat echo.NewHTTPError with bare status codes.

func BadRequest(message string) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, message)
}

func Unauthorized(message string) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusUnauthorized, message)
}

func Forbidden(message string) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusForbidden, message)
}

func NotFound(message string) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusNotFound, message)
}

func BadGateway(message string) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadGateway, message)
}

func ServiceUnavailable(message string) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusServiceUnavailable, message)
}

func Internal(message string) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusInternalServerError, message)
}
