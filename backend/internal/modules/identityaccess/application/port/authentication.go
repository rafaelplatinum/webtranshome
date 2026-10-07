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
	Issue(ctx context.Context, userID int64) (string, error)
}

type TokenVerifier interface {
	Verify(ctx context.Context, token string) (int64, error)
}
