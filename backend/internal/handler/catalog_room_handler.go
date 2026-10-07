package handler

import (
	"errors"
	"net/http"

	"webtranshome/internal/logic/catalog"
	"webtranshome/internal/modules/catalog/application/command"
	"webtranshome/internal/modules/catalog/application/query"
	"webtranshome/internal/modules/catalog/domain/room"
	"webtranshome/internal/shared/response"
	"webtranshome/internal/svc"
	"webtranshome/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

const (
	catalogRoomBadRequestCode = 40022
	catalogRoomNotFoundCode   = 40422
	catalogRoomConflictCode   = 40922
	catalogRoomInternalCode   = 50000
	catalogRoomInvalidJSON    = "HTTP_INVALID_JSON"
	catalogRoomInternalError  = "INTERNAL_ERROR"
)

func CatalogRoomsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := catalog.NewRoomLogic(r.Context(), svcCtx).List()
		if err != nil {
			writeCatalogRoomFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

func CatalogRoomHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var path struct {
			RoomID int64 `path:"roomId"`
		}
		if err := httpx.ParsePath(r, &path); err != nil {
			writeCatalogRoomError(w, r, http.StatusBadRequest, catalogRoomBadRequestCode, query.ErrInvalidRoomID.Error())
			return
		}
		result, err := catalog.NewRoomLogic(r.Context(), svcCtx).Get(path.RoomID)
		if err != nil {
			writeCatalogRoomFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

func CreateCatalogRoomHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateRoomRequest
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			writeCatalogRoomError(w, r, http.StatusBadRequest, catalogRoomBadRequestCode, catalogRoomInvalidJSON)
			return
		}
		result, err := catalog.NewRoomLogic(r.Context(), svcCtx).Create(&req)
		if err != nil {
			writeCatalogRoomFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

func UpdateCatalogRoomHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var path struct {
			RoomID int64 `path:"roomId"`
		}
		if err := httpx.ParsePath(r, &path); err != nil {
			writeCatalogRoomError(w, r, http.StatusBadRequest, catalogRoomBadRequestCode, query.ErrInvalidRoomID.Error())
			return
		}
		var req types.UpdateRoomRequest
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			writeCatalogRoomError(w, r, http.StatusBadRequest, catalogRoomBadRequestCode, catalogRoomInvalidJSON)
			return
		}
		result, err := catalog.NewRoomLogic(r.Context(), svcCtx).Update(&req, path.RoomID)
		if err != nil {
			writeCatalogRoomFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

func writeCatalogRoomFailure(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, room.ErrInvalidRoom), errors.Is(err, query.ErrInvalidRoomID),
		errors.Is(err, command.ErrInvalidRoomID):
		writeCatalogRoomError(w, r, http.StatusBadRequest, catalogRoomBadRequestCode, err.Error())
	case errors.Is(err, room.ErrRoomNotFound):
		writeCatalogRoomError(w, r, http.StatusNotFound, catalogRoomNotFoundCode, err.Error())
	case errors.Is(err, room.ErrRoomSlugConflict):
		writeCatalogRoomError(w, r, http.StatusConflict, catalogRoomConflictCode, err.Error())
	default:
		logx.WithContext(r.Context()).Error(err)
		writeCatalogRoomError(w, r, http.StatusInternalServerError, catalogRoomInternalCode, catalogRoomInternalError)
	}
}

func writeCatalogRoomError(w http.ResponseWriter, r *http.Request, status, code int, message string) {
	httpx.WriteJsonCtx(r.Context(), w, status, response.Body{Code: code, Message: message})
}
