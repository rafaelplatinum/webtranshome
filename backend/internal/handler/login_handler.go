package handler

import (
	"errors"
	"net/http"

	loginlogic "webtranshome/internal/logic/login"
	"webtranshome/internal/modules/identityaccess/application/command"
	"webtranshome/internal/shared/response"
	"webtranshome/internal/svc"
	"webtranshome/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

const (
	invalidJSONErrorCode       = "HTTP_INVALID_JSON"
	internalErrorCode          = "INTERNAL_ERROR"
	badRequestResponseCode     = 40000
	unauthorizedResponseCode   = 40100
	internalServerResponseCode = 50000
)

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginRequest
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusBadRequest, response.Body{
				Code:    badRequestResponseCode,
				Message: invalidJSONErrorCode,
			})
			return
		}

		logic := loginlogic.NewLoginLogic(r.Context(), svcCtx)
		result, err := logic.Login(&req)
		if err != nil {
			switch {
			case errors.Is(err, command.ErrInvalidLoginRequest):
				httpx.WriteJsonCtx(r.Context(), w, http.StatusBadRequest, response.Body{
					Code:    badRequestResponseCode,
					Message: err.Error(),
				})
			case errors.Is(err, command.ErrInvalidCredentials):
				httpx.WriteJsonCtx(r.Context(), w, http.StatusUnauthorized, response.Body{
					Code:    unauthorizedResponseCode,
					Message: err.Error(),
				})
			default:
				logx.WithContext(r.Context()).Error(err)
				httpx.WriteJsonCtx(r.Context(), w, http.StatusInternalServerError, response.Body{
					Code:    internalServerResponseCode,
					Message: internalErrorCode,
				})
			}
			return
		}

		httpx.OkJsonCtx(r.Context(), w, result)
	}
}
