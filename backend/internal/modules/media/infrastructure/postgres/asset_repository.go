package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"webtranshome/internal/modules/media/domain/asset"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type AssetRepository struct{ connection sqlx.SqlConn }

func NewAssetRepository(connection sqlx.SqlConn) *AssetRepository {
	return &AssetRepository{connection: connection}
}

type assetRow struct {
	ID        int64          `db:"id"`
	SlotKey   string         `db:"slot_key"`
	OwnerType sql.NullString `db:"owner_type"`
	OwnerID   sql.NullInt64  `db:"owner_id"`
	FilePath  string         `db:"file_path"`
	AltText   string         `db:"alt_text"`
	SortOrder int            `db:"sort_order"`
	IsActive  bool           `db:"is_active"`
	MIME      string         `db:"mime"`
	Size      int64          `db:"size"`
	Width     int            `db:"width"`
	Height    int            `db:"height"`
	IsDummy   bool           `db:"is_dummy"`
	CreatedBy sql.NullInt64  `db:"created_by"`
	CreatedAt sql.NullTime   `db:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
}

const assetColumns = `id, slot_key, owner_type, owner_id, file_path, alt_text, sort_order, is_active, mime, size, width, height, is_dummy, created_by, created_at, updated_at`

func (r *AssetRepository) List(ctx context.Context, slot string) ([]asset.ImageAsset, error) {
	var rows []assetRow
	query := `SELECT ` + assetColumns + ` FROM crm_schema.image_assets`
	args := []any{}
	if slot != "" {
		query += ` WHERE slot_key = $1`
		args = append(args, slot)
	}
	query += ` ORDER BY slot_key, sort_order, id`
	if err := r.connection.QueryRowsCtx(ctx, &rows, query, args...); err != nil {
		return nil, err
	}
	items := make([]asset.ImageAsset, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapAsset(row))
	}
	return items, nil
}

func (r *AssetRepository) Get(ctx context.Context, id int64) (asset.ImageAsset, error) {
	var row assetRow
	err := r.connection.QueryRowCtx(ctx, &row, `SELECT `+assetColumns+` FROM crm_schema.image_assets WHERE id = $1`, id)
	if errors.Is(err, sqlx.ErrNotFound) {
		return asset.ImageAsset{}, asset.ErrNotFound
	}
	if err != nil {
		return asset.ImageAsset{}, err
	}
	return mapAsset(row), nil
}

func (r *AssetRepository) Create(ctx context.Context, item asset.ImageAsset) (asset.ImageAsset, error) {
	var id int64
	err := r.connection.QueryRowCtx(ctx, &id, `INSERT INTO crm_schema.image_assets
		(slot_key, owner_type, owner_id, file_path, alt_text, sort_order, is_active, mime, size, width, height, is_dummy, created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING id`,
		item.SlotKey, item.OwnerType, item.OwnerID, item.FilePath, item.AltText, item.SortOrder, item.IsActive,
		item.MIME, item.Size, item.Width, item.Height, item.IsDummy, item.CreatedBy)
	if err != nil {
		return asset.ImageAsset{}, err
	}
	return r.Get(ctx, id)
}

func (r *AssetRepository) Update(ctx context.Context, item asset.ImageAsset) (asset.ImageAsset, error) {
	result, err := r.connection.ExecCtx(ctx, `UPDATE crm_schema.image_assets SET alt_text=$2, sort_order=$3, is_active=$4 WHERE id=$1`, item.ID, item.AltText, item.SortOrder, item.IsActive)
	if err != nil {
		return asset.ImageAsset{}, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return asset.ImageAsset{}, err
	}
	if count == 0 {
		return asset.ImageAsset{}, asset.ErrNotFound
	}
	return r.Get(ctx, item.ID)
}

func (r *AssetRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.connection.ExecCtx(ctx, `DELETE FROM crm_schema.image_assets WHERE id=$1`, id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return asset.ErrNotFound
	}
	return nil
}

func (r *AssetRepository) RecordActivity(ctx context.Context, adminID int64, action string, before, after map[string]any, ip string) error {
	var beforeJSON, afterJSON any
	if before != nil {
		value, err := json.Marshal(before)
		if err != nil {
			return err
		}
		beforeJSON = string(value)
	}
	if after != nil {
		value, err := json.Marshal(after)
		if err != nil {
			return err
		}
		afterJSON = string(value)
	}
	_, err := r.connection.ExecCtx(ctx, `INSERT INTO crm_schema.activity_logs (admin_id, action, module, before_data, after_data, ip_address)
		VALUES ($1,$2,'Pengelolaan gambar',$3,$4,$5)`, adminID, action, beforeJSON, afterJSON, ip)
	return err
}

func mapAsset(row assetRow) asset.ImageAsset {
	item := asset.ImageAsset{
		ID: row.ID, SlotKey: row.SlotKey, FilePath: row.FilePath, ImageURL: "/media/" + row.FilePath,
		AltText: row.AltText, SortOrder: row.SortOrder, IsActive: row.IsActive, MIME: row.MIME,
		Size: row.Size, Width: row.Width, Height: row.Height, IsDummy: row.IsDummy,
	}
	if row.OwnerType.Valid {
		item.OwnerType = &row.OwnerType.String
	}
	if row.OwnerID.Valid {
		item.OwnerID = &row.OwnerID.Int64
	}
	if row.CreatedBy.Valid {
		item.CreatedBy = &row.CreatedBy.Int64
	}
	if row.CreatedAt.Valid {
		item.CreatedAt = row.CreatedAt.Time.Format(time.RFC3339)
	}
	if row.UpdatedAt.Valid {
		item.UpdatedAt = row.UpdatedAt.Time.Format(time.RFC3339)
	}
	return item
}
