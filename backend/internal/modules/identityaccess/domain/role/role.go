package role

type Role struct {
	ID          int64
	Code        string
	Name        string
	Description *string
	IsActive    bool
}
