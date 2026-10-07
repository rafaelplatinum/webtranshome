package command

import (
	"context"
	"errors"
	"testing"

	"webtranshome/internal/modules/identityaccess/domain/role"
	"webtranshome/internal/modules/identityaccess/domain/service"
)

type accessControlRepositoryStub struct {
	userID              int64
	roleID              int64
	roleIDs             []int64
	permissionIDs       []int64
	userRoleCalls       int
	rolePermissionCalls int
	createdRole         role.Role
	updatedRole         role.Role
	roleCreateCalls     int
	roleUpdateCalls     int
	err                 error
}

func (r *accessControlRepositoryStub) GetCatalog(context.Context) (service.AccessControlCatalog, error) {
	return service.AccessControlCatalog{}, nil
}

func (r *accessControlRepositoryStub) GetUserAccess(context.Context, int64) (service.UserAccess, error) {
	return service.UserAccess{}, nil
}

func (r *accessControlRepositoryStub) FindUserByEmail(context.Context, string) (service.AccessControlUser, error) {
	return service.AccessControlUser{}, nil
}

func (r *accessControlRepositoryStub) GetUserRoleIDs(context.Context, int64) ([]int64, error) {
	return nil, nil
}

func (r *accessControlRepositoryStub) GetRolePermissionIDs(context.Context, int64) ([]int64, error) {
	return nil, nil
}

func (r *accessControlRepositoryStub) CreateRole(_ context.Context, item role.Role) (role.Role, error) {
	r.roleCreateCalls++
	r.createdRole = item
	item.ID = 19
	return item, nil
}

func (r *accessControlRepositoryStub) UpdateRole(_ context.Context, item role.Role) error {
	r.roleUpdateCalls++
	r.updatedRole = item
	return r.err
}

func (r *accessControlRepositoryStub) ReplaceUserRoles(_ context.Context, userID int64, roleIDs []int64) error {
	r.userRoleCalls++
	r.userID = userID
	r.roleIDs = roleIDs
	return r.err
}

func (r *accessControlRepositoryStub) ReplaceRolePermissions(_ context.Context, roleID int64, permissionIDs []int64) error {
	r.rolePermissionCalls++
	r.roleID = roleID
	r.permissionIDs = permissionIDs
	return r.err
}

func TestReplaceAccessAssignmentsUserRoles(t *testing.T) {
	repository := &accessControlRepositoryStub{}
	command := NewReplaceAccessAssignments(repository)

	if err := command.UserRoles(context.Background(), 11, []int64{2, 4}); err != nil {
		t.Fatalf("UserRoles() error = %v", err)
	}
	if repository.userID != 11 || len(repository.roleIDs) != 2 || repository.roleIDs[0] != 2 || repository.roleIDs[1] != 4 {
		t.Fatalf("repository received userID=%d roleIDs=%v", repository.userID, repository.roleIDs)
	}
}

func TestReplaceAccessAssignmentsRolePermissions(t *testing.T) {
	repository := &accessControlRepositoryStub{}
	command := NewReplaceAccessAssignments(repository)

	if err := command.RolePermissions(context.Background(), 8, []int64{3, 5}); err != nil {
		t.Fatalf("RolePermissions() error = %v", err)
	}
	if repository.roleID != 8 || len(repository.permissionIDs) != 2 || repository.permissionIDs[0] != 3 || repository.permissionIDs[1] != 5 {
		t.Fatalf("repository received roleID=%d permissionIDs=%v", repository.roleID, repository.permissionIDs)
	}
}

func TestReplaceAccessAssignmentsRejectsInvalidInput(t *testing.T) {
	repository := &accessControlRepositoryStub{}
	command := NewReplaceAccessAssignments(repository)

	tests := []struct {
		name string
		run  func() error
	}{
		{name: "invalid user", run: func() error { return command.UserRoles(context.Background(), 0, []int64{1}) }},
		{name: "nil role list", run: func() error { return command.UserRoles(context.Background(), 1, nil) }},
		{name: "non-positive role", run: func() error { return command.UserRoles(context.Background(), 1, []int64{0}) }},
		{name: "duplicate roles", run: func() error { return command.UserRoles(context.Background(), 1, []int64{2, 2}) }},
		{name: "invalid role", run: func() error { return command.RolePermissions(context.Background(), 0, []int64{1}) }},
		{name: "nil permission list", run: func() error { return command.RolePermissions(context.Background(), 1, nil) }},
		{name: "non-positive permission", run: func() error { return command.RolePermissions(context.Background(), 1, []int64{-1}) }},
		{name: "duplicate permissions", run: func() error { return command.RolePermissions(context.Background(), 1, []int64{3, 3}) }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.run(); !errors.Is(err, ErrInvalidAccessAssignment) {
				t.Fatalf("error = %v, want %v", err, ErrInvalidAccessAssignment)
			}
			if repository.userRoleCalls != 0 || repository.rolePermissionCalls != 0 {
				t.Fatalf("repository calls = user roles %d, role permissions %d; want none", repository.userRoleCalls, repository.rolePermissionCalls)
			}
		})
	}
}

func TestReplaceAccessAssignmentsPropagatesRepositoryError(t *testing.T) {
	repositoryFailure := errors.New("repository failure")
	repository := &accessControlRepositoryStub{err: repositoryFailure}
	command := NewReplaceAccessAssignments(repository)

	if err := command.UserRoles(context.Background(), 11, []int64{}); !errors.Is(err, repositoryFailure) {
		t.Fatalf("UserRoles() error = %v, want %v", err, repositoryFailure)
	}
	if err := command.RolePermissions(context.Background(), 8, []int64{}); !errors.Is(err, repositoryFailure) {
		t.Fatalf("RolePermissions() error = %v, want %v", err, repositoryFailure)
	}
}

var _ service.AccessAssignmentRepository = (*accessControlRepositoryStub)(nil)
