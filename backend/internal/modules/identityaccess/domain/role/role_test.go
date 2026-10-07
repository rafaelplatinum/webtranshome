package role

import (
	"errors"
	"strings"
	"testing"
)

func TestNewNormalizesAndValidatesRoleFields(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		roleName    string
		description string
		wantCode    string
		wantErr     error
	}{
		{
			name:        "normalizes valid role",
			code:        " admin_reports ",
			roleName:    " Reports Admin ",
			description: " Reporting ",
			wantCode:    "ADMIN_REPORTS",
		},
		{name: "rejects invalid leading character", code: "1ADMIN", roleName: "Admin", wantErr: ErrInvalidRole},
		{name: "rejects invalid code character", code: "ADMIN-REPORTS", roleName: "Admin", wantErr: ErrInvalidRole},
		{name: "rejects code over limit", code: strings.Repeat("A", 51), roleName: "Admin", wantErr: ErrInvalidRole},
		{name: "rejects empty name", code: "ADMIN", roleName: " ", wantErr: ErrInvalidRole},
		{name: "rejects name over limit", code: "ADMIN", roleName: strings.Repeat("é", 101), wantErr: ErrInvalidRole},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			created, err := New(test.code, test.roleName, test.description)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("New() error = %v, want %v", err, test.wantErr)
			}
			if test.wantErr != nil {
				return
			}
			if created.Code != test.wantCode || created.Name != "Reports Admin" || !created.IsActive {
				t.Fatalf("New() role = %+v", created)
			}
			if created.Description == nil || *created.Description != "Reporting" {
				t.Fatalf("New() description = %v, want Reporting", created.Description)
			}
		})
	}
}

func TestRoleValidateActivation(t *testing.T) {
	tests := []struct {
		name   string
		role   Role
		active bool
		want   error
	}{
		{name: "bootstrap stays active", role: Role{Code: SuperAdminCode}, active: true},
		{name: "bootstrap cannot be disabled", role: Role{Code: SuperAdminCode}, want: ErrProtectedRole},
		{name: "regular role can be disabled", role: Role{Code: AdminCatalog}, want: nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.role.ValidateActivation(test.active)
			if !errors.Is(err, test.want) {
				t.Fatalf("ValidateActivation() error = %v, want %v", err, test.want)
			}
		})
	}
}

func TestRoleValidatePermissions(t *testing.T) {
	tests := []struct {
		name               string
		role               Role
		permissionIDs      []int64
		requiredPermission int64
		want               error
	}{
		{
			name:               "bootstrap keeps management permission",
			role:               Role{Code: SuperAdminCode},
			permissionIDs:      []int64{3, 7},
			requiredPermission: 7,
		},
		{
			name:               "bootstrap cannot lose management permission",
			role:               Role{Code: SuperAdminCode},
			permissionIDs:      []int64{3},
			requiredPermission: 7,
			want:               ErrProtectedRole,
		},
		{
			name:               "regular role has no bootstrap restriction",
			role:               Role{Code: AdminCatalog},
			permissionIDs:      []int64{3},
			requiredPermission: 7,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.role.ValidatePermissions(test.permissionIDs, test.requiredPermission)
			if !errors.Is(err, test.want) {
				t.Fatalf("ValidatePermissions() error = %v, want %v", err, test.want)
			}
		})
	}
}
