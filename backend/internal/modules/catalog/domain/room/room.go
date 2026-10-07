package room

import (
	"errors"
	"strings"
	"unicode/utf8"

	"webtranshome/internal/modules/catalog/domain/slug"
)

var (
	ErrInvalidRoom      = errors.New("CATALOG_INVALID_ROOM")
	ErrRoomNotFound     = errors.New("CATALOG_ROOM_NOT_FOUND")
	ErrRoomSlugConflict = errors.New("CATALOG_ROOM_SLUG_CONFLICT")
)

const maxSortOrder int64 = 1<<31 - 1

type Room struct {
	ID         int64
	Name       string
	Slug       string
	ImageCover string
	SortOrder  int64
	IsActive   bool
}

type Changes struct {
	Name       *string
	Slug       *string
	ImageCover *string
	SortOrder  *int64
	IsActive   *bool
}

func New(name, slug, imageCover string, sortOrder int64) (Room, error) {
	item := Room{
		Name:       strings.TrimSpace(name),
		Slug:       strings.TrimSpace(slug),
		ImageCover: strings.TrimSpace(imageCover),
		SortOrder:  sortOrder,
		IsActive:   true,
	}
	if err := item.Validate(); err != nil {
		return Room{}, err
	}
	return item, nil
}

func (r Room) Validate() error {
	if r.Name == "" || utf8.RuneCountInString(r.Name) > 255 ||
		!slug.IsValid(r.Slug, 255) || r.ImageCover == "" || utf8.RuneCountInString(r.ImageCover) > 255 ||
		r.SortOrder < 0 || r.SortOrder > maxSortOrder {
		return ErrInvalidRoom
	}
	return nil
}

func (r *Room) Apply(changes Changes) error {
	if changes.Empty() {
		return ErrInvalidRoom
	}
	candidate := r.applyChanges(changes)
	if err := candidate.Validate(); err != nil {
		return err
	}
	*r = candidate
	return nil
}

func (c Changes) Empty() bool {
	return c.Name == nil && c.Slug == nil && c.ImageCover == nil &&
		c.SortOrder == nil && c.IsActive == nil
}

func (r Room) applyChanges(changes Changes) Room {
	if changes.Name != nil {
		r.Name = strings.TrimSpace(*changes.Name)
	}
	if changes.Slug != nil {
		r.Slug = strings.TrimSpace(*changes.Slug)
	}
	if changes.ImageCover != nil {
		r.ImageCover = strings.TrimSpace(*changes.ImageCover)
	}
	if changes.SortOrder != nil {
		r.SortOrder = *changes.SortOrder
	}
	if changes.IsActive != nil {
		r.IsActive = *changes.IsActive
	}
	return r
}
