package command

import (
	"context"

	"webtranshome/internal/modules/catalog/application/dto"
	"webtranshome/internal/modules/catalog/domain/brand"
)

type CreateBrand struct {
	writer brand.Writer
}

type CreateBrandInput struct {
	Name    string
	Slug    string
	LogoURL string
}

func NewCreateBrand(writer brand.Writer) *CreateBrand {
	return &CreateBrand{writer: writer}
}

func (c *CreateBrand) Execute(ctx context.Context, input CreateBrandInput) (dto.Brand, error) {
	item, err := brand.New(input.Name, input.Slug, input.LogoURL)
	if err != nil {
		return dto.Brand{}, err
	}
	created, err := c.writer.CreateBrand(ctx, item)
	if err != nil {
		return dto.Brand{}, err
	}
	return brandDTO(created), nil
}
