package postgres

import (
	"context"

	"webtranshome/internal/modules/identityaccess/domain/permission"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type PermissionRepository struct {
	connection sqlx.SqlConn
}

func NewPermissionRepository(connection sqlx.SqlConn) *PermissionRepository {
	return &PermissionRepository{connection: connection}
}

func (r *PermissionRepository) HasUserPermission(ctx context.Context, userID int64, permissionCode string) (bool, error) {
	const query = `SELECT EXISTS (
		SELECT 1
		FROM crm_schema.user_roles ur
		JOIN crm_schema.users u ON u.id = ur.user_id
		JOIN crm_schema.roles r ON r.id = ur.role_id
		JOIN crm_schema.role_permissions rp ON rp.role_id = r.id
		JOIN crm_schema.permissions p ON p.id = rp.permission_id
		JOIN crm_schema.menus m ON m.id = p.menu_id
		WHERE ur.user_id = $1
			AND p.code = $2
			AND u.is_active = TRUE
			AND r.is_active = TRUE
			AND m.is_active = TRUE
	)`

	var hasPermission bool
	if err := r.connection.QueryRowCtx(ctx, &hasPermission, query, userID, permissionCode); err != nil {
		return false, err
	}

	return hasPermission, nil
}

var _ permission.Repository = (*PermissionRepository)(nil)
