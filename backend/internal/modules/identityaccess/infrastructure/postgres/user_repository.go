package postgres

import (
	"context"
	"errors"

	"webtranshome/internal/model"
	"webtranshome/internal/modules/identityaccess/domain/user"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UserRepository struct {
	connection sqlx.SqlConn
	users      model.UsersModel
}

func NewUserRepository(connection sqlx.SqlConn, users model.UsersModel) *UserRepository {
	return &UserRepository{
		connection: connection,
		users:      users,
	}
}

func (r *UserRepository) FindForAuthentication(ctx context.Context, email string) (user.User, error) {
	record, err := r.users.FindOneByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return user.User{}, user.ErrNotFound
		}
		return user.User{}, err
	}

	var roles []struct {
		Code string `db:"code"`
	}
	const query = `
		SELECT DISTINCT r.code
		FROM crm_schema.user_roles ur
		JOIN crm_schema.roles r ON r.id = ur.role_id
		WHERE ur.user_id = $1 AND r.is_active = TRUE`
	if err := r.connection.QueryRowsCtx(ctx, &roles, query, record.Id); err != nil {
		return user.User{}, err
	}

	roleCodes := make([]string, 0, len(roles))
	for _, role := range roles {
		roleCodes = append(roleCodes, role.Code)
	}

	passwordHash := ""
	if record.PasswordHash.Valid {
		passwordHash = record.PasswordHash.String
	}

	return user.User{
		ID:           record.Id,
		Email:        record.Email,
		PasswordHash: passwordHash,
		IsActive:     record.IsActive,
		RoleCodes:    roleCodes,
	}, nil
}

var _ user.Repository = (*UserRepository)(nil)
