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

	passwordHash := ""
	if record.PasswordHash.Valid {
		passwordHash = record.PasswordHash.String
	}

	roleCodes := user.StaffRoleCodes()
	query := `SELECT EXISTS (
		SELECT 1
		FROM crm_schema.user_roles ur
		JOIN crm_schema.roles r ON r.id = ur.role_id
		WHERE ur.user_id = $1 AND r.is_active = TRUE AND r.code IN ($2, $3, $4, $5)
	)`
	var hasStaffRole bool
	if err := r.connection.QueryRowCtx(
		ctx,
		&hasStaffRole,
		query,
		record.Id,
		roleCodes[0],
		roleCodes[1],
		roleCodes[2],
		roleCodes[3],
	); err != nil {
		return user.User{}, err
	}

	return user.User{
		ID:           record.Id,
		Email:        record.Email,
		PasswordHash: passwordHash,
		IsActive:     record.IsActive,
		HasStaffRole: hasStaffRole,
	}, nil
}

var _ user.Repository = (*UserRepository)(nil)
