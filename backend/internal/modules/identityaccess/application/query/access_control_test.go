package query

import (
	"context"
	"errors"
	"testing"

	"webtranshome/internal/modules/identityaccess/domain/menu"
	"webtranshome/internal/modules/identityaccess/domain/permission"
	"webtranshome/internal/modules/identityaccess/domain/role"
	"webtranshome/internal/modules/identityaccess/domain/service"
)

type accessControlRepositoryStub struct {
	userRoleCalls       int
	rolePermissionCalls int
	userAccessCalls     int
	userSearchCalls     int
	searchEmail         string
	userAccess          service.UserAccess
	foundUser           service.AccessControlUser
	err                 error
}

func (r *accessControlRepositoryStub) GetCatalog(context.Context) (service.AccessControlCatalog, error) {
	return service.AccessControlCatalog{}, r.err
}

func (r *accessControlRepositoryStub) GetUserAccess(context.Context, int64) (service.UserAccess, error) {
	r.userAccessCalls++
	return r.userAccess, r.err
}

func (r *accessControlRepositoryStub) FindUserByEmail(_ context.Context, email string) (service.AccessControlUser, error) {
	r.userSearchCalls++
	r.searchEmail = email
	return r.foundUser, r.err
}

func (r *accessControlRepositoryStub) GetUserRoleIDs(context.Context, int64) ([]int64, error) {
	r.userRoleCalls++
	return []int64{2}, r.err
}

func (r *accessControlRepositoryStub) GetRolePermissionIDs(context.Context, int64) ([]int64, error) {
	r.rolePermissionCalls++
	return []int64{3}, r.err
}

func (r *accessControlRepositoryStub) CreateRole(context.Context, role.Role) (role.Role, error) {
	return role.Role{}, nil
}

func (r *accessControlRepositoryStub) UpdateRole(context.Context, role.Role) error {
	return nil
}

func (r *accessControlRepositoryStub) ReplaceUserRoles(context.Context, int64, []int64) error {
	return nil
}

func (r *accessControlRepositoryStub) ReplaceRolePermissions(context.Context, int64, []int64) error {
	return nil
}

func TestAccessControlAssignmentQueries(t *testing.T) {
	repository := &accessControlRepositoryStub{}
	query := NewAccessControl(repository, repository, repository)

	roleIDs, err := query.UserRoleIDs(context.Background(), 7)
	if err != nil || len(roleIDs) != 1 || roleIDs[0] != 2 {
		t.Fatalf("UserRoleIDs() = %v, %v; want [2], nil", roleIDs, err)
	}
	permissionIDs, err := query.RolePermissionIDs(context.Background(), 9)
	if err != nil || len(permissionIDs) != 1 || permissionIDs[0] != 3 {
		t.Fatalf("RolePermissionIDs() = %v, %v; want [3], nil", permissionIDs, err)
	}
	if repository.userRoleCalls != 1 || repository.rolePermissionCalls != 1 {
		t.Fatalf("repository calls = user roles %d, role permissions %d; want 1 each", repository.userRoleCalls, repository.rolePermissionCalls)
	}
}

func TestAccessControlUserAccessMapsEffectiveMenusAndPermissions(t *testing.T) {
	parentID := int64(2)
	repository := &accessControlRepositoryStub{
		userAccess: service.UserAccess{
			Menus: []menu.Menu{
				{ID: 2, Name: "Parent", Code: "parent", IsActive: true},
				{ID: 3, ParentID: &parentID, Name: "Child", Code: "child", Route: stringPointer("/child"), IsActive: true},
			},
			Permissions: []permission.Permission{
				{ID: 5, MenuID: 3, Code: "child.read", Action: "read"},
			},
		},
	}
	query := NewAccessControl(repository, repository, repository)

	access, err := query.UserAccess(context.Background(), 7)
	if err != nil {
		t.Fatalf("UserAccess() error = %v", err)
	}
	if len(access.Menus) != 2 || access.Menus[1].ParentID != 2 {
		t.Fatalf("UserAccess() menus = %+v", access.Menus)
	}
	if len(access.PermissionCodes) != 1 || access.PermissionCodes[0] != "child.read" {
		t.Fatalf("UserAccess() permission codes = %v", access.PermissionCodes)
	}
	if repository.userAccessCalls != 1 {
		t.Fatalf("GetUserAccess() calls = %d, want 1", repository.userAccessCalls)
	}
}

