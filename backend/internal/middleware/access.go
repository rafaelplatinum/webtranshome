package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"webtranshome/internal/modules/identityaccess/application/port"
	"webtranshome/internal/shared/response"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

const (
	authenticationRequiredCode = 40110
	permissionDeniedCode       = 40310
	internalErrorCode          = 50000

	authenticationRequiredMessage = "AUTHENTICATION_REQUIRED"
	permissionDeniedMessage       = "AUTHORIZATION_DENIED"
	internalErrorMessage          = "INTERNAL_ERROR"
)

type PermissionChecker interface {
	Execute(ctx context.Context, userID int64, permissionCode string) (bool, error)
}

type principalContextKey struct{}

func UserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(principalContextKey{}).(int64)
	return userID, ok && userID > 0
}

func Authenticate(verifier port.TokenVerifier) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			token, ok := bearerToken(r.Header.Get("Authorization"))
			if !ok {
				writeAccessError(r, w, http.StatusUnauthorized, authenticationRequiredCode, authenticationRequiredMessage)
				return
			}

			userID, err := verifier.Verify(r.Context(), token)
			if err != nil {
				if errors.Is(err, port.ErrInvalidAccessToken) {
					writeAccessError(r, w, http.StatusUnauthorized, authenticationRequiredCode, authenticationRequiredMessage)
					return
				}

				logx.WithContext(r.Context()).Error(err)
				writeAccessError(r, w, http.StatusInternalServerError, internalErrorCode, internalErrorMessage)
				return
			}
			if userID <= 0 {
				writeAccessError(r, w, http.StatusUnauthorized, authenticationRequiredCode, authenticationRequiredMessage)
				return
			}

			ctx := context.WithValue(r.Context(), principalContextKey{}, userID)
			next(w, r.WithContext(ctx))
		}
	}
}

func RequirePermission(checker PermissionChecker, permissionCode string) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			userID, ok := UserIDFromContext(r.Context())
			if !ok {
				writeAccessError(r, w, http.StatusUnauthorized, authenticationRequiredCode, authenticationRequiredMessage)
				return
			}

			granted, err := checker.Execute(r.Context(), userID, permissionCode)
			if err != nil {
				logx.WithContext(r.Context()).Error(err)
				writeAccessError(r, w, http.StatusInternalServerError, internalErrorCode, internalErrorMessage)
				return
			}
			if !granted {
				writeAccessError(r, w, http.StatusForbidden, permissionDeniedCode, permissionDeniedMessage)
				return
			}

			next(w, r)
		}
	}
}

func bearerToken(header string) (string, bool) {
	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", false
	}

	return parts[1], true
}

func writeAccessError(r *http.Request, w http.ResponseWriter, status, code int, message string) {
	httpx.WriteJsonCtx(r.Context(), w, status, response.Body{
		Code:    code,
		Message: message,
	})
}
