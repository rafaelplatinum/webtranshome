package catalog

import (
	"context"

	"webtranshome/internal/modules/catalog/application/command"
	"webtranshome/internal/modules/catalog/application/dto"
	"webtranshome/internal/modules/catalog/domain/room"
	"webtranshome/internal/svc"
	"webtranshome/internal/types"
)

type RoomLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoomLogic {
	return &RoomLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *RoomLogic) List() (*types.CatalogRoomListResponse, error) {
	items, err := l.svcCtx.ListRooms.Execute(l.ctx)
	if err != nil {
		return nil, err
	}
	result := &types.CatalogRoomListResponse{Rooms: make([]types.CatalogRoom, 0, len(items))}
	for _, item := range items {
		result.Rooms = append(result.Rooms, catalogRoom(item))
	}
	return result, nil
}

func (l *RoomLogic) Get(roomID int64) (*types.CatalogRoom, error) {
	item, err := l.svcCtx.GetRoom.Execute(l.ctx, roomID)
	if err != nil {
		return nil, err
	}
	result := catalogRoom(item)
	return &result, nil
}

func (l *RoomLogic) Create(req *types.CreateRoomRequest) (*types.CatalogRoom, error) {
	if req == nil {
		return nil, room.ErrInvalidRoom
	}
	item, err := l.svcCtx.CreateRoom.Execute(l.ctx, command.CreateRoomInput{
		Name: req.Name, Slug: req.Slug, ImageCover: req.ImageCover, SortOrder: req.SortOrder,
	})
	if err != nil {
		return nil, err
	}
	result := catalogRoom(item)
	return &result, nil
}

func (l *RoomLogic) Update(req *types.UpdateRoomRequest, roomID int64) (*types.CatalogRoom, error) {
	if req == nil {
		return nil, room.ErrInvalidRoom
	}
	item, err := l.svcCtx.UpdateRoom.Execute(l.ctx, command.UpdateRoomInput{
		ID: roomID, Name: req.Name, Slug: req.Slug, ImageCover: req.ImageCover,
		SortOrder: req.SortOrder, IsActive: req.IsActive,
	})
	if err != nil {
		return nil, err
	}
	result := catalogRoom(item)
	return &result, nil
}

func catalogRoom(item dto.Room) types.CatalogRoom {
	return types.CatalogRoom{
		ID: item.ID, Name: item.Name, Slug: item.Slug,
		ImageCover: item.ImageCover, SortOrder: item.SortOrder, IsActive: item.IsActive,
	}
}
