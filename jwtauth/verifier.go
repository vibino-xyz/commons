package jwtauth

import (
	"errors"
	"os"

	gojwt "github.com/golang-jwt/jwt/v5"
)

type Verifier struct {
	secret []byte
	issuer string
}

func NewVerifier(secret, issuer string) (*Verifier, error) {
	if secret == "" {
		return nil, errors.New("jwtauth: secret is required")
	}
	if issuer == "" {
		issuer = DefaultIssuer
	}
	return &Verifier{secret: []byte(secret), issuer: issuer}, nil
}

func NewVerifierFromEnv() (*Verifier, error) {
	return NewVerifier(os.Getenv(EnvSecret), os.Getenv(EnvIssuer))
}

func (v *Verifier) Verify(tokenString string) (*Claims, error) {
	claims := &Claims{}
	_, err := gojwt.ParseWithClaims(tokenString, claims, func(t *gojwt.Token) (any, error) {
		if _, ok := t.Method.(*gojwt.SigningMethodHMAC); !ok {
			return nil, errors.New("jwtauth: unexpected signing method")
		}
		return v.secret, nil
	}, gojwt.WithIssuer(v.issuer), gojwt.WithValidMethods([]string{signingAlg}))
	if err != nil {
		return nil, err
	}
	return claims, nil
}
