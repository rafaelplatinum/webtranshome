package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"webtranshome/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/postgres"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// All queries in this file read the public projection only: no drafts or
// unpublished/scheduled banners are returned.
type publicBannerRow struct {
	ID        int64          `db:"id"`
	Title     string         `db:"title"`
	ImageURL  string         `db:"image_url"`
	LinkURL   sql.NullString `db:"link_url"`
	Position  string         `db:"position"`
	StartAt   sql.NullString `db:"start_at"`
	EndAt     sql.NullString `db:"end_at"`
	Subtitle  sql.NullString `db:"subtitle"`
	VideoURL  sql.NullString `db:"video_url"`
	CTALabel  sql.NullString `db:"cta_label"`
	CTA2Label sql.NullString `db:"cta2_label"`
	CTA2URL   sql.NullString `db:"cta2_url"`
}

type publicBannerItem struct {
	ID        int64   `json:"id"`
	Title     string  `json:"title"`
	ImageURL  string  `json:"image_url"`
	LinkURL   *string `json:"link_url"`
	Position  string  `json:"position"`
	StartAt   *string `json:"start_at"`
	EndAt     *string `json:"end_at"`
	Subtitle  *string `json:"subtitle,omitempty"`
	VideoURL  *string `json:"video_url,omitempty"`
	CTALabel  *string `json:"cta_label,omitempty"`
	CTA2Label *string `json:"cta2_label,omitempty"`
	CTA2URL   *string `json:"cta2_url,omitempty"`
}

type publicBannerListResponse struct {
	Data []publicBannerItem `json:"data"`
}

func publicSQLString(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	result := v.String
	return &result
}

func PublicBannersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	db := postgres.New(svcCtx.Config.Database.DataSource)
	return func(w http.ResponseWriter, r *http.Request) {
		position := strings.TrimSpace(r.URL.Query().Get("position"))
		switch position {
		case "", "HOME_HERO", "HOME_SLIDER", "HOME_SIDE":
		default:
			httpx.WriteJsonCtx(r.Context(), w, http.StatusBadRequest, map[string]string{"message": "Posisi banner tidak valid."})
			return
		}
		var rows []publicBannerRow
		err := db.QueryRowsCtx(r.Context(), &rows, `
            SELECT id, title, image_url, link_url, position,
              to_char(start_at, 'YYYY-MM-DD"T"HH24:MI:SS') AS start_at,
              to_char(end_at, 'YYYY-MM-DD"T"HH24:MI:SS') AS end_at,
              subtitle, video_url, cta_label, cta2_label, cta2_url
            FROM crm_schema.banners
            WHERE is_active = TRUE
              AND (start_at IS NULL OR start_at <= timezone('Asia/Jakarta', now()))
              AND (end_at IS NULL OR end_at >= timezone('Asia/Jakarta', now()))
              AND ($1 = '' OR position = $1)
            ORDER BY CASE position WHEN 'HOME_HERO' THEN 0 WHEN 'HOME_SLIDER' THEN 1 WHEN 'HOME_SIDE' THEN 2 ELSE 3 END,
                     sort_order ASC, id ASC`, position)
		if err != nil {
			publicContentError(w, r, "banner", err)
			return
		}
		result := make([]publicBannerItem, 0, len(rows))
		for _, b := range rows {
			item := publicBannerItem{
				ID: b.ID, Title: b.Title, ImageURL: publicProductMediaURL(r, b.ImageURL),
				LinkURL: publicSQLString(b.LinkURL), Position: b.Position,
				StartAt: publicSQLString(b.StartAt), EndAt: publicSQLString(b.EndAt),
			}
			if b.Position == "HOME_HERO" {
				item.Subtitle = publicSQLString(b.Subtitle)
				item.VideoURL = publicSQLString(b.VideoURL)
				if item.VideoURL != nil {
					value := publicProductMediaURL(r, *item.VideoURL)
					item.VideoURL = &value
				}
				item.CTALabel = publicSQLString(b.CTALabel)
				item.CTA2Label = publicSQLString(b.CTA2Label)
				item.CTA2URL = publicSQLString(b.CTA2URL)
			}
			result = append(result, item)
		}
		httpx.OkJsonCtx(r.Context(), w, publicBannerListResponse{Data: result})
	}
}

