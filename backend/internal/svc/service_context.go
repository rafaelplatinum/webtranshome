package svc

import (
	"webtranshome/internal/config"
	"webtranshome/internal/model"
	"webtranshome/internal/modules/identityaccess/application/command"
	passwordbcrypt "webtranshome/internal/modules/identityaccess/infrastructure/password/bcrypt"
	authpostgres "webtranshome/internal/modules/identityaccess/infrastructure/postgres"
	jwtissuer "webtranshome/internal/modules/identityaccess/infrastructure/token/jwtissuer"

	"github.com/zeromicro/go-zero/core/stores/postgres"
)

type ServiceContext struct {
	Config     config.Config
	UserModel  model.UsersModel
	LoginAdmin *command.Login
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Database.DataSource)
	userModel := model.NewUsersModel(conn)
	userRepository := authpostgres.NewUserRepository(conn, userModel)
	loginAdmin := command.NewLogin(
		userRepository,
		passwordbcrypt.NewVerifier(),
		jwtissuer.NewIssuer(c.Auth.AccessSecret, c.Auth.AccessExpire),
	)

	return &ServiceContext{
		Config:     c,
		UserModel:  userModel,
		LoginAdmin: loginAdmin,
	}
}
