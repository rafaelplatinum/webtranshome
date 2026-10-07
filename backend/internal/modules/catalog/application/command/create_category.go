package command

import (
	"context"

	"webtranshome/internal/modules/catalog/application/dto"
	"webtranshome/internal/modules/catalog/domain/category"
)

var ErrInvalidCategory = category.ErrInvalidCategory

type CreateCategoryInput struct {
	Name      string
	Slug      string
	ParentID  *int64
	ImageURL  string
	SortOrder int64
}

type CreateCategory struct {
	reader category.Reader
	writer category.Writer
}

func NewCreateCategory(reader category.Reader, writer category.Writer) *CreateCategory {
	return &CreateCategory{reader: reader, writer: writer}
}

func (c *CreateCategory) Execute(ctx context.Context, input CreateCategoryInput) (dto.Category, error) {
	item, err := category.New(input.Name, input.Slug, input.ParentID, input.ImageURL, input.SortOrder)
	if err != nil {
		return dto.Category{}, err
	}
	if err := validateActiveParent(ctx, c.reader, item.ParentID); err != nil {
		return dto.Category{}, err
	}
	created, err := c.writer.CreateCategory(ctx, item)
	if err != nil {
		return dto.Category{}, err
	}
	return categoryDTO(created), nil
}

func validateActiveParent(ctx context.Context, reader category.Reader, parentID *int64) error {
	if parentID == nil || *parentID == 0 {
		return nil
	}
	parent, err := reader.FindCategory(ctx, *parentID)
	if err != nil {
		return err
	}
	if !parent.IsActive {
		return category.ErrInvalidParent
	}
	return nil
}

func categoryDTO(item category.Category) dto.Category {
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
