package command

import (
	"context"
	"errors"

	"webtranshome/internal/modules/catalog/application/dto"
	"webtranshome/internal/modules/catalog/domain/brand"
)

var ErrInvalidBrandID = errors.New("CATALOG_INVALID_BRAND_ID")

type UpdateBrand struct {
	reader brand.Reader
	writer brand.Writer
}

type UpdateBrandInput struct {
	ID       int64
	Name     *string
	Slug     *string
	LogoURL  *string
	IsActive *bool
}

func NewUpdateBrand(reader brand.Reader, writer brand.Writer) *UpdateBrand {
	return &UpdateBrand{reader: reader, writer: writer}
}

func (c *UpdateBrand) Execute(ctx context.Context, input UpdateBrandInput) (dto.Brand, error) {
	if input.ID <= 0 {
		return dto.Brand{}, ErrInvalidBrandID
	}
	item, err := c.reader.FindBrand(ctx, input.ID)
	if err != nil {
		return dto.Brand{}, err
	}
	if err := item.Apply(brand.Changes{
		Name: input.Name, Slug: input.Slug, LogoURL: input.LogoURL, IsActive: input.IsActive,
	}); err != nil {
		return dto.Brand{}, err
	}
	if err := c.writer.UpdateBrand(ctx, item); err != nil {
		return dto.Brand{}, err
	}
	return brandDTO(item), nil
}
