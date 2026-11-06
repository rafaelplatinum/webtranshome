package handler

import (
	"errors"
	"net"
	"net/http"
	"strings"

	"webtranshome/internal/middleware"
	"webtranshome/internal/modules/identityaccess/application/passwordreset"
	"webtranshome/internal/shared/response"
	"webtranshome/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type otpRequest struct {
	Channel     string `json:"channel"`
	Destination string `json:"destination"`
	Code        string `json:"code"`
}

type resetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type changePhoneRequest struct {
	PhoneNumber *string `json:"phone_number"`
}

func ForgotPasswordOTPHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req otpRequest
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			writeOTPError(r, w, http.StatusBadRequest, "Permintaan tidak valid.")
			return
		}
		err := svcCtx.PasswordReset.RequestForgot(r.Context(), req.Channel, req.Destination, requestIP(r))
		if errors.Is(err, passwordreset.ErrInvalidInput) {
			writeOTPError(r, w, http.StatusUnprocessableEntity, "Periksa kembali alamat email atau nomor WhatsApp.")
			return
		}
		if err != nil {
			logx.WithContext(r.Context()).Error(err)
			writeOTPError(r, w, http.StatusInternalServerError, "Terjadi kendala. Coba lagi nanti.")
			return
		}
		httpx.OkJsonCtx(r.Context(), w, map[string]any{"message": passwordreset.GenericMessage})
	}
}

func VerifyPasswordOTPHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req otpRequest
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			writeOTPError(r, w, http.StatusBadRequest, "Permintaan tidak valid.")
			return
		}
		token, err := svcCtx.PasswordReset.VerifyForgot(r.Context(), req.Channel, req.Destination, req.Code, requestIP(r))
		if errors.Is(err, passwordreset.ErrRateLimited) {
			writeOTPError(r, w, http.StatusTooManyRequests, "Terlalu banyak percobaan. Coba lagi nanti.")
			return
		}
		if err != nil {
			writeOTPError(r, w, http.StatusUnprocessableEntity, "Kode tidak valid atau sudah kedaluwarsa.")
			return
		}
		httpx.OkJsonCtx(r.Context(), w, map[string]any{"data": map[string]string{"reset_token": token}})
	}
}

func ResetPasswordOTPHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req resetPasswordRequest
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			writeOTPError(r, w, http.StatusBadRequest, "Permintaan tidak valid.")
			return
		}
		if err := svcCtx.PasswordReset.Reset(r.Context(), req.Token, req.Password, requestIP(r)); err != nil {
			if errors.Is(err, passwordreset.ErrInvalidInput) {
				writeOTPError(r, w, http.StatusUnprocessableEntity, "Password minimal 8 karakter.")
				return
			}
			if errors.Is(err, passwordreset.ErrInvalidToken) {
				writeOTPError(r, w, http.StatusUnprocessableEntity, "TOKEN_TIDAK_VALID")
				return
			}
			logx.WithContext(r.Context()).Error(err)
			writeOTPError(r, w, http.StatusInternalServerError, "Terjadi kendala. Coba lagi nanti.")
			return
		}
		httpx.OkJsonCtx(r.Context(), w, map[string]string{"message": "Password diperbarui."})
	}
}

func RequestPhoneOTPHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.UserIDFromContext(r.Context())
		if !ok {
			writeOTPError(r, w, http.StatusUnauthorized, "AUTHENTICATION_REQUIRED")
			return
		}
		if err := svcCtx.PasswordReset.RequestPhone(r.Context(), userID, requestIP(r)); err != nil {
			writePhoneError(r, w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, map[string]string{"message": "Jika nomor valid, kode verifikasi sudah dikirim."})
	}
}

func VerifyPhoneOTPHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.UserIDFromContext(r.Context())
		if !ok {
			writeOTPError(r, w, http.StatusUnauthorized, "AUTHENTICATION_REQUIRED")
			return
		}
		var req struct {
			Code string `json:"code"`
		}
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			writeOTPError(r, w, http.StatusBadRequest, "Permintaan tidak valid.")
			return
		}
		if err := svcCtx.PasswordReset.VerifyPhone(r.Context(), userID, req.Code, requestIP(r)); err != nil {
			writePhoneError(r, w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, map[string]string{"message": "Nomor WhatsApp terverifikasi."})
	}
}

func ChangePhoneHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.UserIDFromContext(r.Context())
		if !ok {
			writeOTPError(r, w, http.StatusUnauthorized, "AUTHENTICATION_REQUIRED")
			return
		}
		var req changePhoneRequest
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			writeOTPError(r, w, http.StatusBadRequest, "Permintaan tidak valid.")
			return
		}
		phone := ""
		if req.PhoneNumber != nil {
			phone = *req.PhoneNumber
		}
		if err := svcCtx.PasswordReset.ChangePhone(r.Context(), userID, phone, requestIP(r)); err != nil {
			if errors.Is(err, passwordreset.ErrInvalidInput) {
				writeOTPError(r, w, http.StatusUnprocessableEntity, "Format nomor WhatsApp tidak valid.")
				return
			}
			logx.WithContext(r.Context()).Error(err)
			writeOTPError(r, w, http.StatusInternalServerError, "Terjadi kendala. Coba lagi nanti.")
			return
		}
		httpx.OkJsonCtx(r.Context(), w, map[string]string{"message": "Nomor WhatsApp diperbarui. Verifikasi nomor untuk memulihkan password melalui WhatsApp."})
	}
}

func writePhoneError(r *http.Request, w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, passwordreset.ErrRateLimited):
		writeOTPError(r, w, http.StatusTooManyRequests, "Terlalu banyak percobaan. Coba lagi nanti.")
	case errors.Is(err, passwordreset.ErrNoPhone), errors.Is(err, passwordreset.ErrInvalidInput):
		writeOTPError(r, w, http.StatusUnprocessableEntity, "Tambahkan nomor WhatsApp valid terlebih dahulu.")
	case errors.Is(err, passwordreset.ErrInvalidCode):
		writeOTPError(r, w, http.StatusUnprocessableEntity, "Kode tidak valid atau sudah kedaluwarsa.")
	default:
		logx.WithContext(r.Context()).Error(err)
		writeOTPError(r, w, http.StatusInternalServerError, "Terjadi kendala. Coba lagi nanti.")
	}
}

func writeOTPError(r *http.Request, w http.ResponseWriter, status int, message string) {
	code := status * 100
	httpx.WriteJsonCtx(r.Context(), w, status, response.Body{Code: code, Message: message})
}

func requestIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return strings.TrimSpace(r.RemoteAddr)
}
