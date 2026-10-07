package command

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestManageRolesCreateNormalizesFields(t *testing.T) {
	repository := &accessControlRepositoryStub{}
	management := NewManageRoles(repository)

	created, err := management.Create(context.Background(), CreateRoleCommand{
		Code:        " admin_reports ",
		Name:        "  Reports Admin  ",
		Description: "  Reads reporting data  ",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if created.ID != 19 || !created.IsActive {
		t.Fatalf("Create() result = %+v", created)
	}
	if repository.createdRole.Code != "ADMIN_REPORTS" ||
		repository.createdRole.Name != "Reports Admin" ||
		repository.createdRole.Description == nil ||
		*repository.createdRole.Description != "Reads reporting data" {
		t.Fatalf("repository received role = %+v", repository.createdRole)
	}
}

func TestManageRolesUsesCharacterLimitsForUnicodeNames(t *testing.T) {
	repository := &accessControlRepositoryStub{}
	management := NewManageRoles(repository)

	if _, err := management.Create(context.Background(), CreateRoleCommand{
		Code: "UNICODE_NAME",
		Name: strings.Repeat("é", 60),
	}); err != nil {
		t.Fatalf("Create() rejected a name below the database character limit: %v", err)
	}
}

func TestManageRolesCreateRejectsInvalidFields(t *testing.T) {
	repository := &accessControlRepositoryStub{}
	management := NewManageRoles(repository)
	tests := []CreateRoleCommand{
		{Code: " ", Name: "Valid name"},
		{Code: "1ADMIN", Name: "Valid name"},
		{Code: "ADMIN-REPORTS", Name: "Valid name"},
		{Code: "VALID", Name: " "},
		{Code: string(make([]byte, 51)), Name: "Valid name"},
		{Code: strings.Repeat("A", 51), Name: "Valid name"},
		{Code: "VALID", Name: string(make([]byte, 101))},
	}

	for _, input := range tests {
		if _, err := management.Create(context.Background(), input); !errors.Is(err, ErrInvalidRole) {
			t.Fatalf("Create(%+v) error = %v, want %v", input, err, ErrInvalidRole)
		}
	}
	if repository.roleCreateCalls != 0 {
		t.Fatalf("repository create calls = %d, want 0", repository.roleCreateCalls)
	}
}

func TestManageRolesUpdateRejectsInvalidFields(t *testing.T) {
	repository := &accessControlRepositoryStub{}
	management := NewManageRoles(repository)
	tests := []UpdateRoleCommand{
		{ID: 0, Name: "Valid name"},
		{ID: 3, Name: " "},
		{ID: 3, Name: string(make([]byte, 101))},
	}

	for _, input := range tests {
		if err := management.Update(context.Background(), input); !errors.Is(err, ErrInvalidRole) {
			t.Fatalf("Update(%+v) error = %v, want %v", input, err, ErrInvalidRole)
		}
	}
	if repository.roleUpdateCalls != 0 {
		t.Fatalf("repository update calls = %d, want 0", repository.roleUpdateCalls)
	}
}
