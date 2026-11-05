package jwtissuer

import (
	"context"
	"math"
	"strconv"
	"time"

	"webtranshome/internal/modules/identityaccess/application/port"

	"github.com/golang-jwt/jwt/v4"
)

type Issuer struct {
	secret     string
	expireSecs int64
}

func NewIssuer(secret string, expireSecs int64) *Issuer {
	return &Issuer{
		secret:     secret,
		expireSecs: expireSecs,
	}
}

func (i *Issuer) Issue(_ context.Context, userID int64, authVersion int64) (string, error) {
	if i.secret == "" || i.expireSecs <= 0 || i.expireSecs > math.MaxInt64/int64(time.Second) {
		return "", port.ErrInvalidTokenConfiguration
	}

	now := time.Now()
	claims := struct {
		UserID      int64 `json:"userId"`
		AuthVersion int64 `json:"authVersion"`
		jwt.RegisteredClaims
	}{
		UserID: userID, AuthVersion: authVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(i.expireSecs) * time.Second)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(i.secret))
}

var _ port.TokenIssuer = (*Issuer)(nil)
