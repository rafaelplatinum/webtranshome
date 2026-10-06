package user

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("IDENTITY_USER_NOT_FOUND")

type Repository interface {
	FindForAuthentication(ctx context.Context, email string) (User, error)
}
