package category

import "context"

type Reader interface {
	ListCategories(ctx context.Context) ([]Category, error)
	FindCategory(ctx context.Context, categoryID int64) (Category, error)
	DescendantIDs(ctx context.Context, categoryID int64) ([]int64, error)
}

type Writer interface {
	CreateCategory(ctx context.Context, item Category) (Category, error)
	UpdateCategory(ctx context.Context, item Category) error
}
