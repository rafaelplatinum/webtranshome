package jwtissuer

import (
	"context"
	"testing"
	"time"

	"webtranshome/internal/modules/identityaccess/application/port"

	"github.com/golang-jwt/jwt/v4"
)

func TestVerifierAcceptsIssuedAccessToken(t *testing.T) {
	issuer := NewIssuer("test-secret", 60)
	token, err := issuer.Issue(context.Background(), 42)
	if err != nil {
		t.Fatalf("Issue() error = %v", err)
	}

	userID, err := NewVerifier("test-secret").Verify(context.Background(), token)
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	if userID != 42 {
		t.Fatalf("Verify() userID = %d, want 42", userID)
	}
}

func TestVerifierRejectsInvalidTokens(t *testing.T) {
	expired := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": 42,
		"sub":    "42",
		"exp":    time.Now().Add(-time.Minute).Unix(),
		"iat":    time.Now().Add(-time.Hour).Unix(),
	})
	expiredToken, err := expired.SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("SignedString() error = %v", err)
	}

	noExpiry := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": 42,
		"sub":    "42",
	})
	noExpiryToken, err := noExpiry.SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("SignedString() error = %v", err)
	}

	for _, token := range []string{
		"",
		"not-a-jwt",
		expiredToken,
		noExpiryToken,
	} {
		if _, err := NewVerifier("test-secret").Verify(context.Background(), token); err == nil {
			t.Errorf("Verify(%q) succeeded, want error", token)
		}
	}
}

func TestVerifierRejectsInvalidConfiguration(t *testing.T) {
	if _, err := NewVerifier("").Verify(context.Background(), "token"); err != port.ErrInvalidTokenConfiguration {
		t.Fatalf("Verify() error = %v, want %v", err, port.ErrInvalidTokenConfiguration)
	}
}
