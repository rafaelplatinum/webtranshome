package postgres

import (
	"context"
	"errors"

	"webtranshome/internal/modules/catalog/domain/room"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type RoomRepository struct {
	connection sqlx.SqlConn
}

func NewRoomRepository(connection sqlx.SqlConn) *RoomRepository {
	return &RoomRepository{connection: connection}
}

type roomRow struct {
	ID         int64  `db:"id"`
	Name       string `db:"name"`
	Slug       string `db:"slug"`
	ImageCover string `db:"image_cover"`
	SortOrder  int64  `db:"sort_order"`
	IsActive   bool   `db:"is_active"`
}

func (r *RoomRepository) ListRooms(ctx context.Context) ([]room.Room, error) {
	var rows []roomRow
	if err := r.connection.QueryRowsCtx(ctx, &rows,
		`SELECT id, name, slug, image_cover, sort_order, is_active
		FROM crm_schema.rooms
		ORDER BY sort_order, name, id`); err != nil {
		return nil, err
	}
	items := make([]room.Room, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapRoom(row))
	}
	return items, nil
}

func (r *RoomRepository) FindRoom(ctx context.Context, roomID int64) (room.Room, error) {
	var row roomRow
	if err := r.connection.QueryRowCtx(ctx, &row,
		`SELECT id, name, slug, image_cover, sort_order, is_active
		FROM crm_schema.rooms WHERE id = $1`, roomID); err != nil {
		return room.Room{}, mapRoomError(err)
	}
	return mapRoom(row), nil
}

func (r *RoomRepository) CreateRoom(ctx context.Context, item room.Room) (room.Room, error) {
	var id int64
	err := r.connection.QueryRowCtx(ctx, &id,
		`INSERT INTO crm_schema.rooms (name, slug, image_cover, sort_order, is_active)
		VALUES ($1, $2, $3, $4, TRUE)
		ON CONFLICT (slug) DO NOTHING
		RETURNING id`,
		item.Name, item.Slug, item.ImageCover, item.SortOrder)
	if errors.Is(err, sqlx.ErrNotFound) {
		return room.Room{}, room.ErrRoomSlugConflict
	}
	if err != nil {
		return room.Room{}, err
	}
	item.ID = id
	return item, nil
}

func (r *RoomRepository) UpdateRoom(ctx context.Context, item room.Room) error {
	return r.connection.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if err := lockRoom(ctx, session, item.ID); err != nil {
			return err
		}
		if err := checkRoomSlug(ctx, session, item.ID, item.Slug); err != nil {
			return err
		}
		return updateRoom(ctx, session, item)
	})
}

func mapRoom(row roomRow) room.Room {
	return room.Room{
		ID: row.ID, Name: row.Name, Slug: row.Slug,
		ImageCover: row.ImageCover, SortOrder: row.SortOrder, IsActive: row.IsActive,
	}
}

func mapRoomError(err error) error {
	if errors.Is(err, sqlx.ErrNotFound) {
		return room.ErrRoomNotFound
	}
	return err
}

var (
	_ room.Reader = (*RoomRepository)(nil)
	_ room.Writer = (*RoomRepository)(nil)
)
