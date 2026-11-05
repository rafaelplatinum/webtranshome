package handler

import (
	"errors"
	"io"
	"mime"
	"net"
	"net/http"
	"path/filepath"
	"strings"

	"webtranshome/internal/middleware"
	"webtranshome/internal/modules/media/application/assets"
	"webtranshome/internal/modules/media/domain/asset"
	"webtranshome/internal/shared/response"
	"webtranshome/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

const maxImageUploadBytes = 2 * 1024 * 1024

func ListImageAssetsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := svcCtx.ImageAssets.List(r.Context(), strings.TrimSpace(r.URL.Query().Get("slot_key")))
		if err != nil {
			writeImageAssetFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, response.Body{Data: items})
	}
}

func CreateImageAssetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.UserIDFromContext(r.Context())
		if !ok {
			writeImageAssetError(w, r, http.StatusUnauthorized, 40110, "AUTHENTICATION_REQUIRED")
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxImageUploadBytes+64*1024)
		if err := r.ParseMultipartForm(128 * 1024); err != nil {
			writeImageAssetError(w, r, http.StatusBadRequest, 40080, "HTTP_INVALID_MULTIPART")
			return
		}
		file, header, err := r.FormFile("file")
		if err != nil {
			writeImageAssetError(w, r, http.StatusBadRequest, 40080, "IMAGE_FILE_REQUIRED")
			return
		}
		defer file.Close()
		data, err := io.ReadAll(io.LimitReader(file, maxImageUploadBytes+1))
		if err != nil || len(data) > maxImageUploadBytes {
			writeImageAssetError(w, r, http.StatusBadRequest, 40080, asset.ErrInvalidAsset.Error())
			return
		}
		item, err := svcCtx.ImageAssets.Upload(r.Context(), r.FormValue("slot_key"), r.FormValue("alt_text"), header.Filename, data, userID, clientIP(r))
		if err != nil {
			writeImageAssetFailure(w, r, err)
			return
		}
		httpx.WriteJsonCtx(r.Context(), w, http.StatusCreated, response.Body{Data: item})
	}
}

func UpdateImageAssetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.UserIDFromContext(r.Context())
		if !ok {
			writeImageAssetError(w, r, http.StatusUnauthorized, 40110, "AUTHENTICATION_REQUIRED")
			return
		}
		var path struct {
			ID int64 `path:"id"`
		}
		if err := httpx.ParsePath(r, &path); err != nil {
			writeImageAssetError(w, r, http.StatusBadRequest, 40080, asset.ErrInvalidAsset.Error())
			return
		}
		var request struct {
			AltText   *string `json:"alt_text"`
			SortOrder *int    `json:"sort_order"`
			IsActive  *bool   `json:"is_active"`
		}
		if err := httpx.ParseJsonBody(r, &request); err != nil {
			writeImageAssetError(w, r, http.StatusBadRequest, 40080, "HTTP_INVALID_JSON")
			return
		}
		item, err := svcCtx.ImageAssets.Update(r.Context(), path.ID, userID, assets.Changes{AltText: request.AltText, SortOrder: request.SortOrder, IsActive: request.IsActive}, clientIP(r))
		if err != nil {
			writeImageAssetFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, response.Body{Data: item})
	}
}

func DeleteImageAssetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.UserIDFromContext(r.Context())
		if !ok {
			writeImageAssetError(w, r, http.StatusUnauthorized, 40110, "AUTHENTICATION_REQUIRED")
			return
		}
		var path struct {
			ID int64 `path:"id"`
		}
		if err := httpx.ParsePath(r, &path); err != nil {
			writeImageAssetError(w, r, http.StatusBadRequest, 40080, asset.ErrInvalidAsset.Error())
			return
		}
		if err := svcCtx.ImageAssets.Delete(r.Context(), path.ID, userID, clientIP(r)); err != nil {
			writeImageAssetFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, response.Body{Message: "Gambar dihapus."})
	}
}

func ImageAssetFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var path struct {
			Filename string `path:"filename"`
		}
		if err := httpx.ParsePath(r, &path); err != nil {
			http.NotFound(w, r)
			return
		}
		filePath, err := svcCtx.ImageAssets.Open(path.Filename)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if contentType := mime.TypeByExtension(filepath.Ext(filePath)); contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
		w.Header().Set("X-Content-Type-Options", "nosniff")
		http.ServeFile(w, r, filePath)
	}
}

func writeImageAssetFailure(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, asset.ErrInvalidAsset), errors.Is(err, asset.ErrInvalidImage):
		writeImageAssetError(w, r, http.StatusBadRequest, 40080, err.Error())
	case errors.Is(err, asset.ErrNotFound):
		writeImageAssetError(w, r, http.StatusNotFound, 40480, err.Error())
	default:
		logx.WithContext(r.Context()).Error(err)
		writeImageAssetError(w, r, http.StatusInternalServerError, 50000, "INTERNAL_ERROR")
	}
}

func writeImageAssetError(w http.ResponseWriter, r *http.Request, status, code int, message string) {
	httpx.WriteJsonCtx(r.Context(), w, status, response.Body{Code: code, Message: message})
}

func clientIP(r *http.Request) string {
	value := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	if value != "" {
		return strings.TrimSpace(strings.Split(value, ",")[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}
