package handler

import (
	"errors"
	"net/http"

	"webtranshome/internal/logic/identityaccess"
	"webtranshome/internal/middleware"
	"webtranshome/internal/modules/identityaccess/application/command"
	"webtranshome/internal/modules/identityaccess/application/query"
	"webtranshome/internal/modules/identityaccess/domain/service"
	"webtranshome/internal/shared/response"
	"webtranshome/internal/svc"
	"webtranshome/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

const (
	accessControlBadRequestCode      = 40010
	accessControlUnauthorizedCode    = 40110
	accessControlNotFoundCode        = 40410
	accessControlConflictCode        = 40910
	accessControlInternalErrorCode   = 50000
	accessControlInvalidJSONMessage  = "HTTP_INVALID_JSON"
	accessControlUnauthorizedMessage = "AUTHENTICATION_REQUIRED"
	accessControlInternalMessage     = "INTERNAL_ERROR"
)

func AccessControlCatalogHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logic := identityaccess.NewAccessControlLogic(r.Context(), svcCtx)
		result, err := logic.Catalog()
		if err != nil {
			writeAccessControlFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

func CreateRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateRoleRequest
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			writeAccessControlError(w, r, http.StatusBadRequest, accessControlBadRequestCode, accessControlInvalidJSONMessage)
			return
		}

		logic := identityaccess.NewAccessControlLogic(r.Context(), svcCtx)
		result, err := logic.CreateRole(&req)
		if err != nil {
			writeAccessControlFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

func SearchAccessControlUsersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchAccessControlUsersRequest
		if err := httpx.ParseForm(r, &req); err != nil {
			writeAccessControlError(w, r, http.StatusBadRequest, accessControlBadRequestCode, query.ErrInvalidAccessControlQuery.Error())
			return
		}

		logic := identityaccess.NewAccessControlLogic(r.Context(), svcCtx)
		result, err := logic.FindUserByEmail(req.Email)
		if err != nil {
			writeAccessControlFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

func UpdateRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var path struct {
			RoleID int64 `path:"roleId"`
		}
		if err := httpx.ParsePath(r, &path); err != nil {
			writeAccessControlError(w, r, http.StatusBadRequest, accessControlBadRequestCode, command.ErrInvalidRole.Error())
			return
		}

		var req types.UpdateRoleRequest
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			writeAccessControlError(w, r, http.StatusBadRequest, accessControlBadRequestCode, accessControlInvalidJSONMessage)
			return
		}

		logic := identityaccess.NewAccessControlLogic(r.Context(), svcCtx)
		if err := logic.UpdateRole(&req, path.RoleID); err != nil {
			writeAccessControlFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, types.AccessControlMutationResponse{Success: true})
	}
}

func CurrentUserAccessHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.UserIDFromContext(r.Context())
		if !ok {
			writeAccessControlError(w, r, http.StatusUnauthorized, accessControlUnauthorizedCode, accessControlUnauthorizedMessage)
			return
		}

		logic := identityaccess.NewAccessControlLogic(r.Context(), svcCtx)
		result, err := logic.UserAccess(userID)
		if err != nil {
			writeAccessControlFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

func UserRolesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var path struct {
			UserID int64 `path:"userId"`
		}
		if err := httpx.ParsePath(r, &path); err != nil {
			writeAccessControlError(w, r, http.StatusBadRequest, accessControlBadRequestCode, query.ErrInvalidAccessControlQuery.Error())
			return
		}

		logic := identityaccess.NewAccessControlLogic(r.Context(), svcCtx)
		result, err := logic.UserRoleIDs(path.UserID)
		if err != nil {
			writeAccessControlFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

func ReplaceUserRolesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var path struct {
			UserID int64 `path:"userId"`
		}
		if err := httpx.ParsePath(r, &path); err != nil {
			writeAccessControlError(w, r, http.StatusBadRequest, accessControlBadRequestCode, query.ErrInvalidAccessControlQuery.Error())
			return
		}

		var req types.ReplaceUserRolesRequest
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			writeAccessControlError(w, r, http.StatusBadRequest, accessControlBadRequestCode, accessControlInvalidJSONMessage)
			return
		}

		logic := identityaccess.NewAccessControlLogic(r.Context(), svcCtx)
		if err := logic.ReplaceUserRoles(&req, path.UserID); err != nil {
			writeAccessControlFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, types.AccessControlMutationResponse{Success: true})
	}
}

func RolePermissionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var path struct {
			RoleID int64 `path:"roleId"`
		}
		if err := httpx.ParsePath(r, &path); err != nil {
			writeAccessControlError(w, r, http.StatusBadRequest, accessControlBadRequestCode, query.ErrInvalidAccessControlQuery.Error())
			return
		}

		logic := identityaccess.NewAccessControlLogic(r.Context(), svcCtx)
		result, err := logic.RolePermissionIDs(path.RoleID)
		if err != nil {
			writeAccessControlFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

func ReplaceRolePermissionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var path struct {
			RoleID int64 `path:"roleId"`
		}
		if err := httpx.ParsePath(r, &path); err != nil {
			writeAccessControlError(w, r, http.StatusBadRequest, accessControlBadRequestCode, query.ErrInvalidAccessControlQuery.Error())
			return
		}

		var req types.ReplaceRolePermissionsRequest
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			writeAccessControlError(w, r, http.StatusBadRequest, accessControlBadRequestCode, accessControlInvalidJSONMessage)
			return
		}

		logic := identityaccess.NewAccessControlLogic(r.Context(), svcCtx)
		if err := logic.ReplaceRolePermissions(&req, path.RoleID); err != nil {
			writeAccessControlFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, types.AccessControlMutationResponse{Success: true})
	}
}

func writeAccessControlFailure(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, query.ErrInvalidAccessControlQuery),
		errors.Is(err, command.ErrInvalidAccessAssignment),
		errors.Is(err, command.ErrInvalidRole):
		writeAccessControlError(w, r, http.StatusBadRequest, accessControlBadRequestCode, err.Error())
	case errors.Is(err, service.ErrUserNotFound),
		errors.Is(err, service.ErrRoleNotFound),
		errors.Is(err, service.ErrPermissionNotFound):
		writeAccessControlError(w, r, http.StatusNotFound, accessControlNotFoundCode, err.Error())
	case errors.Is(err, service.ErrProtectedRole),
		errors.Is(err, service.ErrRoleCodeConflict):
		writeAccessControlError(w, r, http.StatusConflict, accessControlConflictCode, err.Error())
	default:
		logx.WithContext(r.Context()).Error(err)
		writeAccessControlError(w, r, http.StatusInternalServerError, accessControlInternalErrorCode, accessControlInternalMessage)
	}
}

func writeAccessControlError(w http.ResponseWriter, r *http.Request, status, code int, message string) {
	httpx.WriteJsonCtx(r.Context(), w, status, response.Body{
		Code:    code,
		Message: message,
	})
}
