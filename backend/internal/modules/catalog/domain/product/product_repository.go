package product

import "context"

type Reader interface {
	ListProducts(ctx context.Context) ([]Product, error)
	FindProduct(ctx context.Context, productID int64) (Product, error)
}

type Writer interface {
	CreateProduct(ctx context.Context, item Product) (Product, error)
	UpdateProduct(ctx context.Context, item Product) error
}

type RoomAssignmentWriter interface {
	ReplaceProductRooms(ctx context.Context, productID int64, roomIDs []int64) error
}
