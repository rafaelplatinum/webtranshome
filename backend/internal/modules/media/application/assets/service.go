package assets

import (
	"context"
	"errors"
	"strings"

	"webtranshome/internal/modules/media/domain/asset"
)

type Repository interface {
	asset.Repository
}

type Storage interface {
	asset.FileStorage
}

type Service struct {
	repository Repository
	storage    Storage
}

type Changes struct {
	AltText   *string
	SortOrder *int
	IsActive  *bool
}

func New(repository Repository, storage Storage) *Service {
	return &Service{repository: repository, storage: storage}
}

func (s *Service) List(ctx context.Context, slot string) ([]asset.ImageAsset, error) {
	return s.repository.List(ctx, slot)
}

func (s *Service) Upload(ctx context.Context, slot, alt, originalName string, data []byte, adminID int64, ip string) (asset.ImageAsset, error) {
	alt = strings.TrimSpace(alt)
	if err := asset.ValidateMetadata(slot, alt); err != nil || adminID <= 0 {
		return asset.ImageAsset{}, asset.ErrInvalidAsset
	}
	stored, err := s.storage.Save(originalName, data)
	if err != nil {
		return asset.ImageAsset{}, err
	}
	createdBy := adminID
	created := asset.ImageAsset{
		SlotKey: slot, FilePath: stored.Path, ImageURL: "/media/" + stored.Path, AltText: alt,
		SortOrder: 1, IsActive: true, MIME: stored.MIME, Size: stored.Size,
		Width: stored.Width, Height: stored.Height, CreatedBy: &createdBy,
	}
	item, err := s.repository.Create(ctx, created)
	if err != nil {
		_ = s.storage.Delete(stored.Path)
		return asset.ImageAsset{}, err
	}
	if err := s.repository.RecordActivity(ctx, adminID, "image.create", map[string]any{"slot_key": slot}, map[string]any{"slot_key": slot, "alt_text": alt, "size": stored.Size}, ip); err != nil {
		return item, err
	}
	return item, nil
}

func (s *Service) Update(ctx context.Context, id, adminID int64, changes Changes, ip string) (asset.ImageAsset, error) {
	if id <= 0 || adminID <= 0 || changes.AltText == nil && changes.SortOrder == nil && changes.IsActive == nil {
		return asset.ImageAsset{}, asset.ErrInvalidAsset
	}
	item, err := s.repository.Get(ctx, id)
	if err != nil {
		return asset.ImageAsset{}, err
	}
	before := map[string]any{"alt_text": item.AltText, "sort_order": item.SortOrder, "is_active": item.IsActive}
	if changes.AltText != nil {
		alt := strings.TrimSpace(*changes.AltText)
		if len([]rune(alt)) == 0 || len([]rune(alt)) > 255 {
			return asset.ImageAsset{}, asset.ErrInvalidAsset
		}
		item.AltText = alt
	}
	if changes.SortOrder != nil {
		if *changes.SortOrder < 1 {
			return asset.ImageAsset{}, asset.ErrInvalidAsset
		}
		item.SortOrder = *changes.SortOrder
	}
	if changes.IsActive != nil {
		item.IsActive = *changes.IsActive
	}
	updated, err := s.repository.Update(ctx, item)
	if err != nil {
		return asset.ImageAsset{}, err
	}
	if err := s.repository.RecordActivity(ctx, adminID, "image.update", before, map[string]any{"slot_key": updated.SlotKey, "alt_text": updated.AltText, "sort_order": updated.SortOrder, "is_active": updated.IsActive}, ip); err != nil {
		return updated, err
	}
	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id, adminID int64, ip string) error {
	if id <= 0 || adminID <= 0 {
		return asset.ErrInvalidAsset
	}
	item, err := s.repository.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repository.Delete(ctx, id); err != nil {
		return err
	}
	if err := s.repository.RecordActivity(ctx, adminID, "image.delete", map[string]any{"slot_key": item.SlotKey, "alt_text": item.AltText}, nil, ip); err != nil {
		return err
	}
	return s.storage.Delete(item.FilePath)
}

func (s *Service) Open(filename string) (string, error) {
	path, err := s.storage.Open(filename)
	if err != nil {
		if errors.Is(err, asset.ErrInvalidAsset) {
			return "", err
		}
		return "", err
	}
	return path, nil
}
