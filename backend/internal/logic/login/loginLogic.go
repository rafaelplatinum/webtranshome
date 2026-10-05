package login

import (
	"context"
	"errors"
	"math"
	"strconv"
	"time"

	"webtranshome/internal/model"
	"webtranshome/internal/svc"
	"webtranshome/internal/types"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	if req == nil || req.Email == "" || req.Password == "" {
		return nil, ErrInvalidRequest
	}

	user, err := l.svcCtx.UserModel.FindOneByEmail(l.ctx, req.Email)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !user.IsActive || !user.PasswordHash.Valid {
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(req.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := l.generateToken(user.Id)
	if err != nil {
		return nil, err
	}

	return &types.LoginResponse{
		AccessToken: token,
	}, nil
}

func (l *LoginLogic) generateToken(userId int64) (string, error) {
	secret := l.svcCtx.Config.Auth.AccessSecret
	expire := l.svcCtx.Config.Auth.AccessExpire
	if secret == "" {
		return "", ErrJWTConfiguration
	}
	if expire <= 0 || expire > math.MaxInt64/int64(time.Second) {
		return "", ErrJWTConfiguration
	}

	now := time.Now()
	claims := struct {
		UserID int64 `json:"userId"`
		jwt.RegisteredClaims
	}{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userId, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expire) * time.Second)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
