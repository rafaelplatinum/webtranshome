package postgres

import (
	"context"
	"database/sql"
	"errors"

	"webtranshome/internal/modules/catalog/domain/brand"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type BrandRepository struct {
	connection sqlx.SqlConn
}

func NewBrandRepository(connection sqlx.SqlConn) *BrandRepository {
	return &BrandRepository{connection: connection}
}

type brandRow struct {
	ID       int64          `db:"id"`
	Name     string         `db:"name"`
	Slug     string         `db:"slug"`
	LogoURL  sql.NullString `db:"logo_url"`
	IsActive bool           `db:"is_active"`
}

func (r *BrandRepository) ListBrands(ctx context.Context) ([]brand.Brand, error) {
	var rows []brandRow
	if err := r.connection.QueryRowsCtx(ctx, &rows,
		`SELECT id, name, slug, logo_url, is_active
		FROM crm_schema.brands
		ORDER BY name, id`); err != nil {
		return nil, err
	}
	items := make([]brand.Brand, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapBrand(row))
	}
	return items, nil
}

func (r *BrandRepository) FindBrand(ctx context.Context, brandID int64) (brand.Brand, error) {
	var row brandRow
	if err := r.connection.QueryRowCtx(ctx, &row,
		`SELECT id, name, slug, logo_url, is_active
		FROM crm_schema.brands WHERE id = $1`, brandID); err != nil {
		return brand.Brand{}, mapBrandError(err)
	}
	return mapBrand(row), nil
}

func (r *BrandRepository) CreateBrand(ctx context.Context, item brand.Brand) (brand.Brand, error) {
	var id int64
	err := r.connection.QueryRowCtx(ctx, &id,
		`INSERT INTO crm_schema.brands (name, slug, logo_url, is_active)
		VALUES ($1, $2, $3, TRUE)
		ON CONFLICT (slug) DO NOTHING
		RETURNING id`,
		item.Name, item.Slug, nullableBrandLogo(item.LogoURL))
	if errors.Is(err, sqlx.ErrNotFound) {
		return brand.Brand{}, brand.ErrBrandSlugConflict
	}
	if err != nil {
		return brand.Brand{}, err
	}
	item.ID = id
	return item, nil
}

func (r *BrandRepository) UpdateBrand(ctx context.Context, item brand.Brand) error {
	return r.connection.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if err := lockBrand(ctx, session, item.ID); err != nil {
			return err
		}
		if err := checkBrandSlug(ctx, session, item.ID, item.Slug); err != nil {
			return err
		}
		return updateBrand(ctx, session, item)
	})
}

func mapBrand(row brandRow) brand.Brand {
	return brand.Brand{
		ID: row.ID, Name: row.Name, Slug: row.Slug,
		LogoURL: nullableString(row.LogoURL), IsActive: row.IsActive,
	}
}

func mapBrandError(err error) error {
	if errors.Is(err, sqlx.ErrNotFound) {
		return brand.ErrBrandNotFound
	}
	return err
}

func nullableBrandLogo(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *value, Valid: true}
}

var (
	_ brand.Reader = (*BrandRepository)(nil)
	_ brand.Writer = (*BrandRepository)(nil)
)
