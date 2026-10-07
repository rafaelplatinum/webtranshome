package category

import (
	"errors"
	"strings"
	"unicode/utf8"

	"webtranshome/internal/modules/catalog/domain/slug"
)

var (
	ErrInvalidCategory      = errors.New("CATALOG_INVALID_CATEGORY")
	ErrInvalidParent        = errors.New("CATALOG_INVALID_CATEGORY_PARENT")
	ErrCategoryNotFound     = errors.New("CATALOG_CATEGORY_NOT_FOUND")
	ErrCategorySlugConflict = errors.New("CATALOG_CATEGORY_SLUG_CONFLICT")
)

const maxSortOrder int64 = 1<<31 - 1

type Category struct {
	ID        int64
	Name      string
	Slug      string
	ParentID  *int64
	ImageURL  *string
	SortOrder int64
	IsActive  bool
}

type Changes struct {
	Name      *string
	Slug      *string
	ParentID  *int64
	ImageURL  *string
	SortOrder *int64
	IsActive  *bool
}

func New(name, slug string, parentID *int64, imageURL string, sortOrder int64) (Category, error) {
	category := Category{
		Name:      strings.TrimSpace(name),
		Slug:      strings.TrimSpace(slug),
		ParentID:  copyInt64(parentID),
		ImageURL:  optionalString(imageURL),
		SortOrder: sortOrder,
		IsActive:  true,
	}
	if err := category.Validate(); err != nil {
		return Category{}, err
	}
	return category, nil
}

func (c Category) Validate() error {
	if c.Name == "" || utf8.RuneCountInString(c.Name) > 255 ||
		!slug.IsValid(c.Slug, 255) || c.SortOrder < 0 || c.SortOrder > maxSortOrder {
		return ErrInvalidCategory
	}
	if c.ParentID != nil && *c.ParentID <= 0 {
		return ErrInvalidParent
	}
	if c.ImageURL != nil && utf8.RuneCountInString(*c.ImageURL) > 255 {
		return ErrInvalidCategory
	}
	return nil
}

func (c *Category) Apply(changes Changes, descendants []int64) error {
	if changes.Empty() {
		return ErrInvalidCategory
	}
	candidate := c.applyChanges(changes)
	if err := candidate.Validate(); err != nil {
		return err
	}
	if candidate.ParentID != nil && (*candidate.ParentID == c.ID || contains(descendants, *candidate.ParentID)) {
		return ErrInvalidParent
	}
	*c = candidate
	return nil
}

func (c Changes) Empty() bool {
	return c.Name == nil && c.Slug == nil && c.ParentID == nil &&
		c.ImageURL == nil && c.SortOrder == nil && c.IsActive == nil
}

func (c Category) applyChanges(changes Changes) Category {
	candidate := c
	if changes.Name != nil {
		candidate.Name = strings.TrimSpace(*changes.Name)
	}
	if changes.Slug != nil {
		candidate.Slug = strings.TrimSpace(*changes.Slug)
	}
	if changes.ParentID != nil {
		if *changes.ParentID == 0 {
			candidate.ParentID = nil
		} else {
			candidate.ParentID = copyInt64(changes.ParentID)
		}
	}
	if changes.ImageURL != nil {
		candidate.ImageURL = optionalString(*changes.ImageURL)
	}
	if changes.SortOrder != nil {
		candidate.SortOrder = *changes.SortOrder
	}
	if changes.IsActive != nil {
		candidate.IsActive = *changes.IsActive
	}
	return candidate
}

func contains(values []int64, expected int64) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}

func copyInt64(value *int64) *int64 {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func optionalString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}
