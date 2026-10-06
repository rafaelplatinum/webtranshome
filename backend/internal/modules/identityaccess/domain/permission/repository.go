package permission

import "context"

type Repository interface {
	HasUserPermission(ctx context.Context, userID int64, permissionCode string) (bool, error)
}