type publicArticleRow struct {
	ID           int64          `db:"id"`
	Type         string         `db:"type"`
	Title        string         `db:"title"`
	Slug         string         `db:"slug"`
	ThumbnailURL sql.NullString `db:"thumbnail_url"`
	Content      string         `db:"content"`
	Status       string         `db:"status"`
	PublishedAt  sql.NullString `db:"published_at"`
}

type publicArticleItem struct {
	ID           int64   `json:"id"`
	Type         string  `json:"type"`
	Title        string  `json:"title"`
	Slug         string  `json:"slug"`
	ThumbnailURL *string `json:"thumbnail_url"`
	Excerpt      string  `json:"excerpt"`
	PublishedAt  *string `json:"published_at"`
}

type publicArticleDetail struct {
	publicArticleItem
	Content string `json:"content"`
	Status  string `json:"status"`
}

type publicArticleMeta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type publicArticleListResponse struct {
	Data []publicArticleItem `json:"data"`
	Meta publicArticleMeta   `json:"meta"`
}

type publicArticleDetailResponse struct {
	Data publicArticleDetail `json:"data"`
}

var articleWhitespace = regexp.MustCompile(`\s+`)

// The mock extracts the first non-heading, non-list paragraph with a ~160
// character limit; keep that behavior so cards do not expose raw markdown.
func publicArticleExcerpt(content string) string {
	text := ""
	for _, paragraph := range regexp.MustCompile(`\n\s*\n`).Split(strings.ReplaceAll(content, "\r\n", "\n"), -1) {
		line := strings.TrimSpace(paragraph)
		if line != "" && !strings.HasPrefix(line, "## ") && !strings.HasPrefix(line, "- ") {
			text = articleWhitespace.ReplaceAllString(line, " ")
			break
		}
	}
	runes := []rune(text)
	if len(runes) <= 160 {
		return text
	}
	cut := strings.TrimSpace(string(runes[:157]))
	if idx := strings.LastIndex(cut, " "); idx > 0 {
		cut = cut[:idx]
	}
	return cut + "…"
}

func publicArticleFromRow(r *http.Request, row publicArticleRow) publicArticleItem {
	thumb := publicSQLString(row.ThumbnailURL)
	if thumb != nil {
		v := publicProductMediaURL(r, *thumb)
		thumb = &v
	}
	return publicArticleItem{ID: row.ID, Type: row.Type, Title: row.Title, Slug: row.Slug,
		ThumbnailURL: thumb, Excerpt: publicArticleExcerpt(row.Content), PublishedAt: publicSQLString(row.PublishedAt)}
}

func publicArticleFilter(raw string) ([]string, bool) {
	if strings.TrimSpace(raw) == "" {
		return nil, true
	}
	found := make(map[string]bool)
	types := make([]string, 0, 4)
	for _, part := range strings.Split(raw, ",") {
		value := strings.ToUpper(strings.TrimSpace(part))
		switch value {
		case "ARTIKEL", "PROMO", "EVENT", "HALAMAN":
		default:
			return nil, false
		}
		if !found[value] {
			found[value] = true
			types = append(types, value)
		}
	}
	return types, true
}

func publicArticlePage(q string, fallback int) (int, bool) {
	if strings.TrimSpace(q) == "" {
		return fallback, true
	}
	n, err := strconv.Atoi(q)
	return n, err == nil && n > 0 && n <= 10000
}

func PublicArticlesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	db := postgres.New(svcCtx.Config.Database.DataSource)
	return func(w http.ResponseWriter, r *http.Request) {
		types, valid := publicArticleFilter(r.URL.Query().Get("type"))
		page, goodPage := publicArticlePage(r.URL.Query().Get("page"), 1)
		perPage, goodSize := publicArticlePage(r.URL.Query().Get("per_page"), 9)
		if !valid || !goodPage || !goodSize || perPage > 100 {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusBadRequest, map[string]string{"message": "Parameter artikel tidak valid (per_page maksimal 100)."})
			return
		}
		// PostgreSQL ANY($1) on a text[] requires a pg array driver;
		// instead build only fixed, whitelisted type literals as parameters.
		// The dynamic OR predicates are still parameterized.
		where := "status = 'PUBLISHED' AND (published_at IS NULL OR published_at <= timezone('Asia/Jakarta',now()))"
		args := make([]any, 0, 6)
		if len(types) == 0 {
			where += " AND type <> 'HALAMAN'"
		} else {
			orParts := make([]string, 0, len(types))
			for i, t := range types {
				orParts = append(orParts, "type = $"+strconv.Itoa(i+1))
				args = append(args, t)
			}
			where += " AND (" + strings.Join(orParts, " OR ") + ")"
		}
		var total int64
		if err := db.QueryRowCtx(r.Context(), &total, "SELECT COUNT(*) FROM crm_schema.articles WHERE "+where, args...); err != nil {
			publicContentError(w, r, "artikel", err)
			return
		}
		pages := int((total + int64(perPage) - 1) / int64(perPage))
		if pages < 1 {
			pages = 1
		}
		if page > pages {
			page = pages
		}
		query := `SELECT id,type,title,slug,thumbnail_url,content,status,
                to_char(published_at,'YYYY-MM-DD"T"HH24:MI:SS') AS published_at
            FROM crm_schema.articles WHERE ` + where + `
            ORDER BY published_at DESC NULLS LAST,id DESC
            LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
		args = append(args, perPage, (page-1)*perPage)
		var rows []publicArticleRow
		if err := db.QueryRowsCtx(r.Context(), &rows, query, args...); err != nil {
			publicContentError(w, r, "artikel", err)
			return
		}
		data := make([]publicArticleItem, 0, len(rows))
		for _, row := range rows {
			data = append(data, publicArticleFromRow(r, row))
		}
		httpx.OkJsonCtx(r.Context(), w, publicArticleListResponse{Data: data, Meta: publicArticleMeta{Page: page, PerPage: perPage, Total: total, TotalPages: pages}})
	}
}

func PublicArticleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	db := postgres.New(svcCtx.Config.Database.DataSource)
	return func(w http.ResponseWriter, r *http.Request) {
		var path struct {
			Slug string `path:"slug"`
		}
		if err := httpx.ParsePath(r, &path); err != nil || strings.TrimSpace(path.Slug) == "" {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusBadRequest, map[string]string{"message": "Slug artikel tidak valid."})
			return
		}
		// Intentional: draft previews require a separately designed,
		// cryptographically verified preview token. Do not accept ?preview=
		// as authentication; currently only PUBLISHED can be accessed.
		var row publicArticleRow
		err := db.QueryRowCtx(r.Context(), &row, `SELECT id,type,title,slug,thumbnail_url,content,status,
               to_char(published_at,'YYYY-MM-DD"T"HH24:MI:SS') AS published_at
             FROM crm_schema.articles
             WHERE slug=$1 AND status='PUBLISHED'
               AND (published_at IS NULL OR published_at <= timezone('Asia/Jakarta',now()))`, path.Slug)
		if errors.Is(err, sqlx.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusNotFound, map[string]string{"message": "Artikel tidak ditemukan."})
			return
		}
		if err != nil {
			publicContentError(w, r, "artikel", err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, publicArticleDetailResponse{Data: publicArticleDetail{publicArticleItem: publicArticleFromRow(r, row), Content: row.Content, Status: row.Status}})
	}
}

func publicContentError(w http.ResponseWriter, r *http.Request, part string, err error) {
	logx.WithContext(r.Context()).Errorf("public %s query failed: %v", part, err)
	httpx.WriteJsonCtx(r.Context(), w, http.StatusInternalServerError, map[string]string{"message": "Gagal memuat data " + part + "."})
}