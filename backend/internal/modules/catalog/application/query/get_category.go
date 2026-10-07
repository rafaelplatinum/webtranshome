package query

import (
	"context"
	"errors"

	"webtranshome/internal/modules/catalog/application/dto"
	"webtranshome/internal/modules/catalog/domain/category"
)

var ErrInvalidCategoryID = errors.New("CATALOG_INVALID_CATEGORY_ID")

type GetCategory struct {
	reader category.Reader
}

func NewGetCategory(reader category.Reader) *GetCategory {
	return &GetCategory{reader: reader}
}

func (q *GetCategory) Execute(ctx context.Context, categoryID int64) (dto.Category, error) {
	if categoryID <= 0 {
		return dto.Category{}, ErrInvalidCategoryID
	}
	item, err := q.reader.FindCategory(ctx, categoryID)
	if err != nil {
		return dto.Category{}, err
	}
	return toCategoryDTO(item), nil
}
