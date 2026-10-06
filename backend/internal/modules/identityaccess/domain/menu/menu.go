package menu

type Menu struct {
	ID        int64
	ParentID  *int64
	Name      string
	Code      string
	Route     *string
	Icon      *string
	SortOrder int
	IsActive  bool
}
