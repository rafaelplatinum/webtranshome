package command

import (
	"context"

	"webtranshome/internal/modules/catalog/application/dto"
	"webtranshome/internal/modules/catalog/domain/room"
)

type CreateRoom struct {
	writer room.Writer
}

type CreateRoomInput struct {
	Name       string
	Slug       string
	ImageCover string
	SortOrder  int64
}

func NewCreateRoom(writer room.Writer) *CreateRoom {
	return &CreateRoom{writer: writer}
}

func (c *CreateRoom) Execute(ctx context.Context, input CreateRoomInput) (dto.Room, error) {
	item, err := room.New(input.Name, input.Slug, input.ImageCover, input.SortOrder)
	if err != nil {
		return dto.Room{}, err
	}
	created, err := c.writer.CreateRoom(ctx, item)
	if err != nil {
		return dto.Room{}, err
	}
	return roomDTO(created), nil
}
