package user

import "strings"

const (
	RoleSuperAdmin      = "SUPER_ADMIN"
	RoleAdminCatalog    = "ADMIN_KATALOG"
	RoleAdminMembership = "ADMIN_MEMBERSHIP"
	RoleAdminContent    = "ADMIN_KONTEN"
)

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	IsActive     bool
	RoleCodes    []string
}

func (u User) IsAdmin() bool {
	for _, roleCode := range u.RoleCodes {
		if roleCode == RoleSuperAdmin ||
			roleCode == RoleAdminCatalog ||
			roleCode == RoleAdminMembership ||
			roleCode == RoleAdminContent {
			return true
		}

		if strings.HasPrefix(roleCode, "ADMIN_") {
			return true
		}
	}

	return false
}
