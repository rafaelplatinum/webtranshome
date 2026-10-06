package query

import (
	"context"
	"errors"
	"strings"

	"webtranshome/internal/modules/identityaccess/domain/permission"
)

var ErrInvalidPermissionCheck = errors.New("IDENTITY_INVALID_PERMISSION_CHECK")

type CheckPermission struct {
	permissions permission.Repository
}

func NewCheckPermission(permissions permission.Repository) *CheckPermission {
	return &CheckPermission{permissions: permissions}
}

func (q *CheckPermission) Execute(ctx context.Context, userID int64, permissionCode string) (bool, error) {
	permissionCode = strings.TrimSpace(permissionCode)
	if userID <= 0 || permissionCode == "" {
		return false, ErrInvalidPermissionCheck
	}

	return q.permissions.HasUserPermission(ctx, userID, permissionCode)
}
