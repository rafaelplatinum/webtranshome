package svc

import (
	"webtranshome/internal/config"
	"webtranshome/internal/model"

	"github.com/zeromicro/go-zero/core/stores/postgres"
)

type ServiceContext struct {
	Config    config.Config
	UserModel model.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Database.DataSource)

	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUsersModel(conn),
	}
}
