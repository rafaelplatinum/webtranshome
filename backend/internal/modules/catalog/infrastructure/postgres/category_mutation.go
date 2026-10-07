package postgres

import (
	"context"

	"webtranshome/internal/modules/catalog/domain/category"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func lockCategory(ctx context.Context, session sqlx.Session, categoryID int64) error {
	var id int64
	if err := session.QueryRowCtx(ctx, &id,
		`SELECT id FROM crm_schema.categories WHERE id = $1 FOR UPDATE`, categoryID); err != nil {
		return mapCategoryError(err)
	}
	return nil
}

func checkCategorySlug(ctx context.Context, session sqlx.Session, categoryID int64, slug string) error {
	var exists bool
	if err := session.QueryRowCtx(ctx, &exists,
		`SELECT EXISTS (
			SELECT 1 FROM crm_schema.categories WHERE slug = $1 AND id <> $2
		)`, slug, categoryID); err != nil {
		return err
	}
	if exists {
		return category.ErrCategorySlugConflict
	}
	return nil
}

func updateCategory(ctx context.Context, session sqlx.Session, item category.Category) error {
	parentID, imageURL := nullableCategoryFields(item)
	result, err := session.ExecCtx(ctx,
		`UPDATE crm_schema.categories
		SET name = $2, slug = $3, parent_id = $4, image_url = $5,
			sort_order = $6, is_active = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND NOT EXISTS (
			SELECT 1 FROM crm_schema.categories WHERE slug = $3 AND id <> $1
		)`,
		item.ID, item.Name, item.Slug, parentID, imageURL, item.SortOrder, item.IsActive)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return category.ErrCategorySlugConflict
	}
	return nil
}
