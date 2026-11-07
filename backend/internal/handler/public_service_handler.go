package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"webtranshome/internal/svc"

	"github.com/zeromicro/go-zero/core/stores/postgres"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type publicServiceRow struct {
	ID          int64          `db:"id"`
	Name        string         `db:"name"`
	Slug        string         `db:"slug"`
	Summary     string         `db:"summary"`
	Description sql.NullString `db:"description"`
	ImageURL    sql.NullString `db:"image_url"`
	ServiceArea sql.NullString `db:"service_area"`
}

type publicServiceItem struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Summary     string  `json:"summary"`
	ImageURL    *string `json:"image_url"`
	ServiceArea *string `json:"service_area"`
}

type publicServiceCategory struct {
	Name string `json:"name" db:"name"`
	Slug string `json:"slug" db:"slug"`
}

type publicServiceDetail struct {
	publicServiceItem
	Description string                  `json:"description"`
	Categories  []publicServiceCategory `json:"categories"`
}

type publicServiceListResponse struct {
	Data []publicServiceItem `json:"data"`
}
type publicServiceDetailResponse struct {
	Data publicServiceDetail `json:"data"`
}

func publicServiceFromRow(r *http.Request, v publicServiceRow) publicServiceItem {
	photo := publicSQLString(v.ImageURL)
	if photo != nil {
		url := publicProductMediaURL(r, *photo)
		photo = &url
	}
	return publicServiceItem{ID: v.ID, Name: v.Name, Slug: v.Slug, Summary: v.Summary, ImageURL: photo, ServiceArea: publicSQLString(v.ServiceArea)}
}

func PublicServicesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	db := postgres.New(svcCtx.Config.Database.DataSource)
	return func(w http.ResponseWriter, r *http.Request) {
		category := strings.TrimSpace(r.URL.Query().Get("category"))
		limit := 0
		if value := strings.TrimSpace(r.URL.Query().Get("limit")); value != "" {
			parsed, err := strconv.Atoi(value)
			if err != nil || parsed < 1 || parsed > 100 {
				httpx.WriteJsonCtx(r.Context(), w, http.StatusBadRequest, map[string]string{"message": "Limit jasa harus 1 sampai 100."})
				return
			}
			limit = parsed
		}
		// Unknown category naturally returns an empty list; if a child is
		// selected, services attached to its active parent also match.
		query := `SELECT s.id,s.name,s.slug,s.summary,s.description,s.image_url,s.service_area
             FROM crm_schema.services s
             WHERE s.is_active=TRUE
               AND ($1='' OR EXISTS (
                   SELECT 1 FROM crm_schema.categories c
                   JOIN crm_schema.service_categories sc ON sc.category_id IN (c.id,c.parent_id)
                   JOIN crm_schema.categories attached ON attached.id=sc.category_id AND attached.is_active=TRUE
                   WHERE c.slug=$1 AND c.is_active=TRUE AND sc.service_id=s.id
               ))
             ORDER BY s.sort_order ASC,s.id ASC`
		if limit > 0 {
			query += " LIMIT " + strconv.Itoa(limit)
		}
		var rows []publicServiceRow
		if err := db.QueryRowsCtx(r.Context(), &rows, query, category); err != nil {
			publicContentError(w, r, "jasa", err)
			return
		}
		items := make([]publicServiceItem, 0, len(rows))
		for _, row := range rows {
			items = append(items, publicServiceFromRow(r, row))
		}
		httpx.OkJsonCtx(r.Context(), w, publicServiceListResponse{Data: items})
	}
}

func PublicServiceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	db := postgres.New(svcCtx.Config.Database.DataSource)
	return func(w http.ResponseWriter, r *http.Request) {
		var path struct {
			Slug string `path:"slug"`
		}
		if err := httpx.ParsePath(r, &path); err != nil || strings.TrimSpace(path.Slug) == "" {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusBadRequest, map[string]string{"message": "Slug jasa tidak valid."})
			return
		}
		var row publicServiceRow
		err := db.QueryRowCtx(r.Context(), &row, `SELECT id,name,slug,summary,description,image_url,service_area
              FROM crm_schema.services WHERE slug=$1 AND is_active=TRUE`, path.Slug)
		if errors.Is(err, sqlx.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusNotFound, map[string]string{"message": "Jasa tidak ditemukan."})
			return
		}
		if err != nil {
			publicContentError(w, r, "jasa", err)
			return
		}
		var categories []publicServiceCategory
		err = db.QueryRowsCtx(r.Context(), &categories, `SELECT c.name,c.slug FROM crm_schema.categories c
            JOIN crm_schema.service_categories sc ON sc.category_id=c.id
            WHERE sc.service_id=$1 AND c.is_active=TRUE
            ORDER BY c.sort_order,c.name,c.id`, row.ID)
		if err != nil {
			publicContentError(w, r, "jasa", err)
			return
		}
		if categories == nil {
			categories = []publicServiceCategory{}
		}
		description := ""
		if row.Description.Valid {
			description = row.Description.String
		}
		httpx.OkJsonCtx(r.Context(), w, publicServiceDetailResponse{Data: publicServiceDetail{
			publicServiceItem: publicServiceFromRow(r, row), Description: description, Categories: categories,
		}})
	}
}