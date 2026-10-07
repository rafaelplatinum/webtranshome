package query

import (
	"context"

	"webtranshome/internal/modules/catalog/application/dto"
	"webtranshome/internal/modules/catalog/domain/room"
)

type ListRooms struct {
	reader room.Reader
}

func NewListRooms(reader room.Reader) *ListRooms {
	return &ListRooms{reader: reader}
}

func (q *ListRooms) Execute(ctx context.Context) ([]dto.Room, error) {
	items, err := q.reader.ListRooms(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]dto.Room, 0, len(items))
	for _, item := range items {
		result = append(result, roomDTO(item))
	}
	return result, nil
}

func roomDTO(item room.Room) dto.Room {
	return dto.Room{
		ID: item.ID, Name: item.Name, Slug: item.Slug,
		ImageCover: item.ImageCover, SortOrder: item.SortOrder, IsActive: item.IsActive,
	}
}
