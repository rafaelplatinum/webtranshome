package query

import (
	"context"
	"errors"

	"webtranshome/internal/modules/catalog/application/dto"
	"webtranshome/internal/modules/catalog/domain/brand"
)

var ErrInvalidBrandID = errors.New("CATALOG_INVALID_BRAND_ID")

type GetBrand struct {
	reader brand.Reader
}

func NewGetBrand(reader brand.Reader) *GetBrand {
	return &GetBrand{reader: reader}
}

func (q *GetBrand) Execute(ctx context.Context, brandID int64) (dto.Brand, error) {
	if brandID <= 0 {
		return dto.Brand{}, ErrInvalidBrandID
	}
	item, err := q.reader.FindBrand(ctx, brandID)
	if err != nil {
		return dto.Brand{}, err
	}
	return brandDTO(item), nil
}
