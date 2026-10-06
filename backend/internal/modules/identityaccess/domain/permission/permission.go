package permission

type Permission struct {
	ID          int64
	MenuID      int64
	Code        string
	Action      string
	Description *string
}
