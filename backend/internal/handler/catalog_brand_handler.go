package handler

import (
	"errors"
	"net/http"

	"webtranshome/internal/logic/catalog"
	"webtranshome/internal/modules/catalog/application/command"
	"webtranshome/internal/modules/catalog/application/query"
	"webtranshome/internal/modules/catalog/domain/brand"
	"webtranshome/internal/shared/response"
	"webtranshome/internal/svc"
	"webtranshome/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

const (
	catalogBrandBadRequestCode = 40021
	catalogBrandNotFoundCode   = 40421
	catalogBrandConflictCode   = 40921
	catalogBrandInternalCode   = 50000
	catalogBrandInvalidJSON    = "HTTP_INVALID_JSON"
	catalogBrandInternalError  = "INTERNAL_ERROR"
)

func CatalogBrandsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := catalog.NewBrandLogic(r.Context(), svcCtx).List()
		if err != nil {
			writeCatalogBrandFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

func CatalogBrandHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var path struct {
			BrandID int64 `path:"brandId"`
		}
		if err := httpx.ParsePath(r, &path); err != nil {
			writeCatalogBrandError(w, r, http.StatusBadRequest, catalogBrandBadRequestCode, query.ErrInvalidBrandID.Error())
			return
		}
		result, err := catalog.NewBrandLogic(r.Context(), svcCtx).Get(path.BrandID)
		if err != nil {
			writeCatalogBrandFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

func CreateCatalogBrandHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateBrandRequest
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			writeCatalogBrandError(w, r, http.StatusBadRequest, catalogBrandBadRequestCode, catalogBrandInvalidJSON)
			return
		}
		result, err := catalog.NewBrandLogic(r.Context(), svcCtx).Create(&req)
		if err != nil {
			writeCatalogBrandFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

func UpdateCatalogBrandHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var path struct {
			BrandID int64 `path:"brandId"`
		}
		if err := httpx.ParsePath(r, &path); err != nil {
			writeCatalogBrandError(w, r, http.StatusBadRequest, catalogBrandBadRequestCode, query.ErrInvalidBrandID.Error())
			return
		}
		var req types.UpdateBrandRequest
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			writeCatalogBrandError(w, r, http.StatusBadRequest, catalogBrandBadRequestCode, catalogBrandInvalidJSON)
			return
		}
		result, err := catalog.NewBrandLogic(r.Context(), svcCtx).Update(&req, path.BrandID)
		if err != nil {
			writeCatalogBrandFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

func writeCatalogBrandFailure(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, brand.ErrInvalidBrand), errors.Is(err, query.ErrInvalidBrandID),
		errors.Is(err, command.ErrInvalidBrandID):
		writeCatalogBrandError(w, r, http.StatusBadRequest, catalogBrandBadRequestCode, err.Error())
	case errors.Is(err, brand.ErrBrandNotFound):
		writeCatalogBrandError(w, r, http.StatusNotFound, catalogBrandNotFoundCode, err.Error())
	case errors.Is(err, brand.ErrBrandSlugConflict):
		writeCatalogBrandError(w, r, http.StatusConflict, catalogBrandConflictCode, err.Error())
	default:
		logx.WithContext(r.Context()).Error(err)
		writeCatalogBrandError(w, r, http.StatusInternalServerError, catalogBrandInternalCode, catalogBrandInternalError)
	}
}

func writeCatalogBrandError(w http.ResponseWriter, r *http.Request, status, code int, message string) {
	httpx.WriteJsonCtx(r.Context(), w, status, response.Body{Code: code, Message: message})
}
