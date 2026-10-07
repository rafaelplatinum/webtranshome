package brand

import (
	"errors"
	"strings"
	"unicode/utf8"

	"webtranshome/internal/modules/catalog/domain/slug"
)

var (
	ErrInvalidBrand      = errors.New("CATALOG_INVALID_BRAND")
	ErrBrandNotFound     = errors.New("CATALOG_BRAND_NOT_FOUND")
	ErrBrandSlugConflict = errors.New("CATALOG_BRAND_SLUG_CONFLICT")
)

type Brand struct {
	ID       int64
	Name     string
	Slug     string
	LogoURL  *string
	IsActive bool
}

type Changes struct {
	Name     *string
	Slug     *string
	LogoURL  *string
	IsActive *bool
}

func New(name, slug, logoURL string) (Brand, error) {
	item := Brand{
		Name:     strings.TrimSpace(name),
		Slug:     strings.TrimSpace(slug),
		LogoURL:  optionalString(logoURL),
		IsActive: true,
	}
	if err := item.Validate(); err != nil {
		return Brand{}, err
	}
	return item, nil
}

func (b Brand) Validate() error {
	if b.Name == "" || utf8.RuneCountInString(b.Name) > 150 ||
		!slug.IsValid(b.Slug, 150) || (b.LogoURL != nil && utf8.RuneCountInString(*b.LogoURL) > 255) {
		return ErrInvalidBrand
	}
	return nil
}

func (b *Brand) Apply(changes Changes) error {
	if changes.Empty() {
		return ErrInvalidBrand
	}
	candidate := b.applyChanges(changes)
	if err := candidate.Validate(); err != nil {
		return err
	}
	*b = candidate
	return nil
}

func (c Changes) Empty() bool {
	return c.Name == nil && c.Slug == nil && c.LogoURL == nil && c.IsActive == nil
}

func (b Brand) applyChanges(changes Changes) Brand {
	if changes.Name != nil {
		b.Name = strings.TrimSpace(*changes.Name)
	}
	if changes.Slug != nil {
		b.Slug = strings.TrimSpace(*changes.Slug)
	}
	if changes.LogoURL != nil {
		b.LogoURL = optionalString(*changes.LogoURL)
	}
	if changes.IsActive != nil {
		b.IsActive = *changes.IsActive
	}
	return b
}

func optionalString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}
