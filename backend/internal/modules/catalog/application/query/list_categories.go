package query

import (
	"context"

	"webtranshome/internal/modules/catalog/application/dto"
	"webtranshome/internal/modules/catalog/domain/category"
)

type ListCategories struct {
	reader category.Reader
}

func NewListCategories(reader category.Reader) *ListCategories {
	return &ListCategories{reader: reader}
}

func (q *ListCategories) Execute(ctx context.Context) ([]dto.Category, error) {
	items, err := q.reader.ListCategories(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]dto.Category, 0, len(items))
	for _, item := range items {
		result = append(result, toCategoryDTO(item))
	}
	return result, nil
}

func toCategoryDTO(item category.Category) dto.Category {
	imageURL := ""
	if item.ImageURL != nil {
		imageURL = *item.ImageURL
	}
	return dto.Category{
		ID:        item.ID,
		Name:      item.Name,
		Slug:      item.Slug,
		ParentID:  item.ParentID,
		ImageURL:  imageURL,
		SortOrder: item.SortOrder,
		IsActive:  item.IsActive,
	}
}
