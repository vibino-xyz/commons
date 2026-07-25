package jwtauth

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

// Organization role values carried in the access token's `role` claim. Any
// service that trusts nexy-issued tokens can authorize against these without
// calling back into nexy.
const (
	RoleOwner  = "OWNER"
	RoleAdmin  = "ADMIN"
	RoleMember = "MEMBER"
)

// RequireMember returns the caller's claims, requiring a valid token scoped to
// an organization. Use in handlers any organization member may call.
func RequireMember(c *echo.Context) (*Claims, error) {
	claims := ClaimsFromContext(c)
	if claims == nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	if claims.OrganizationId == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "no active organization")
	}
	return claims, nil
}

// RequireManager additionally requires an OWNER or ADMIN role.
func RequireManager(c *echo.Context) (*Claims, error) {
	claims, err := RequireMember(c)
	if err != nil {
		return nil, err
	}
	if claims.Role != RoleOwner && claims.Role != RoleAdmin {
		return nil, echo.NewHTTPError(http.StatusForbidden, "requires an admin or owner role")
	}
	return claims, nil
}
