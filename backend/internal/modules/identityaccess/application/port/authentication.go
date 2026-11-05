package port

import (
	"context"
	"errors"
)

var (
	ErrInvalidTokenConfiguration = errors.New("AUTH_JWT_CONFIGURATION_INVALID")
	ErrInvalidAccessToken        = errors.New("AUTH_ACCESS_TOKEN_INVALID")
	ErrPasswordMismatch          = errors.New("AUTH_PASSWORD_MISMATCH")
)

type PasswordVerifier interface {
	Verify(hash, password string) error
}

type TokenIssuer interface {
	Issue(ctx context.Context, userID int64, authVersion int64) (string, error)
}

type VerifiedToken struct {
	UserID      int64
	AuthVersion int64
}

type TokenVerifier interface {
	Verify(ctx context.Context, token string) (VerifiedToken, error)
}
