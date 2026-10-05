package login

import "errors"

var (
	ErrInvalidCredentials = errors.New("AUTH_INVALID_CREDENTIALS")
	ErrInvalidRequest     = errors.New("AUTH_INVALID_REQUEST")
	ErrJWTConfiguration   = errors.New("AUTH_JWT_CONFIGURATION_INVALID")
)
