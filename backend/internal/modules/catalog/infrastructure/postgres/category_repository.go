package postgres

import (
	"context"
	"database/sql"
	"errors"

	"webtranshome/internal/modules/catalog/domain/category"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type CategoryRepository struct {
	connection sqlx.SqlConn
}

func NewCategoryRepository(connection sqlx.SqlConn) *CategoryRepository {
	return &CategoryRepository{connection: connection}
}

type categoryRow struct {
	ID        int64          `db:"id"`
	Name      string         `db:"name"`
	Slug      string         `db:"slug"`
	ParentID  sql.NullInt64  `db:"parent_id"`
	ImageURL  sql.NullString `db:"image_url"`
	SortOrder int64          `db:"sort_order"`
	IsActive  bool           `db:"is_active"`
}

func (r *CategoryRepository) ListCategories(ctx context.Context) ([]category.Category, error) {
	var rows []categoryRow
	if err := r.connection.QueryRowsCtx(ctx, &rows,
		`SELECT id, name, slug, parent_id, image_url, sort_order, is_active
		FROM crm_schema.categories
		ORDER BY sort_order, name, id`); err != nil {
		return nil, err
	}
	items := make([]category.Category, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapCategory(row))
	}
	return items, nil
}

func (r *CategoryRepository) FindCategory(ctx context.Context, categoryID int64) (category.Category, error) {
	var row categoryRow
	if err := r.connection.QueryRowCtx(ctx, &row,
		`SELECT id, name, slug, parent_id, image_url, sort_order, is_active
		FROM crm_schema.categories WHERE id = $1`, categoryID); err != nil {
		return category.Category{}, mapCategoryError(err)
	}
	return mapCategory(row), nil
}

func (r *CategoryRepository) DescendantIDs(ctx context.Context, categoryID int64) ([]int64, error) {
	var rows []struct {
		ID int64 `db:"id"`
	}
	const query = `WITH RECURSIVE descendants AS (
		SELECT id FROM crm_schema.categories WHERE parent_id = $1
		UNION
		SELECT child.id
		FROM crm_schema.categories child
		JOIN descendants parent ON child.parent_id = parent.id
	)
	SELECT id FROM descendants`
	if err := r.connection.QueryRowsCtx(ctx, &rows, query, categoryID); err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ID)
	}
	return ids, nil
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, item category.Category) (category.Category, error) {
	parentID, imageURL := nullableCategoryFields(item)
	var id int64
	err := r.connection.QueryRowCtx(ctx, &id,
		`INSERT INTO crm_schema.categories (name, slug, parent_id, image_url, sort_order, is_active)
		VALUES ($1, $2, $3, $4, $5, TRUE)
		ON CONFLICT (slug) DO NOTHING
		RETURNING id`,
		item.Name, item.Slug, parentID, imageURL, item.SortOrder)
	if errors.Is(err, sqlx.ErrNotFound) {
		return category.Category{}, category.ErrCategorySlugConflict
	}
	if err != nil {
		return category.Category{}, err
	}
	item.ID = id
	return item, nil
}

func (r *CategoryRepository) UpdateCategory(ctx context.Context, item category.Category) error {
	return r.connection.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if err := lockCategory(ctx, session, item.ID); err != nil {
			return err
		}
		if err := checkCategorySlug(ctx, session, item.ID, item.Slug); err != nil {
			return err
		}
		return updateCategory(ctx, session, item)
	})
}

func mapCategory(row categoryRow) category.Category {
	return category.Category{
		ID:        row.ID,
		Name:      row.Name,
		Slug:      row.Slug,
		ParentID:  nullableInt64(row.ParentID),
		ImageURL:  nullableString(row.ImageURL),
		SortOrder: row.SortOrder,
		IsActive:  row.IsActive,
	}
}

func mapCategoryError(err error) error {
	if errors.Is(err, sqlx.ErrNotFound) {
		return category.ErrCategoryNotFound
	}
	return err
}

func nullableCategoryFields(item category.Category) (sql.NullInt64, sql.NullString) {
	var parentID sql.NullInt64
	if item.ParentID != nil {
		parentID = sql.NullInt64{Int64: *item.ParentID, Valid: true}
	}
	var imageURL sql.NullString
	if item.ImageURL != nil {
		imageURL = sql.NullString{String: *item.ImageURL, Valid: true}
	}
	return parentID, imageURL
}

func nullableInt64(value sql.NullInt64) *int64 {
	if !value.Valid {
		return nil
	}
	return &value.Int64
}

func nullableString(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}
