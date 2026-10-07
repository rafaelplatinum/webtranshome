package query

import (
	"context"
	"errors"

	"webtranshome/internal/modules/catalog/application/dto"
	"webtranshome/internal/modules/catalog/domain/room"
)

var ErrInvalidRoomID = errors.New("CATALOG_INVALID_ROOM_ID")

type GetRoom struct {
	reader room.Reader
}

func NewGetRoom(reader room.Reader) *GetRoom {
	return &GetRoom{reader: reader}
}

func (q *GetRoom) Execute(ctx context.Context, roomID int64) (dto.Room, error) {
	if roomID <= 0 {
		return dto.Room{}, ErrInvalidRoomID
	}
	item, err := q.reader.FindRoom(ctx, roomID)
	if err != nil {
		return dto.Room{}, err
	}
	return roomDTO(item), nil
}
