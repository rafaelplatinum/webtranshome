package jwtissuer

import (
	"context"
	"strconv"

	"webtranshome/internal/modules/identityaccess/application/port"

	"github.com/golang-jwt/jwt/v4"
)

type Verifier struct {
	secret string
}

func NewVerifier(secret string) *Verifier {
	return &Verifier{secret: secret}
}

func (v *Verifier) Verify(_ context.Context, tokenString string) (port.VerifiedToken, error) {
	if v.secret == "" {
		return port.VerifiedToken{}, port.ErrInvalidTokenConfiguration
	}

	claims := struct {
		UserID      int64 `json:"userId"`
		AuthVersion int64 `json:"authVersion"`
		jwt.RegisteredClaims
	}{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, port.ErrInvalidAccessToken
		}
		return []byte(v.secret), nil
	})
	if err != nil || token == nil || !token.Valid {
		return port.VerifiedToken{}, port.ErrInvalidAccessToken
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil || userID <= 0 || claims.UserID != userID || claims.ExpiresAt == nil {
		return port.VerifiedToken{}, port.ErrInvalidAccessToken
	}

	return port.VerifiedToken{UserID: userID, AuthVersion: claims.AuthVersion}, nil
}

var _ port.TokenVerifier = (*Verifier)(nil)
