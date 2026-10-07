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

func (v *Verifier) Verify(_ context.Context, tokenString string) (int64, error) {
	if v.secret == "" {
		return 0, port.ErrInvalidTokenConfiguration
	}

	claims := struct {
		UserID int64 `json:"userId"`
		jwt.RegisteredClaims
	}{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, port.ErrInvalidAccessToken
		}
		return []byte(v.secret), nil
	})
	if err != nil || token == nil || !token.Valid {
		return 0, port.ErrInvalidAccessToken
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil || userID <= 0 || claims.UserID != userID || claims.ExpiresAt == nil {
		return 0, port.ErrInvalidAccessToken
	}

	return userID, nil
}

var _ port.TokenVerifier = (*Verifier)(nil)
