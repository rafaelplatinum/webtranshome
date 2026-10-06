package query

import (
	"context"
	"errors"
	"testing"
)

type permissionRepositoryStub struct {
	called  bool
	granted bool
	err     error
	userID  int64
	code    string
}

func (r *permissionRepositoryStub) HasUserPermission(_ context.Context, userID int64, code string) (bool, error) {
	r.called = true
	r.userID = userID
	r.code = code
	return r.granted, r.err
}

func TestCheckPermissionExecute(t *testing.T) {
	repositoryFailure := errors.New("repository failure")
	tests := []struct {
		name               string
		userID             int64
		code               string
		granted            bool
		repository         error
		want               bool
		wantErr            error
		wantRepositoryCode string
		wantRepositoryCall bool
	}{
		{name: "granted", userID: 7, code: "catalog.read", granted: true, want: true, wantRepositoryCode: "catalog.read", wantRepositoryCall: true},
		{name: "denied", userID: 7, code: "catalog.read", want: false, wantRepositoryCode: "catalog.read", wantRepositoryCall: true},
		{name: "trim permission code", userID: 7, code: " catalog.read ", granted: true, want: true, wantRepositoryCode: "catalog.read", wantRepositoryCall: true},
		{name: "invalid user", userID: 0, code: "catalog.read", wantErr: ErrInvalidPermissionCheck},
		{name: "empty code", userID: 7, code: " ", wantErr: ErrInvalidPermissionCheck},
		{name: "repository error", userID: 7, code: "catalog.read", repository: repositoryFailure, wantErr: repositoryFailure, wantRepositoryCode: "catalog.read", wantRepositoryCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := &permissionRepositoryStub{granted: test.granted, err: test.repository}
			query := NewCheckPermission(repository)

			got, err := query.Execute(context.Background(), test.userID, test.code)
			if got != test.want {
				t.Fatalf("Execute() = %t, want %t", got, test.want)
			}
			if test.wantErr != nil && !errors.Is(err, test.wantErr) {
				t.Fatalf("Execute() error = %v, want %v", err, test.wantErr)
			}
			if test.wantErr == nil && err != nil {
				t.Fatalf("Execute() error = %v, want nil", err)
			}
			if repository.called != test.wantRepositoryCall {
				t.Fatalf("repository called = %t, want %t", repository.called, test.wantRepositoryCall)
			}
			if test.wantRepositoryCall && (repository.userID != test.userID || repository.code != test.wantRepositoryCode) {
				t.Fatalf("repository received userID=%d code=%q", repository.userID, repository.code)
			}
		})
	}
}
