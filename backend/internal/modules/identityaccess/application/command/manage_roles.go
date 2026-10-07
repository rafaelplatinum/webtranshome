package command

import (
	"context"
	"strings"
	"unicode/utf8"

	"webtranshome/internal/modules/identityaccess/domain/role"
	"webtranshome/internal/modules/identityaccess/domain/service"
)

var ErrInvalidRole = role.ErrInvalidRole

type ManageRoles struct {
	repository service.RoleRepository
}

func NewManageRoles(repository service.RoleRepository) *ManageRoles {
	return &ManageRoles{repository: repository}
}

type CreateRoleCommand struct {
	Code        string
	Name        string
	Description string
}

type UpdateRoleCommand struct {
	ID          int64
	Name        string
	Description string
	IsActive    bool
}

func (c *ManageRoles) Create(ctx context.Context, input CreateRoleCommand) (role.Role, error) {
	newRole, err := role.New(input.Code, input.Name, input.Description)
	if err != nil {
		return role.Role{}, ErrInvalidRole
	}

	return c.repository.CreateRole(ctx, newRole)
}

func (c *ManageRoles) Update(ctx context.Context, input UpdateRoleCommand) error {
	name := strings.TrimSpace(input.Name)
	if input.ID <= 0 || name == "" || utf8.RuneCountInString(name) > 100 {
		return ErrInvalidRole
	}

	return c.repository.UpdateRole(ctx, role.Role{
		ID:          input.ID,
		Name:        name,
		Description: optionalDescription(input.Description),
		IsActive:    input.IsActive,
	})
}

func optionalDescription(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}
