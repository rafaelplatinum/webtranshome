package command

import (
	"context"
	"errors"

	"webtranshome/internal/modules/catalog/application/dto"
	"webtranshome/internal/modules/catalog/domain/room"
)

var ErrInvalidRoomID = errors.New("CATALOG_INVALID_ROOM_ID")

type UpdateRoom struct {
	reader room.Reader
	writer room.Writer
}

type UpdateRoomInput struct {
	ID         int64
	Name       *string
	Slug       *string
	ImageCover *string
	SortOrder  *int64
	IsActive   *bool
}

func NewUpdateRoom(reader room.Reader, writer room.Writer) *UpdateRoom {
	return &UpdateRoom{reader: reader, writer: writer}
}

func (c *UpdateRoom) Execute(ctx context.Context, input UpdateRoomInput) (dto.Room, error) {
	if input.ID <= 0 {
		return dto.Room{}, ErrInvalidRoomID
	}
	item, err := c.reader.FindRoom(ctx, input.ID)
	if err != nil {
		return dto.Room{}, err
	}
	if err := item.Apply(room.Changes{
		Name: input.Name, Slug: input.Slug, ImageCover: input.ImageCover,
		SortOrder: input.SortOrder, IsActive: input.IsActive,
	}); err != nil {
		return dto.Room{}, err
	}
	if err := c.writer.UpdateRoom(ctx, item); err != nil {
		return dto.Room{}, err
	}
	return roomDTO(item), nil
}
