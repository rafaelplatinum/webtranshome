package command

import (
	"context"
	"errors"
	"strings"

	"webtranshome/internal/modules/identityaccess/application/port"
	"webtranshome/internal/modules/identityaccess/domain/user"
)

var (
	ErrInvalidLoginRequest = errors.New("AUTH_INVALID_REQUEST")
	ErrInvalidCredentials  = errors.New("AUTH_INVALID_CREDENTIALS")
)

type LoginCommand struct {
	Email    string
	Password string
}

type LoginResult struct {
	AccessToken string
}

type Login struct {
	users            user.Repository
	passwordVerifier port.PasswordVerifier
	tokenIssuer      port.TokenIssuer
}

func NewLogin(
	users user.Repository,
	passwordVerifier port.PasswordVerifier,
	tokenIssuer port.TokenIssuer,
) *Login {
	return &Login{
		users:            users,
		passwordVerifier: passwordVerifier,
		tokenIssuer:      tokenIssuer,
	}
}

func (l *Login) Execute(ctx context.Context, command LoginCommand) (LoginResult, error) {
	if strings.TrimSpace(command.Email) == "" || command.Password == "" {
		return LoginResult{}, ErrInvalidLoginRequest
	}

	account, err := l.users.FindForAuthentication(ctx, strings.TrimSpace(command.Email))
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return LoginResult{}, ErrInvalidCredentials
		}
		return LoginResult{}, err
	}

	if !account.IsActive || !account.HasStaffRole || account.PasswordHash == "" {
		return LoginResult{}, ErrInvalidCredentials
	}

	if err := l.passwordVerifier.Verify(account.PasswordHash, command.Password); err != nil {
		if errors.Is(err, port.ErrPasswordMismatch) {
			return LoginResult{}, ErrInvalidCredentials
		}
		return LoginResult{}, err
	}

	token, err := l.tokenIssuer.Issue(ctx, account.ID)
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{AccessToken: token}, nil
}
