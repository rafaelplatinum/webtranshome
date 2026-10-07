package postgres

import (
	"webtranshome/internal/modules/identityaccess/domain/service"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type AccessControlRepository struct {
	connection sqlx.SqlConn
}

func NewAccessControlRepository(connection sqlx.SqlConn) *AccessControlRepository {
	return &AccessControlRepository{connection: connection}
}

var (
	_ service.AccessControlCatalogRepository = (*AccessControlRepository)(nil)
	_ service.UserAccessRepository           = (*AccessControlRepository)(nil)
	_ service.AccessControlQueryRepository   = (*AccessControlRepository)(nil)
	_ service.RoleRepository                 = (*AccessControlRepository)(nil)
	_ service.AccessAssignmentRepository     = (*AccessControlRepository)(nil)
)
