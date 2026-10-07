package types

type AccessControlRole struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IsActive    bool   `json:"isActive"`
}

type AccessControlMenu struct {
	ID        int64  `json:"id"`
	ParentID  int64  `json:"parentId,omitempty"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	Route     string `json:"route,omitempty"`
	Icon      string `json:"icon,omitempty"`
	SortOrder int64  `json:"sortOrder"`
	IsActive  bool   `json:"isActive"`
}

type AccessControlPermission struct {
	ID          int64  `json:"id"`
	MenuID      int64  `json:"menuId"`
	Code        string `json:"code"`
	Action      string `json:"action"`
	Description string `json:"description,omitempty"`
}

type AccessControlCatalogResponse struct {
	Roles       []AccessControlRole       `json:"roles"`
	Menus       []AccessControlMenu       `json:"menus"`
	Permissions []AccessControlPermission `json:"permissions"`
}

type ReplaceUserRolesRequest struct {
	RoleIDs []int64 `json:"roleIds"`
}

type ReplaceRolePermissionsRequest struct {
	PermissionIDs []int64 `json:"permissionIds"`
}

type CreateRoleRequest struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type UpdateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IsActive    bool   `json:"isActive"`
}

type AccessControlRoleAssignmentsResponse struct {
	RoleIDs []int64 `json:"roleIds"`
}

type AccessControlPermissionAssignmentsResponse struct {
	PermissionIDs []int64 `json:"permissionIds"`
}

type AccessControlUserResponse struct {
	ID       int64   `json:"id"`
	Email    string  `json:"email"`
	IsActive bool    `json:"isActive"`
	RoleIDs  []int64 `json:"roleIds"`
}

type SearchAccessControlUsersRequest struct {
	Email string `form:"email"`
}

type AccessControlMutationResponse struct {
	Success bool `json:"success"`
}

type UserAccessResponse struct {
	Menus           []AccessControlMenu `json:"menus"`
	PermissionCodes []string            `json:"permissionCodes"`
}
