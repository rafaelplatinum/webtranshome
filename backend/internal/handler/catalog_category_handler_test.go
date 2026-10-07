package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"webtranshome/internal/modules/catalog/domain/category"
)

func TestWriteCatalogCategoryFailure(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{name: "invalid category", err: category.ErrInvalidCategory, wantStatus: http.StatusBadRequest},
		{name: "missing category", err: category.ErrCategoryNotFound, wantStatus: http.StatusNotFound},
		{name: "slug conflict", err: category.ErrCategorySlugConflict, wantStatus: http.StatusConflict},
		{name: "unexpected failure", err: errors.New("database unavailable"), wantStatus: http.StatusInternalServerError},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertCategoryFailureStatus(t, test.err, test.wantStatus)
		})
	}
}

func assertCategoryFailureStatus(t *testing.T, err error, wantStatus int) {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, "/admin/catalog/categories", nil)
	recorder := httptest.NewRecorder()
	writeCatalogCategoryFailure(recorder, request, err)
	if recorder.Code != wantStatus {
		t.Fatalf("status = %d, want %d", recorder.Code, wantStatus)
	}
}
