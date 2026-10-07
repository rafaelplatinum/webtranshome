package room

import "context"

type Reader interface {
	ListRooms(ctx context.Context) ([]Room, error)
	FindRoom(ctx context.Context, roomID int64) (Room, error)
}

type Writer interface {
	CreateRoom(ctx context.Context, item Room) (Room, error)
	UpdateRoom(ctx context.Context, item Room) error
}
