package command

import (
	"webtranshome/internal/modules/catalog/application/dto"
	"webtranshome/internal/modules/catalog/domain/room"
)

func roomDTO(item room.Room) dto.Room {
	return dto.Room{
		ID: item.ID, Name: item.Name, Slug: item.Slug,
		ImageCover: item.ImageCover, SortOrder: item.SortOrder, IsActive: item.IsActive,
	}
}
