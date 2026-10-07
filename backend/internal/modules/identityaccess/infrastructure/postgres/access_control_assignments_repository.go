package postgres

import (
	"context"
	"errors"

	"webtranshome/internal/modules/identityaccess/domain/permission"
	"webtranshome/internal/modules/identityaccess/domain/role"
	"webtranshome/internal/modules/identityaccess/domain/service"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func (r *AccessControlRepository) ReplaceUserRoles(ctx context.Context, userID int64, roleIDs []int64) error {
	return r.connection.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var active bool
		if err := session.QueryRowCtx(ctx, &active,
			`SELECT is_active FROM crm_schema.users WHERE id = $1`, userID); err != nil {
			if errors.Is(err, sqlx.ErrNotFound) {
				return service.ErrUserNotFound
			}
			return err
		}
		if !active {
			return service.ErrUserNotFound
		}

		for _, roleID := range roleIDs {
			var roleActive bool
			if err := session.QueryRowCtx(ctx, &roleActive,
				`SELECT is_active FROM crm_schema.roles WHERE id = $1`, roleID); err != nil {
				if errors.Is(err, sqlx.ErrNotFound) {
					return service.ErrRoleNotFound
				}
				return err
			}
			if !roleActive {
				return service.ErrRoleNotFound
			}
		}

		if _, err := session.ExecCtx(ctx,
			`DELETE FROM crm_schema.user_roles WHERE user_id = $1`, userID); err != nil {
			return err
		}
		for _, roleID := range roleIDs {
			if _, err := session.ExecCtx(ctx,
				`INSERT INTO crm_schema.user_roles (user_id, role_id) VALUES ($1, $2)`,
				userID, roleID); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *AccessControlRepository) ReplaceRolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error {
	return r.connection.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var roleRecord struct {
			Code     string `db:"code"`
			IsActive bool   `db:"is_active"`
		}
		if err := session.QueryRowCtx(ctx, &roleRecord,
			`SELECT code, is_active FROM crm_schema.roles WHERE id = $1`, roleID); err != nil {
			if errors.Is(err, sqlx.ErrNotFound) {
				return service.ErrRoleNotFound
			}
			return err
		}
		if !roleRecord.IsActive {
			return service.ErrRoleNotFound
		}

		for _, permissionID := range permissionIDs {
			var exists bool
			if err := session.QueryRowCtx(ctx, &exists,
				`SELECT EXISTS (
					SELECT 1
					FROM crm_schema.permissions p
					JOIN crm_schema.menus m ON m.id = p.menu_id
					WHERE p.id = $1 AND m.is_active = TRUE
				)`, permissionID); err != nil {
				return err
			}
			if !exists {
				return service.ErrPermissionNotFound
			}
		}

		protectedRole := role.Role{Code: roleRecord.Code}
		if protectedRole.IsBootstrap() {
			var managePermissionID int64
			if err := session.QueryRowCtx(ctx, &managePermissionID,
				`SELECT id FROM crm_schema.permissions WHERE code = $1`,
				permission.AccessControlManageCode); err != nil {
				if errors.Is(err, sqlx.ErrNotFound) {
					return service.ErrPermissionNotFound
				}
				return err
			}
			if err := protectedRole.ValidatePermissions(permissionIDs, managePermissionID); err != nil {
				return err
			}
		}

		if _, err := session.ExecCtx(ctx,
			`DELETE FROM crm_schema.role_permissions WHERE role_id = $1`, roleID); err != nil {
			return err
		}
		for _, permissionID := range permissionIDs {
			if _, err := session.ExecCtx(ctx,
				`INSERT INTO crm_schema.role_permissions (role_id, permission_id) VALUES ($1, $2)`,
				roleID, permissionID); err != nil {
				return err
			}
		}

		return nil
	})
}
