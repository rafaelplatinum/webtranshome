package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"webtranshome/internal/modules/catalog/domain/room"
)

func TestWriteCatalogRoomFailure(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{name: "invalid room", err: room.ErrInvalidRoom, wantStatus: http.StatusBadRequest},
		{name: "missing room", err: room.ErrRoomNotFound, wantStatus: http.StatusNotFound},
		{name: "slug conflict", err: room.ErrRoomSlugConflict, wantStatus: http.StatusConflict},
		{name: "unexpected failure", err: errors.New("database unavailable"), wantStatus: http.StatusInternalServerError},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/admin/catalog/rooms", nil)
			recorder := httptest.NewRecorder()
			writeCatalogRoomFailure(recorder, request, test.err)
			if recorder.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", recorder.Code, test.wantStatus)
			}
		})
	}
}
