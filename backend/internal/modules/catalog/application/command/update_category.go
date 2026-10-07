package command

import (
	"context"
	"errors"

	"webtranshome/internal/modules/catalog/application/dto"
	"webtranshome/internal/modules/catalog/domain/category"
)

var ErrInvalidCategoryID = errors.New("CATALOG_INVALID_CATEGORY_ID")

type UpdateCategoryInput struct {
	ID        int64
	Name      *string
	Slug      *string
	ParentID  *int64
	ImageURL  *string
	SortOrder *int64
	IsActive  *bool
}

type UpdateCategory struct {
	reader category.Reader
	writer category.Writer
}

func NewUpdateCategory(reader category.Reader, writer category.Writer) *UpdateCategory {
	return &UpdateCategory{reader: reader, writer: writer}
}

func (c *UpdateCategory) Execute(ctx context.Context, input UpdateCategoryInput) (dto.Category, error) {
	if input.ID <= 0 {
		return dto.Category{}, ErrInvalidCategoryID
	}
	item, err := c.reader.FindCategory(ctx, input.ID)
	if err != nil {
		return dto.Category{}, err
	}
	if err := c.applyChanges(ctx, &item, input); err != nil {
		return dto.Category{}, err
	}
	if err := c.writer.UpdateCategory(ctx, item); err != nil {
		return dto.Category{}, err
	}
	return categoryDTO(item), nil
}

func (c *UpdateCategory) applyChanges(ctx context.Context, item *category.Category, input UpdateCategoryInput) error {
	descendants, err := c.parentDescendants(ctx, *item, input.ParentID)
	if err != nil {
		return err
	}
	if parentWasSet(input.ParentID) {
		if err := validateActiveParent(ctx, c.reader, input.ParentID); err != nil {
			return err
		}
	}
	return item.Apply(category.Changes{
		Name: input.Name, Slug: input.Slug, ParentID: input.ParentID,
		ImageURL: input.ImageURL, SortOrder: input.SortOrder, IsActive: input.IsActive,
	}, descendants)
}

func (c *UpdateCategory) parentDescendants(ctx context.Context, item category.Category, parentID *int64) ([]int64, error) {
	if !parentWasSet(parentID) {
		return nil, nil
	}
	return c.reader.DescendantIDs(ctx, item.ID)
}

func parentWasSet(parentID *int64) bool {
	return parentID != nil && *parentID > 0
}
