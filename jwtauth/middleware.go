package jwtauth

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
)

const contextKey = "jwtauth.claims"

func Middleware(v *Verifier) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			token := BearerToken(c)
			if token == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing bearer token")
			}

			claims, err := v.Verify(token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
			}

			c.Set(contextKey, claims)
			return next(c)
		}
	}
}

func BearerToken(c *echo.Context) string {
	const prefix = "Bearer "
	header := c.Request().Header.Get("Authorization")
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimPrefix(header, prefix)
}

func ClaimsFromContext(c *echo.Context) *Claims {
	claims, _ := c.Get(contextKey).(*Claims)
	return claims
}
