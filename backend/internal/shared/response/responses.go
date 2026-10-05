package response

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// Body adalah standar format amplop untuk semua API Transhome
type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success membungkus data balikan dengan format 200 OK standar
func Success(w http.ResponseWriter, data interface{}) {
	httpx.OkJson(w, Body{
		Code:    20000,
		Message: "Success",
		Data:    data,
	})
}

// ErrorHandler memetakan error domain internal menjadi HTTP error code dan amplop yang rapi
func ErrorHandler(err error) (int, interface{}) {
	// Default error jika tidak dikenali
	statusCode := http.StatusInternalServerError
	resp := Body{
		Code:    50000,
		Message: "Terjadi kesalahan pada server",
	}

	// TODO: Nanti kita tambahkan switch-case di sini untuk memetakan error spesifik
	// Contoh: ErrDuplicateReferenceNo -> HTTP 409, Code 40920
	// Contoh: ErrInsufficientPoints -> HTTP 422, Code 42220

	if err != nil {
		resp.Message = err.Error()
	}

	return statusCode, resp
}
