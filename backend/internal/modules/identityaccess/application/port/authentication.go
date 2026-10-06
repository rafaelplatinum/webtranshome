package port

import (
	"context"
	"errors"
)

var (
	ErrInvalidTokenConfiguration = errors.New("AUTH_JWT_CONFIGURATION_INVALID")
	ErrPasswordMismatch          = errors.New("AUTH_PASSWORD_MISMATCH")
)

type PasswordVerifier interface {
	Verify(hash, password string) error
}

type TokenIssuer interface {
	Issue(ctx context.Context, userID int64) (string, error)
}
