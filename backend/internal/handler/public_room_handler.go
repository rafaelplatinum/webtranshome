package handler

import (
	"net/http"

	"webtranshome/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/postgres"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type publicRoomItem struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	ImageCover   string `json:"image_cover"`
	ProductCount int64  `json:"product_count"`
}

type publicRoomListResponse struct {
	Data []publicRoomItem `json:"data"`
}

type publicRoomDetailResponse struct {
	Data publicRoomItem `json:"data"`
}

func PublicRoomsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	conn := postgres.New(svcCtx.Config.Database.DataSource)
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := loadPublicRooms(r, svcCtx, conn)
		if err != nil {
			logx.WithContext(r.Context()).Errorf("list public rooms: %v", err)
			httpx.WriteJsonCtx(r.Context(), w, http.StatusInternalServerError, map[string]string{"message": "Gagal memuat ruangan."})
			return
		}
		httpx.OkJsonCtx(r.Context(), w, publicRoomListResponse{Data: items})
	}
}

func PublicRoomHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	conn := postgres.New(svcCtx.Config.Database.DataSource)
	return func(w http.ResponseWriter, r *http.Request) {
		var path struct {
			Slug string `path:"slug"`
		}
		if err := httpx.ParsePath(r, &path); err != nil || path.Slug == "" {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusBadRequest, map[string]string{"message": "Slug ruangan tidak valid."})
			return
		}

		items, err := loadPublicRooms(r, svcCtx, conn)
		if err != nil {
			logx.WithContext(r.Context()).Errorf("get public room: %v", err)
			httpx.WriteJsonCtx(r.Context(), w, http.StatusInternalServerError, map[string]string{"message": "Gagal memuat ruangan."})
			return
		}
		for _, item := range items {
			if item.Slug == path.Slug {
				httpx.OkJsonCtx(r.Context(), w, publicRoomDetailResponse{Data: item})
				return
			}
		}
		httpx.WriteJsonCtx(r.Context(), w, http.StatusNotFound, map[string]string{"message": "Ruangan tidak ditemukan."})
	}
}

func loadPublicRooms(r *http.Request, svcCtx *svc.ServiceContext, conn sqlx.SqlConn) ([]publicRoomItem, error) {
	rooms, err := svcCtx.ListRooms.Execute(r.Context())
	if err != nil {
		return nil, err
	}

	// A single aggregation query avoids doing one COUNT per room.
	var rows []struct {
		RoomID       int64 `db:"room_id"`
		ProductCount int64 `db:"product_count"`
	}
	err = conn.QueryRowsCtx(r.Context(), &rows,
		`SELECT pr.room_id, COUNT(*)::BIGINT AS product_count
         FROM crm_schema.product_rooms pr
         JOIN crm_schema.products p ON p.id = pr.product_id AND p.is_active = TRUE
         GROUP BY pr.room_id`)
	if err != nil {
		return nil, err
	}
	counts := make(map[int64]int64, len(rows))
	for _, row := range rows {
		counts[row.RoomID] = row.ProductCount
	}

	result := make([]publicRoomItem, 0, len(rooms))
	for _, room := range rooms {
		if !room.IsActive {
			continue
		}
		result = append(result, publicRoomItem{
			ID: room.ID, Name: room.Name, Slug: room.Slug,
			ImageCover: room.ImageCover, ProductCount: counts[room.ID],
		})
	}
	return result, nil
}