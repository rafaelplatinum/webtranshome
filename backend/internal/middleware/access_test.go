package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"webtranshome/internal/modules/identityaccess/application/port"
)

type tokenVerifierStub struct {
	userID int64
	err    error
	token  string
}

func (v *tokenVerifierStub) Verify(_ context.Context, token string) (int64, error) {
	v.token = token
	return v.userID, v.err
}

type permissionCheckerStub struct {
	userID  int64
	code    string
	granted bool
	err     error
	called  bool
}

func (c *permissionCheckerStub) Execute(_ context.Context, userID int64, code string) (bool, error) {
	c.called = true
	c.userID = userID
	c.code = code
	return c.granted, c.err
}

func TestAuthenticate(t *testing.T) {
	tests := []struct {
		name       string
		header     string
		userID     int64
		verifyErr  error
		wantStatus int
		wantNext   bool
		wantToken  string
	}{
		{name: "valid bearer", header: "Bearer access-token", userID: 42, wantStatus: http.StatusNoContent, wantNext: true, wantToken: "access-token"},
		{name: "case insensitive bearer", header: "bearer access-token", userID: 42, wantStatus: http.StatusNoContent, wantNext: true, wantToken: "access-token"},
		{name: "missing bearer", wantStatus: http.StatusUnauthorized},
		{name: "invalid token", header: "Bearer invalid", verifyErr: port.ErrInvalidAccessToken, wantStatus: http.StatusUnauthorized, wantToken: "invalid"},
		{name: "verifier failure", header: "Bearer token", verifyErr: errors.New("verifier failure"), wantStatus: http.StatusInternalServerError, wantToken: "token"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			verifier := &tokenVerifierStub{userID: test.userID, err: test.verifyErr}
			nextCalled := false
			next := Authenticate(verifier)(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				if userID, ok := UserIDFromContext(r.Context()); !ok || userID != 42 {
					t.Fatalf("UserIDFromContext() = %d, %t; want 42, true", userID, ok)
				}
				w.WriteHeader(http.StatusNoContent)
			})
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			request.Header.Set("Authorization", test.header)
			recorder := httptest.NewRecorder()

			next(recorder, request)

			if recorder.Code != test.wantStatus || nextCalled != test.wantNext {
				t.Fatalf("status=%d nextCalled=%t; want status=%d nextCalled=%t", recorder.Code, nextCalled, test.wantStatus, test.wantNext)
			}
			if verifier.token != test.wantToken {
				t.Fatalf("verified token = %q, want %q", verifier.token, test.wantToken)
			}
		})
	}
}

func TestRequirePermission(t *testing.T) {
	tests := []struct {
		name            string
		withPrincipal   bool
		granted         bool
		checkErr        error
		wantStatus      int
		wantCheckerCall bool
	}{
		{name: "granted", withPrincipal: true, granted: true, wantStatus: http.StatusNoContent, wantCheckerCall: true},
		{name: "denied", withPrincipal: true, wantStatus: http.StatusForbidden, wantCheckerCall: true},
		{name: "no principal", wantStatus: http.StatusUnauthorized},
		{name: "checker failure", withPrincipal: true, checkErr: errors.New("repository failure"), wantStatus: http.StatusInternalServerError, wantCheckerCall: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			checker := &permissionCheckerStub{granted: test.granted, err: test.checkErr}
			nextCalled := false
			next := RequirePermission(checker, "catalog.product.read")(func(w http.ResponseWriter, _ *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusNoContent)
			})
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			if test.withPrincipal {
				request = request.WithContext(context.WithValue(request.Context(), principalContextKey{}, int64(42)))
			}
			recorder := httptest.NewRecorder()

			next(recorder, request)

			if recorder.Code != test.wantStatus || nextCalled != (test.wantStatus == http.StatusNoContent) {
				t.Fatalf("status=%d nextCalled=%t; want status=%d", recorder.Code, nextCalled, test.wantStatus)
			}
			if checker.called != test.wantCheckerCall {
				t.Fatalf("checker called=%t, want %t", checker.called, test.wantCheckerCall)
			}
			if test.wantCheckerCall && (checker.userID != 42 || checker.code != "catalog.product.read") {
				t.Fatalf("checker received userID=%d code=%q", checker.userID, checker.code)
			}
		})
	}
}
