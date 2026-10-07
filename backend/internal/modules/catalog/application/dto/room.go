package dto

type Room struct {
	ID         int64
	Name       string
	Slug       string
	ImageCover string
	SortOrder  int64
	IsActive   bool
}
