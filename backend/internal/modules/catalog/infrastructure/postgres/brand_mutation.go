package postgres

import (
	"context"

	"webtranshome/internal/modules/catalog/domain/brand"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func lockBrand(ctx context.Context, session sqlx.Session, brandID int64) error {
	var id int64
	if err := session.QueryRowCtx(ctx, &id,
		`SELECT id FROM crm_schema.brands WHERE id = $1 FOR UPDATE`, brandID); err != nil {
		return mapBrandError(err)
	}
	return nil
}

func checkBrandSlug(ctx context.Context, session sqlx.Session, brandID int64, slug string) error {
	var exists bool
	if err := session.QueryRowCtx(ctx, &exists,
		`SELECT EXISTS (
			SELECT 1 FROM crm_schema.brands WHERE slug = $1 AND id <> $2
		)`, slug, brandID); err != nil {
		return err
	}
	if exists {
		return brand.ErrBrandSlugConflict
	}
	return nil
}

func updateBrand(ctx context.Context, session sqlx.Session, item brand.Brand) error {
	result, err := session.ExecCtx(ctx,
		`UPDATE crm_schema.brands
		SET name = $2, slug = $3, logo_url = $4, is_active = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND NOT EXISTS (
			SELECT 1 FROM crm_schema.brands WHERE slug = $3 AND id <> $1
		)`,
		item.ID, item.Name, item.Slug, nullableBrandLogo(item.LogoURL), item.IsActive)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return brand.ErrBrandSlugConflict
	}
	return nil
}
