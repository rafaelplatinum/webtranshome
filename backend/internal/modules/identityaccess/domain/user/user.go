package user

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
