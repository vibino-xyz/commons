package jwtauth

import (
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
)

const (
	EnvSecret    = "JWT_SECRET"
	EnvIssuer    = "JWT_ISSUER"
	EnvAccessTTL = "ACCESS_TOKEN_TTL"
)

const (
	DefaultIssuer    = "nexy"
	DefaultAccessTTL = 15 * time.Minute
)

const signingAlg = "HS256"

type Claims struct {
	gojwt.RegisteredClaims
	UserId         string `json:"uid"`
	OrganizationId string `json:"oid,omitempty"`
	OnboardingStep string `json:"step,omitempty"`
}
