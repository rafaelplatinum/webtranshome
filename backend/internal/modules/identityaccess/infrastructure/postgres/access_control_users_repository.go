package postgres

import (
	"context"
	"errors"

	"webtranshome/internal/modules/identityaccess/domain/service"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func (r *AccessControlRepository) FindUserByEmail(ctx context.Context, email string) (service.AccessControlUser, error) {
	var found struct {
		ID       int64  `db:"id"`
		Email    string `db:"email"`
		IsActive bool   `db:"is_active"`
	}
	if err := r.connection.QueryRowCtx(ctx, &found,
		`SELECT id, email, is_active
		FROM crm_schema.users
		WHERE email = $1`,
		email); err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return service.AccessControlUser{}, service.ErrUserNotFound
		}
		return service.AccessControlUser{}, err
	}

	roleIDs, err := r.GetUserRoleIDs(ctx, found.ID)
	if err != nil {
		return service.AccessControlUser{}, err
	}

	return service.AccessControlUser{
		ID:       found.ID,
		Email:    found.Email,
		IsActive: found.IsActive,
		RoleIDs:  roleIDs,
	}, nil
}

func (r *AccessControlRepository) GetUserRoleIDs(ctx context.Context, userID int64) ([]int64, error) {
	var exists bool
	if err := r.connection.QueryRowCtx(ctx, &exists,
		`SELECT EXISTS (SELECT 1 FROM crm_schema.users WHERE id = $1)`, userID); err != nil {
		return nil, err
	}
	if !exists {
		return nil, service.ErrUserNotFound
	}

	var rows []struct {
		ID int64 `db:"role_id"`
	}
	if err := r.connection.QueryRowsCtx(ctx, &rows,
		`SELECT role_id FROM crm_schema.user_roles WHERE user_id = $1 ORDER BY role_id`, userID); err != nil {
		return nil, err
	}

	ids := make([]int64, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ID)
	}
	return ids, nil
}

func (r *AccessControlRepository) GetRolePermissionIDs(ctx context.Context, roleID int64) ([]int64, error) {
	var exists bool
	if err := r.connection.QueryRowCtx(ctx, &exists,
		`SELECT EXISTS (SELECT 1 FROM crm_schema.roles WHERE id = $1)`, roleID); err != nil {
		return nil, err
	}
	if !exists {
		return nil, service.ErrRoleNotFound
	}

	var rows []struct {
		ID int64 `db:"permission_id"`
	}
	if err := r.connection.QueryRowsCtx(ctx, &rows,
		`SELECT permission_id FROM crm_schema.role_permissions WHERE role_id = $1 ORDER BY permission_id`, roleID); err != nil {
		return nil, err
	}

	ids := make([]int64, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ID)
	}
	return ids, nil
}
