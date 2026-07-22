package id

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

const (
	randomBytes = 16
	PrefixLen   = 3
	Length      = PrefixLen + 1 + randomBytes*2
)

func New(prefix string) (string, error) {
	if err := validatePrefix(prefix); err != nil {
		return "", err
	}

	b := make([]byte, randomBytes)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate id: %w", err)
	}

	return prefix + "_" + hex.EncodeToString(b), nil
}

func validatePrefix(prefix string) error {
	if len(prefix) != PrefixLen {
		return fmt.Errorf("id prefix %q must be exactly %d letters", prefix, PrefixLen)
	}
	for _, r := range prefix {
		if r < 'a' || r > 'z' {
			return fmt.Errorf("id prefix %q must contain only lowercase letters", prefix)
		}
	}
	return nil
}
