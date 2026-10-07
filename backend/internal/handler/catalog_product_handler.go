package handler

import (
	"errors"
	"net/http"

	"webtranshome/internal/logic/catalog"
	"webtranshome/internal/modules/catalog/application/command"
	"webtranshome/internal/modules/catalog/application/query"
	"webtranshome/internal/modules/catalog/domain/product"
	"webtranshome/internal/shared/response"
	"webtranshome/internal/svc"
	"webtranshome/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

const (
	catalogProductBadRequestCode = 40025
	catalogProductNotFoundCode   = 40425
	catalogProductConflictCode   = 40925
	catalogProductInternalCode   = 50000
	catalogProductInvalidJSON    = "HTTP_INVALID_JSON"
	catalogProductInternalError  = "INTERNAL_ERROR"
)

func CatalogProductsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := catalog.NewProductLogic(r.Context(), svcCtx).List()
		if err != nil {
			writeCatalogProductFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

func CatalogProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var path struct {
			ProductID int64 `path:"productId"`
		}
		if err := httpx.ParsePath(r, &path); err != nil {
			writeCatalogProductError(w, r, http.StatusBadRequest, catalogProductBadRequestCode, query.ErrInvalidProductID.Error())
			return
		}
		result, err := catalog.NewProductLogic(r.Context(), svcCtx).Get(path.ProductID)
		if err != nil {
			writeCatalogProductFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

func CreateCatalogProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateProductRequest
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			writeCatalogProductError(w, r, http.StatusBadRequest, catalogProductBadRequestCode, catalogProductInvalidJSON)
			return
		}
		result, err := catalog.NewProductLogic(r.Context(), svcCtx).Create(&req)
		if err != nil {
			writeCatalogProductFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

func UpdateCatalogProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var path struct {
			ProductID int64 `path:"productId"`
		}
		if err := httpx.ParsePath(r, &path); err != nil {
			writeCatalogProductError(w, r, http.StatusBadRequest, catalogProductBadRequestCode, query.ErrInvalidProductID.Error())
			return
		}
		var req types.UpdateProductRequest
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			writeCatalogProductError(w, r, http.StatusBadRequest, catalogProductBadRequestCode, catalogProductInvalidJSON)
			return
		}
		result, err := catalog.NewProductLogic(r.Context(), svcCtx).Update(&req, path.ProductID)
		if err != nil {
			writeCatalogProductFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

func ReplaceCatalogProductRoomsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var path struct {
			ProductID int64 `path:"productId"`
		}
		if err := httpx.ParsePath(r, &path); err != nil {
			writeCatalogProductError(w, r, http.StatusBadRequest, catalogProductBadRequestCode, query.ErrInvalidProductID.Error())
			return
		}
		var req types.ReplaceProductRoomsRequest
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			writeCatalogProductError(w, r, http.StatusBadRequest, catalogProductBadRequestCode, catalogProductInvalidJSON)
			return
		}
		if err := catalog.NewProductLogic(r.Context(), svcCtx).ReplaceRooms(&req, path.ProductID); err != nil {
			writeCatalogProductFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, types.CatalogMutationResponse{Success: true})
	}
}

func writeCatalogProductFailure(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, product.ErrInvalidProduct), errors.Is(err, product.ErrInvalidCategoryID),
		errors.Is(err, product.ErrInvalidBrandID), errors.Is(err, product.ErrInvalidRoomID),
		errors.Is(err, product.ErrInvalidStockStatus), errors.Is(err, query.ErrInvalidProductID),
		errors.Is(err, command.ErrInvalidProductID):
		writeCatalogProductError(w, r, http.StatusBadRequest, catalogProductBadRequestCode, err.Error())
	case errors.Is(err, product.ErrProductNotFound):
		writeCatalogProductError(w, r, http.StatusNotFound, catalogProductNotFoundCode, err.Error())
	case errors.Is(err, product.ErrProductSKUConflict), errors.Is(err, product.ErrProductSlugConflict):
		writeCatalogProductError(w, r, http.StatusConflict, catalogProductConflictCode, err.Error())
	default:
		logx.WithContext(r.Context()).Error(err)
		writeCatalogProductError(w, r, http.StatusInternalServerError, catalogProductInternalCode, catalogProductInternalError)
	}
}

func writeCatalogProductError(w http.ResponseWriter, r *http.Request, status, code int, message string) {
	httpx.WriteJsonCtx(r.Context(), w, status, response.Body{Code: code, Message: message})
}
