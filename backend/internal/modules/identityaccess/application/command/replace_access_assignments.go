package command

import (
	"context"
	"errors"

	"webtranshome/internal/modules/identityaccess/domain/service"
)

var ErrInvalidAccessAssignment = errors.New("IDENTITY_INVALID_ACCESS_ASSIGNMENT")

type ReplaceAccessAssignments struct {
	repository service.AccessAssignmentRepository
}

func NewReplaceAccessAssignments(repository service.AccessAssignmentRepository) *ReplaceAccessAssignments {
	return &ReplaceAccessAssignments{repository: repository}
}

func (c *ReplaceAccessAssignments) UserRoles(ctx context.Context, userID int64, roleIDs []int64) error {
	if userID <= 0 || roleIDs == nil || !validIDs(roleIDs) {
		return ErrInvalidAccessAssignment
	}

	return c.repository.ReplaceUserRoles(ctx, userID, roleIDs)
}

func (c *ReplaceAccessAssignments) RolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error {
	if roleID <= 0 || permissionIDs == nil || !validIDs(permissionIDs) {
		return ErrInvalidAccessAssignment
	}

	return c.repository.ReplaceRolePermissions(ctx, roleID, permissionIDs)
}

func validIDs(ids []int64) bool {
	seen := make(map[int64]struct{}, len(ids))
	for _, id := range ids {
		if id <= 0 {
			return false
		}
		if _, exists := seen[id]; exists {
			return false
		}
		seen[id] = struct{}{}
	}

	return true
}
