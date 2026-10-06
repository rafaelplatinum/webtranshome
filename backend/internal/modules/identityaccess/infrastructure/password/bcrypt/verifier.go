package bcrypt

import (
	"errors"

	"webtranshome/internal/modules/identityaccess/application/port"

	"golang.org/x/crypto/bcrypt"
)

type Verifier struct{}

func NewVerifier() Verifier {
	return Verifier{}
}

func (Verifier) Verify(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return port.ErrPasswordMismatch
	}
	return err
}
