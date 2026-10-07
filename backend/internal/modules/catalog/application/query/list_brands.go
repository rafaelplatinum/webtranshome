package query

import (
	"context"

	"webtranshome/internal/modules/catalog/application/dto"
	"webtranshome/internal/modules/catalog/domain/brand"
)

type ListBrands struct {
	reader brand.Reader
}

func NewListBrands(reader brand.Reader) *ListBrands {
	return &ListBrands{reader: reader}
}

func (q *ListBrands) Execute(ctx context.Context) ([]dto.Brand, error) {
	items, err := q.reader.ListBrands(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]dto.Brand, 0, len(items))
	for _, item := range items {
		result = append(result, brandDTO(item))
	}
	return result, nil
}

func brandDTO(item brand.Brand) dto.Brand {
	logoURL := ""
	if item.LogoURL != nil {
		logoURL = *item.LogoURL
	}
	return dto.Brand{
		ID: item.ID, Name: item.Name, Slug: item.Slug, LogoURL: logoURL, IsActive: item.IsActive,
	}
}
