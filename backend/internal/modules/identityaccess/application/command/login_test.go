package command

import (
	"context"
	"errors"
	"testing"

	"webtranshome/internal/modules/identityaccess/application/port"
	"webtranshome/internal/modules/identityaccess/domain/user"
)

type userRepositoryStub struct {
	account user.User
	err     error
}

func (s userRepositoryStub) FindForAuthentication(context.Context, string) (user.User, error) {
	return s.account, s.err
}

type passwordVerifierStub struct {
	err      error
	hash     string
	password string
}

func (s *passwordVerifierStub) Verify(hash, password string) error {
	s.hash = hash
	s.password = password
	return s.err
}

type tokenIssuerStub struct {
	token  string
	userID int64
	err    error
}

func (s *tokenIssuerStub) Issue(_ context.Context, userID int64) (string, error) {
	s.userID = userID
	return s.token, s.err
}

func TestLoginExecuteIssuesTokenForActiveAdmin(t *testing.T) {
	verifier := &passwordVerifierStub{}
	issuer := &tokenIssuerStub{token: "signed-token"}
	login := NewLogin(
		userRepositoryStub{
			account: user.User{
				ID:           42,
				Email:        "admin@example.test",
				PasswordHash: "stored-hash",
				IsActive:     true,
				RoleCodes:    []string{user.RoleAdminCatalog},
			},
		},
		verifier,
		issuer,
	)

	result, err := login.Execute(context.Background(), LoginCommand{
		Email:    " admin@example.test ",
		Password: "submitted-password",
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if result.AccessToken != "signed-token" {
		t.Fatalf("AccessToken = %q, want token from issuer", result.AccessToken)
	}
	if verifier.hash != "stored-hash" || verifier.password != "submitted-password" {
		t.Fatal("password verifier did not receive the stored hash and submitted password")
	}
	if issuer.userID != 42 {
		t.Fatalf("token issuer userID = %d, want 42", issuer.userID)
	}
}

func TestLoginExecuteRejectsNonAdminWithoutVerifyingPassword(t *testing.T) {
	verifier := &passwordVerifierStub{}
	login := NewLogin(
		userRepositoryStub{
			account: user.User{
				ID:           42,
				PasswordHash: "stored-hash",
				IsActive:     true,
				RoleCodes:    []string{"MEMBER"},
			},
		},
		verifier,
		&tokenIssuerStub{},
	)

	_, err := login.Execute(context.Background(), LoginCommand{
		Email:    "user@example.test",
		Password: "submitted-password",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("Execute() error = %v, want ErrInvalidCredentials", err)
	}
	if verifier.hash != "" {
		t.Fatal("password verifier was called for a non-admin user")
	}
}

func TestLoginExecuteRejectsInvalidPassword(t *testing.T) {
	verifier := &passwordVerifierStub{err: port.ErrPasswordMismatch}
	issuer := &tokenIssuerStub{}
	login := NewLogin(
		userRepositoryStub{
			account: user.User{
				ID:           42,
				PasswordHash: "stored-hash",
				IsActive:     true,
				RoleCodes:    []string{user.RoleSuperAdmin},
			},
		},
		verifier,
		issuer,
	)

	_, err := login.Execute(context.Background(), LoginCommand{
		Email:    "admin@example.test",
		Password: "wrong-password",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("Execute() error = %v, want ErrInvalidCredentials", err)
	}
	if issuer.userID != 0 {
		t.Fatal("token issuer was called for invalid credentials")
	}
}

func TestLoginExecutePropagatesPasswordVerifierFailures(t *testing.T) {
	verifierErr := errors.New("verifier unavailable")
	login := NewLogin(
		userRepositoryStub{
			account: user.User{
				ID:           42,
				PasswordHash: "stored-hash",
				IsActive:     true,
				RoleCodes:    []string{user.RoleSuperAdmin},
			},
		},
		&passwordVerifierStub{err: verifierErr},
		&tokenIssuerStub{},
	)

	_, err := login.Execute(context.Background(), LoginCommand{
		Email:    "admin@example.test",
		Password: "submitted-password",
	})
	if !errors.Is(err, verifierErr) {
		t.Fatalf("Execute() error = %v, want verifier error", err)
	}
}

func TestLoginExecutePropagatesRepositoryFailures(t *testing.T) {
	repositoryErr := errors.New("repository unavailable")
	login := NewLogin(
		userRepositoryStub{err: repositoryErr},
		&passwordVerifierStub{},
		&tokenIssuerStub{},
	)

	_, err := login.Execute(context.Background(), LoginCommand{
		Email:    "admin@example.test",
		Password: "submitted-password",
	})
	if !errors.Is(err, repositoryErr) {
		t.Fatalf("Execute() error = %v, want repository error", err)
	}
}
