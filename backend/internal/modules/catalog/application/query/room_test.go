package query

import (
	"context"
	"errors"
	"testing"

	"webtranshome/internal/modules/catalog/domain/room"
)

type roomReaderStub struct {
	items []room.Room
	item  room.Room
	err   error
}

func (r *roomReaderStub) ListRooms(context.Context) ([]room.Room, error) {
	return r.items, r.err
}

func (r *roomReaderStub) FindRoom(context.Context, int64) (room.Room, error) {
	return r.item, r.err
}

func TestListRooms(t *testing.T) {
	tests := []struct {
		name  string
		items []room.Room
		err   error
		want  int
	}{
		{name: "maps room", items: []room.Room{{ID: 5, Name: "Living Room", Slug: "living-room", ImageCover: "/cover.png", IsActive: true}}, want: 1},
		{name: "propagates repository error", err: errors.New("database unavailable")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := NewListRooms(&roomReaderStub{items: test.items, err: test.err}).Execute(context.Background())
			if !errors.Is(err, test.err) {
				t.Fatalf("Execute() error = %v, want %v", err, test.err)
			}
			if len(result) != test.want {
				t.Fatalf("Execute() rooms = %d, want %d", len(result), test.want)
			}
		})
	}
}

func TestGetRoom(t *testing.T) {
	tests := []struct {
		name      string
		id        int64
		item      room.Room
		err       error
		wantError error
	}{
		{name: "rejects invalid id", wantError: ErrInvalidRoomID},
		{name: "returns room", id: 5, item: room.Room{ID: 5, Name: "Living Room", Slug: "living-room", ImageCover: "/cover.png"}},
		{name: "propagates not found", id: 8, err: room.ErrRoomNotFound, wantError: room.ErrRoomNotFound},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := NewGetRoom(&roomReaderStub{item: test.item, err: test.err}).Execute(context.Background(), test.id)
			if !errors.Is(err, test.wantError) {
				t.Fatalf("Execute() error = %v, want %v", err, test.wantError)
			}
			if test.wantError == nil && result.ID != test.item.ID {
				t.Fatalf("Execute() room = %+v", result)
			}
		})
	}
}

var _ room.Reader = (*roomReaderStub)(nil)
