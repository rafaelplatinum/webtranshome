package postgres

import (
	"errors"
	"testing"

	"webtranshome/internal/modules/catalog/domain/room"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func TestMapRoomError(t *testing.T) {
	databaseError := errors.New("database unavailable")
	tests := []struct {
		name string
		err  error
		want error
	}{
		{name: "maps missing row", err: sqlx.ErrNotFound, want: room.ErrRoomNotFound},
		{name: "preserves infrastructure error", err: databaseError, want: databaseError},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := mapRoomError(test.err); !errors.Is(got, test.want) {
				t.Fatalf("mapRoomError() = %v, want %v", got, test.want)
			}
		})
	}
}
