package jwtissuer

import (
	"context"
	"strconv"
	"testing"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v4"
)

func TestIssuerSignsTokenWithUserAndExpiryClaims(t *testing.T) {
	issuer := NewIssuer("test-secret", 60)
	signedToken, err := issuer.Issue(context.Background(), 42)
	if err != nil {
		t.Fatalf("Issue() error = %v", err)
	}

	claims := &struct {
		UserID int64 `json:"userId"`
		jwtlib.RegisteredClaims
	}{}
	token, err := jwtlib.ParseWithClaims(signedToken, claims, func(token *jwtlib.Token) (interface{}, error) {
		if token.Method != jwtlib.SigningMethodHS256 {
			t.Fatalf("signing method = %v, want HS256", token.Method)
		}
		return []byte("test-secret"), nil
	})
	if err != nil {
		t.Fatalf("ParseWithClaims() error = %v", err)
	}
	if !token.Valid {
		t.Fatal("issued token is invalid")
	}
	if claims.UserID != 42 || claims.Subject != strconv.FormatInt(42, 10) {
		t.Fatalf("claims identify user %d with subject %q, want user 42", claims.UserID, claims.Subject)
	}
	if claims.ExpiresAt == nil || !claims.ExpiresAt.Time.After(time.Now()) {
		t.Fatal("token has no future expiration")
	}
}

func TestIssuerRejectsInvalidConfiguration(t *testing.T) {
	for _, issuer := range []*Issuer{
		NewIssuer("", 60),
		NewIssuer("test-secret", 0),
		NewIssuer("test-secret", -1),
	} {
		if _, err := issuer.Issue(context.Background(), 42); err == nil {
			t.Fatal("Issue() succeeded with invalid JWT configuration")
		}
	}
}
