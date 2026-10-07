package query

import (
	"context"
	"errors"

	"webtranshome/internal/modules/catalog/application/dto"
	"webtranshome/internal/modules/catalog/domain/product"
)

var ErrInvalidProductID = errors.New("CATALOG_INVALID_PRODUCT_ID")

type GetProduct struct {
	reader product.Reader
}

func NewGetProduct(reader product.Reader) *GetProduct {
	return &GetProduct{reader: reader}
}

func (q *GetProduct) Execute(ctx context.Context, productID int64) (dto.Product, error) {
	if productID <= 0 {
		return dto.Product{}, ErrInvalidProductID
	}
	item, err := q.reader.FindProduct(ctx, productID)
	if err != nil {
		return dto.Product{}, err
	}
	return toProductDTO(item), nil
}
