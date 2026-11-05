package asset

import (
	"errors"
	"strings"
)

var (
	ErrInvalidAsset = errors.New("invalid image asset")
	ErrNotFound     = errors.New("image asset not found")
	ErrInvalidImage = errors.New("image must be a valid JPEG, PNG, or WebP file of at least 64x64 pixels")
)

var AllowedSlots = map[string]struct{}{
	"home.hero": {}, "home.promo_banner": {}, "home.side_banner": {}, "category.tile": {},
	"product.main": {}, "product.gallery": {}, "brand.logo": {}, "room.cover": {},
	"service.cover": {}, "article.cover": {}, "promo.cover": {}, "event.cover": {},
	"store.front": {}, "store.map": {}, "site.logo": {}, "site.favicon": {}, "site.og_image": {},
}

type ImageAsset struct {
	ID        int64   `json:"id"`
	SlotKey   string  `json:"slot_key"`
	OwnerType *string `json:"owner_type,omitempty"`
	OwnerID   *int64  `json:"owner_id,omitempty"`
	FilePath  string  `json:"file_path"`
	ImageURL  string  `json:"image_url"`
	AltText   string  `json:"alt_text"`
	SortOrder int     `json:"sort_order"`
	IsActive  bool    `json:"is_active"`
	MIME      string  `json:"mime"`
	Size      int64   `json:"size"`
	Width     int     `json:"width"`
	Height    int     `json:"height"`
	IsDummy   bool    `json:"is_dummy"`
	CreatedBy *int64  `json:"created_by,omitempty"`
	CreatedAt string  `json:"created_at,omitempty"`
	UpdatedAt string  `json:"updated_at,omitempty"`
}

func ValidateMetadata(slot, alt string) error {
	if _, ok := AllowedSlots[slot]; !ok || len(strings.TrimSpace(alt)) == 0 || len([]rune(strings.TrimSpace(alt))) > 255 {
		return ErrInvalidAsset
	}
	return nil
}