func TestAccessControlFindUserByEmail(t *testing.T) {
	repository := &accessControlRepositoryStub{
		foundUser: service.AccessControlUser{
			ID:       12,
			Email:    "member@example.test",
			IsActive: true,
			RoleIDs:  []int64{4},
		},
	}
	query := NewAccessControl(repository, repository, repository)

	found, err := query.FindUserByEmail(context.Background(), " member@example.test ")
	if err != nil {
		t.Fatalf("FindUserByEmail() error = %v", err)
	}
	if found.ID != 12 || found.Email != "member@example.test" || !found.IsActive ||
		len(found.RoleIDs) != 1 || found.RoleIDs[0] != 4 {
		t.Fatalf("FindUserByEmail() = %+v", found)
	}
	if repository.userSearchCalls != 1 || repository.searchEmail != "member@example.test" {
		t.Fatalf("user search calls = %d email = %q, want 1 and trimmed email", repository.userSearchCalls, repository.searchEmail)
	}
}

func TestAccessControlQueriesRejectInvalidIDs(t *testing.T) {
	repository := &accessControlRepositoryStub{}
	query := NewAccessControl(repository, repository, repository)

	if _, err := query.UserAccess(context.Background(), 0); !errors.Is(err, ErrInvalidAccessControlQuery) {
		t.Fatalf("UserAccess() error = %v, want %v", err, ErrInvalidAccessControlQuery)
	}
	if _, err := query.UserRoleIDs(context.Background(), 0); !errors.Is(err, ErrInvalidAccessControlQuery) {
		t.Fatalf("UserRoleIDs() error = %v, want %v", err, ErrInvalidAccessControlQuery)
	}
	if _, err := query.RolePermissionIDs(context.Background(), -1); !errors.Is(err, ErrInvalidAccessControlQuery) {
		t.Fatalf("RolePermissionIDs() error = %v, want %v", err, ErrInvalidAccessControlQuery)
	}
	if _, err := query.FindUserByEmail(context.Background(), " "); !errors.Is(err, ErrInvalidAccessControlQuery) {
		t.Fatalf("FindUserByEmail() error = %v, want %v", err, ErrInvalidAccessControlQuery)
	}
	if repository.userAccessCalls != 0 || repository.userSearchCalls != 0 ||
		repository.userRoleCalls != 0 || repository.rolePermissionCalls != 0 {
		t.Fatalf("repository calls = user access %d, user search %d, user roles %d, role permissions %d; want none",
			repository.userAccessCalls, repository.userSearchCalls, repository.userRoleCalls, repository.rolePermissionCalls)
	}
}

func TestAccessControlQueriesPropagateRepositoryErrors(t *testing.T) {
	repositoryFailure := errors.New("repository failure")
	repository := &accessControlRepositoryStub{err: repositoryFailure}
	query := NewAccessControl(repository, repository, repository)

	if _, err := query.Catalog(context.Background()); !errors.Is(err, repositoryFailure) {
		t.Fatalf("Catalog() error = %v, want %v", err, repositoryFailure)
	}
	if _, err := query.UserAccess(context.Background(), 7); !errors.Is(err, repositoryFailure) {
		t.Fatalf("UserAccess() error = %v, want %v", err, repositoryFailure)
	}
	if _, err := query.UserRoleIDs(context.Background(), 7); !errors.Is(err, repositoryFailure) {
		t.Fatalf("UserRoleIDs() error = %v, want %v", err, repositoryFailure)
	}
	if _, err := query.RolePermissionIDs(context.Background(), 9); !errors.Is(err, repositoryFailure) {
		t.Fatalf("RolePermissionIDs() error = %v, want %v", err, repositoryFailure)
	}
	if _, err := query.FindUserByEmail(context.Background(), "member@example.test"); !errors.Is(err, repositoryFailure) {
		t.Fatalf("FindUserByEmail() error = %v, want %v", err, repositoryFailure)
	}
}

var (
	_ service.AccessControlCatalogRepository = (*accessControlRepositoryStub)(nil)
	_ service.UserAccessRepository           = (*accessControlRepositoryStub)(nil)
	_ service.AccessControlQueryRepository   = (*accessControlRepositoryStub)(nil)
)

func stringPointer(value string) *string {
	return &value
}
