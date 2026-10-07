package svc

import (
	"webtranshome/internal/config"
	"webtranshome/internal/middleware"
	"webtranshome/internal/model"
	"webtranshome/internal/modules/identityaccess/application/command"
	"webtranshome/internal/modules/identityaccess/application/query"
	passwordbcrypt "webtranshome/internal/modules/identityaccess/infrastructure/password/bcrypt"
	authpostgres "webtranshome/internal/modules/identityaccess/infrastructure/postgres"
	jwtissuer "webtranshome/internal/modules/identityaccess/infrastructure/token/jwtissuer"

	"github.com/zeromicro/go-zero/core/stores/postgres"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config                   config.Config
	UserModel                model.UsersModel
	LoginAdmin               *command.Login
	AccessControlQueries     *query.AccessControl
	ReplaceAccessAssignments *command.ReplaceAccessAssignments
	ManageRoles              *command.ManageRoles
	Authenticate             rest.Middleware
	RequirePermission        func(string) rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Database.DataSource)
	userModel := model.NewUsersModel(conn)
	userRepository := authpostgres.NewUserRepository(conn, userModel)
	permissionRepository := authpostgres.NewPermissionRepository(conn)
	accessControlRepository := authpostgres.NewAccessControlRepository(conn)
	loginAdmin := command.NewLogin(
		userRepository,
		passwordbcrypt.NewVerifier(),
		jwtissuer.NewIssuer(c.Auth.AccessSecret, c.Auth.AccessExpire),
	)
	checkPermission := query.NewCheckPermission(permissionRepository)
	tokenVerifier := jwtissuer.NewVerifier(c.Auth.AccessSecret)
	accessControlQueries := query.NewAccessControl(
		accessControlRepository,
		accessControlRepository,
		accessControlRepository,
	)
	replaceAccessAssignments := command.NewReplaceAccessAssignments(accessControlRepository)
	manageRoles := command.NewManageRoles(accessControlRepository)

	return &ServiceContext{
		Config:                   c,
		UserModel:                userModel,
		LoginAdmin:               loginAdmin,
		AccessControlQueries:     accessControlQueries,
		ReplaceAccessAssignments: replaceAccessAssignments,
		ManageRoles:              manageRoles,
		Authenticate:             middleware.Authenticate(tokenVerifier),
		RequirePermission: func(code string) rest.Middleware {
			return middleware.RequirePermission(checkPermission, code)
		},
	}
}
