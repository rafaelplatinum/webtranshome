package bcrypt

import (
	"errors"
	"testing"

	"webtranshome/internal/modules/identityaccess/application/port"

	"golang.org/x/crypto/bcrypt"
)

func TestVerifierAcceptsMatchingPassword(t *testing.T) {
	password := "test-password"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}

	if err := NewVerifier().Verify(string(hash), password); err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
}

func TestVerifierRejectsNonMatchingPassword(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}

	err = NewVerifier().Verify(string(hash), "incorrect-password")
	if !errors.Is(err, port.ErrPasswordMismatch) {
		t.Fatalf("Verify() error = %v, want password mismatch", err)
	}
}
