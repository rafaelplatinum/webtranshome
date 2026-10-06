package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"webtranshome/internal/modules/identityaccess/application/command"
	"webtranshome/internal/modules/identityaccess/application/port"
	"webtranshome/internal/modules/identityaccess/domain/user"
	"webtranshome/internal/svc"
	"webtranshome/internal/types"
)

type loginUserRepositoryStub struct {
	account user.User
	err     error
}

func (s loginUserRepositoryStub) FindForAuthentication(context.Context, string) (user.User, error) {
	return s.account, s.err
}

type loginPasswordVerifierStub struct {
	err error
}

func (s loginPasswordVerifierStub) Verify(string, string) error {
	return s.err
}

type loginTokenIssuerStub struct {
	token string
	err   error
}

func (s loginTokenIssuerStub) Issue(context.Context, int64) (string, error) {
	return s.token, s.err
}

func TestLoginHandlerReturnsTokenForValidAdmin(t *testing.T) {
	handler := newLoginHandler(user.User{
		ID:           42,
		Email:        "admin@example.test",
		PasswordHash: "stored-hash",
		IsActive:     true,
		HasStaffRole: true,
	}, nil, "signed-token", nil)

	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(
		`{"email":"admin@example.test","password":"submitted-password"}`,
	))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}

	var result types.LoginResponse
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if result.AccessToken != "signed-token" {
		t.Fatalf("access token = %q, want token from issuer", result.AccessToken)
	}
}

func TestLoginHandlerMapsInvalidCredentialsToUnauthorized(t *testing.T) {
	handler := newLoginHandler(user.User{
		ID:           42,
		PasswordHash: "stored-hash",
		IsActive:     true,
		HasStaffRole: true,
	}, port.ErrPasswordMismatch, "", nil)

	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(
		`{"email":"admin@example.test","password":"incorrect-password"}`,
	))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusUnauthorized)
	}
	var body struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Message != command.ErrInvalidCredentials.Error() {
		t.Fatalf("message = %q, want stable authentication error code", body.Message)
	}
}

func TestLoginHandlerRejectsMalformedJSON(t *testing.T) {
	handler := newLoginHandler(user.User{}, nil, "", nil)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString("{"))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}

func TestLoginHandlerDoesNotExposeInternalErrors(t *testing.T) {
	handler := newLoginHandler(user.User{
		ID:           42,
		PasswordHash: "stored-hash",
		IsActive:     true,
		HasStaffRole: true,
	}, nil, "", errors.New("token signing failure"))
	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(
		`{"email":"admin@example.test","password":"submitted-password"}`,
	))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusInternalServerError)
	}
	var body struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Message != internalErrorCode {
		t.Fatalf("message = %q, want generic internal error code", body.Message)
	}
}

func newLoginHandler(account user.User, verifierErr error, token string, tokenIssuerErr error) http.Handler {
	serviceContext := &svc.ServiceContext{
		LoginAdmin: command.NewLogin(
			loginUserRepositoryStub{account: account},
			loginPasswordVerifierStub{err: verifierErr},
			loginTokenIssuerStub{token: token, err: tokenIssuerErr},
		),
	}
	return LoginHandler(serviceContext)
}

var _ user.Repository = loginUserRepositoryStub{}
var _ port.PasswordVerifier = loginPasswordVerifierStub{}
var _ port.TokenIssuer = loginTokenIssuerStub{}
