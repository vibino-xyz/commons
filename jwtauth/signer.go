package jwtauth

import (
	"errors"
	"os"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
)

type Signer struct {
	secret    []byte
	issuer    string
	accessTTL time.Duration
}

func NewSigner(secret, issuer string, accessTTL time.Duration) (*Signer, error) {
	if secret == "" {
		return nil, errors.New("jwtauth: secret is required")
	}
	if issuer == "" {
		issuer = DefaultIssuer
	}
	if accessTTL <= 0 {
		accessTTL = DefaultAccessTTL
	}
	return &Signer{secret: []byte(secret), issuer: issuer, accessTTL: accessTTL}, nil
}

func NewSignerFromEnv() (*Signer, error) {
	return NewSigner(os.Getenv(EnvSecret), os.Getenv(EnvIssuer), durationFromEnv(EnvAccessTTL, DefaultAccessTTL))
}

func (s *Signer) AccessTTL() time.Duration { return s.accessTTL }

func (s *Signer) Sign(userId, organizationId, onboardingStep, role string) (string, time.Time, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(s.accessTTL)

	claims := Claims{
		RegisteredClaims: gojwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   userId,
			IssuedAt:  gojwt.NewNumericDate(now),
			ExpiresAt: gojwt.NewNumericDate(expiresAt),
		},
		UserId:         userId,
		OrganizationId: organizationId,
		OnboardingStep: onboardingStep,
		Role:           role,
	}

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expiresAt, nil
}

func durationFromEnv(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
