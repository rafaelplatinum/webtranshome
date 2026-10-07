package permission

const (
	AccessControlManageCode   = "rbac.manage"
	CatalogBrandManageCode    = "catalog.brand.manage"
	CatalogCategoryManageCode = "catalog.category.manage"
	CatalogRoomManageCode     = "catalog.room.manage"
)

type Permission struct {
	ID          int64
	MenuID      int64
	Code        string
	Action      string
	Description *string
}
