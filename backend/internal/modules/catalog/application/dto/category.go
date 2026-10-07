package dto

type Category struct {
	ID        int64
	Name      string
	Slug      string
	ParentID  *int64
	ImageURL  string
	SortOrder int64
	IsActive  bool
}
