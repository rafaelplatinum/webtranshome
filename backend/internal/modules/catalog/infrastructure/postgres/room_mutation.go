package postgres

import (
	"context"

	"webtranshome/internal/modules/catalog/domain/room"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func lockRoom(ctx context.Context, session sqlx.Session, roomID int64) error {
	var id int64
	if err := session.QueryRowCtx(ctx, &id,
		`SELECT id FROM crm_schema.rooms WHERE id = $1 FOR UPDATE`, roomID); err != nil {
		return mapRoomError(err)
	}
	return nil
}

func checkRoomSlug(ctx context.Context, session sqlx.Session, roomID int64, slug string) error {
	var exists bool
	if err := session.QueryRowCtx(ctx, &exists,
		`SELECT EXISTS (
			SELECT 1 FROM crm_schema.rooms WHERE slug = $1 AND id <> $2
		)`, slug, roomID); err != nil {
		return err
	}
	if exists {
		return room.ErrRoomSlugConflict
	}
	return nil
}

func updateRoom(ctx context.Context, session sqlx.Session, item room.Room) error {
	result, err := session.ExecCtx(ctx,
		`UPDATE crm_schema.rooms
		SET name = $2, slug = $3, image_cover = $4, sort_order = $5,
			is_active = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND NOT EXISTS (
			SELECT 1 FROM crm_schema.rooms WHERE slug = $3 AND id <> $1
		)`,
		item.ID, item.Name, item.Slug, item.ImageCover, item.SortOrder, item.IsActive)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return room.ErrRoomSlugConflict
	}
	return nil
}
