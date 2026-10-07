package role

import (
	"errors"
	"strings"
	"unicode/utf8"
)

const (
	SuperAdminCode = "SUPER_ADMIN"
	AdminCatalog   = "ADMIN_KATALOG"
	AdminMember    = "ADMIN_MEMBERSHIP"
	AdminContent   = "ADMIN_KONTEN"
)

var (
	ErrInvalidRole   = errors.New("IDENTITY_INVALID_ROLE")
	ErrProtectedRole = errors.New("IDENTITY_PROTECTED_ROLE")
)

type Role struct {
	ID          int64
	Code        string
	Name        string
	Description *string
	IsActive    bool
}

func New(code, name, description string) (Role, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	name = strings.TrimSpace(name)
	if !validCode(code) || name == "" || utf8.RuneCountInString(name) > 100 {
		return Role{}, ErrInvalidRole
	}
	description = strings.TrimSpace(description)
	var descriptionValue *string
	if description != "" {
		descriptionValue = &description
	}
	return Role{
		Code:        code,
		Name:        name,
		Description: descriptionValue,
		IsActive:    true,
	}, nil
}

func validCode(code string) bool {
	if len(code) == 0 || len(code) > 50 || code[0] < 'A' || code[0] > 'Z' {
		return false
	}
	for i := 1; i < len(code); i++ {
		char := code[i]
		if (char < 'A' || char > 'Z') && (char < '0' || char > '9') && char != '_' {
			return false
		}
	}
	return true
}

func (r Role) IsBootstrap() bool {
	return r.Code == SuperAdminCode
}

func (r Role) ValidateActivation(active bool) error {
	if r.IsBootstrap() && !active {
		return ErrProtectedRole
	}
	return nil
}

func (r Role) ValidatePermissions(permissionIDs []int64, requiredPermissionID int64) error {
	if !r.IsBootstrap() {
		return nil
	}
	for _, permissionID := range permissionIDs {
		if permissionID == requiredPermissionID {
			return nil
		}
	}
	return ErrProtectedRole
}
