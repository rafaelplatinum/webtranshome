package command

import (
	"context"

	"webtranshome/internal/modules/catalog/domain/product"
	"webtranshome/internal/modules/catalog/domain/room"
)

type ReplaceProductRoomsInput struct {
	ProductID int64
	RoomIDs   []int64
}

type ReplaceProductRooms struct {
	productReader product.Reader
	roomReader    room.Reader
	writer        product.RoomAssignmentWriter
}

func NewReplaceProductRooms(productReader product.Reader, roomReader room.Reader, writer product.RoomAssignmentWriter) *ReplaceProductRooms {
	return &ReplaceProductRooms{productReader: productReader, roomReader: roomReader, writer: writer}
}

func (c *ReplaceProductRooms) Execute(ctx context.Context, input ReplaceProductRoomsInput) error {
	assignment, err := product.NewRoomAssignment(input.ProductID, input.RoomIDs)
	if err != nil {
		return err
	}
	if _, err := c.productReader.FindProduct(ctx, assignment.ProductID); err != nil {
		return err
	}
	for _, roomID := range assignment.RoomIDs {
		if _, err := c.roomReader.FindRoom(ctx, roomID); err != nil {
			return err
		}
	}
	return c.writer.ReplaceProductRooms(ctx, assignment.ProductID, assignment.RoomIDs)
}
