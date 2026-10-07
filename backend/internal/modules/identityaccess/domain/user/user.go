package user

import "webtranshome/internal/modules/identityaccess/domain/role"

const (
	RoleSuperAdmin      = role.SuperAdminCode
	RoleAdminCatalog    = role.AdminCatalog
	RoleAdminMembership = role.AdminMember
	RoleAdminContent    = role.AdminContent
)

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	IsActive     bool
	HasStaffRole bool
}

func StaffRoleCodes() []string {
	return []string{
		RoleSuperAdmin,
		RoleAdminCatalog,
		RoleAdminMembership,
		RoleAdminContent,
	}
}
