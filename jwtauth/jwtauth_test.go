package jwtauth

import (
	"testing"
	"time"
)

func TestSignAndVerifyRoundTrip(t *testing.T) {
	signer, err := NewSigner("test-secret", "nexy", time.Minute)
	if err != nil {
		t.Fatalf("NewSigner: %v", err)
	}
	verifier, err := NewVerifier("test-secret", "nexy")
	if err != nil {
		t.Fatalf("NewVerifier: %v", err)
	}

	token, expiresAt, err := signer.Sign("usr_123", "org_456", "COMPLETE", RoleOwner)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	if !expiresAt.After(time.Now()) {
		t.Fatalf("expiresAt should be in the future, got %v", expiresAt)
	}

	claims, err := verifier.Verify(token)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if claims.UserId != "usr_123" || claims.OrganizationId != "org_456" || claims.OnboardingStep != "COMPLETE" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
	if claims.Role != RoleOwner {
		t.Fatalf("expected role %q, got %q", RoleOwner, claims.Role)
	}
	if claims.Subject != "usr_123" || claims.Issuer != "nexy" {
		t.Fatalf("unexpected registered claims: %+v", claims.RegisteredClaims)
	}
}

func TestVerifyRejectsWrongSecret(t *testing.T) {
	signer, _ := NewSigner("secret-a", "nexy", time.Minute)
	verifier, _ := NewVerifier("secret-b", "nexy")

	token, _, _ := signer.Sign("usr_1", "", "PASSWORD", "")
	if _, err := verifier.Verify(token); err == nil {
		t.Fatal("expected verification to fail with a mismatched secret")
	}
}

func TestVerifyRejectsWrongIssuer(t *testing.T) {
	signer, _ := NewSigner("secret", "other", time.Minute)
	verifier, _ := NewVerifier("secret", "nexy")

	token, _, _ := signer.Sign("usr_1", "", "PASSWORD", "")
	if _, err := verifier.Verify(token); err == nil {
		t.Fatal("expected verification to fail with a mismatched issuer")
	}
}

func TestVerifyRejectsExpiredToken(t *testing.T) {
	// A tiny positive TTL (non-positive is floored to the default by NewSigner).
	signer, _ := NewSigner("secret", "nexy", time.Millisecond)
	verifier, _ := NewVerifier("secret", "nexy")

	token, _, _ := signer.Sign("usr_1", "", "PASSWORD", "")
	time.Sleep(20 * time.Millisecond)
	if _, err := verifier.Verify(token); err == nil {
		t.Fatal("expected verification to fail for an expired token")
	}
}

func TestVerifyRejectsTamperedToken(t *testing.T) {
	signer, _ := NewSigner("secret", "nexy", time.Minute)
	verifier, _ := NewVerifier("secret", "nexy")

	token, _, _ := signer.Sign("usr_1", "", "PASSWORD", "")
	if _, err := verifier.Verify(token + "x"); err == nil {
		t.Fatal("expected verification to fail for a tampered token")
	}
}

func TestNewSignerAndVerifierRequireSecret(t *testing.T) {
	if _, err := NewSigner("", "nexy", time.Minute); err == nil {
		t.Fatal("expected NewSigner to require a secret")
	}
	if _, err := NewVerifier("", "nexy"); err == nil {
		t.Fatal("expected NewVerifier to require a secret")
	}
}
