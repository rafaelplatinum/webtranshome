package postgres

import (
	"context"
	"errors"

	"webtranshome/internal/modules/identityaccess/domain/role"
	"webtranshome/internal/modules/identityaccess/domain/service"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func (r *AccessControlRepository) CreateRole(ctx context.Context, item role.Role) (role.Role, error) {
	var id int64
	err := r.connection.QueryRowCtx(ctx, &id,
		`INSERT INTO crm_schema.roles (code, name, description, is_active)
		VALUES ($1, $2, $3, TRUE)
		ON CONFLICT (code) DO NOTHING
		RETURNING id`,
		item.Code, item.Name, item.Description)
	if errors.Is(err, sqlx.ErrNotFound) {
		return role.Role{}, service.ErrRoleCodeConflict
	}
	if err != nil {
		return role.Role{}, err
	}

	item.ID = id
	item.IsActive = true
	return item, nil
}

func (r *AccessControlRepository) UpdateRole(ctx context.Context, item role.Role) error {
	return r.connection.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var code string
		if err := session.QueryRowCtx(ctx, &code,
			`SELECT code FROM crm_schema.roles WHERE id = $1 FOR UPDATE`, item.ID); err != nil {
			if errors.Is(err, sqlx.ErrNotFound) {
				return service.ErrRoleNotFound
			}
			return err
		}
		if err := (role.Role{Code: code}).ValidateActivation(item.IsActive); err != nil {
			return err
		}

		_, err := session.ExecCtx(ctx,
			`UPDATE crm_schema.roles
			SET name = $2, description = $3, is_active = $4
			WHERE id = $1`,
			item.ID, item.Name, item.Description, item.IsActive)
		return err
	})
}
