package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"webtranshome/internal/modules/catalog/domain/brand"
)

func TestWriteCatalogBrandFailure(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{name: "invalid brand", err: brand.ErrInvalidBrand, wantStatus: http.StatusBadRequest},
		{name: "missing brand", err: brand.ErrBrandNotFound, wantStatus: http.StatusNotFound},
		{name: "slug conflict", err: brand.ErrBrandSlugConflict, wantStatus: http.StatusConflict},
		{name: "unexpected failure", err: errors.New("database unavailable"), wantStatus: http.StatusInternalServerError},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/admin/catalog/brands", nil)
			recorder := httptest.NewRecorder()
			writeCatalogBrandFailure(recorder, request, test.err)
			if recorder.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", recorder.Code, test.wantStatus)
			}
		})
	}
}
