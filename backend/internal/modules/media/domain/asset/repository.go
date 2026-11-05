package asset

import "context"

type Repository interface {
	List(context.Context, string) ([]ImageAsset, error)
	Get(context.Context, int64) (ImageAsset, error)
	Create(context.Context, ImageAsset) (ImageAsset, error)
	Update(context.Context, ImageAsset) (ImageAsset, error)
	Delete(context.Context, int64) error
	RecordActivity(context.Context, int64, string, map[string]any, map[string]any, string) error
}

type FileStorage interface {
	Save(string, []byte) (StoredFile, error)
	Delete(string) error
	Open(string) (string, error)
}

type StoredFile struct {
	Path   string
	MIME   string
	Size   int64
	Width  int
	Height int
}
