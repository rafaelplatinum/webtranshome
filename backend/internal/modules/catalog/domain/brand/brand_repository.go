package brand

import "context"

type Reader interface {
	ListBrands(ctx context.Context) ([]Brand, error)
	FindBrand(ctx context.Context, brandID int64) (Brand, error)
}

type Writer interface {
	CreateBrand(ctx context.Context, item Brand) (Brand, error)
	UpdateBrand(ctx context.Context, item Brand) error
}
