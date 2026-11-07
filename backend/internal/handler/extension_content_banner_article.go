package handler

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"webtranshome/internal/svc"
)

func extBannerMap(rows interface{ Scan(...any) error }) (map[string]any, error) {
	var id int64
	var title, image, position string
	var link, subtitle, video, cta1, cta2, cta2URL sql.NullString
	var start, end sql.NullTime
	var order int
	var active bool
	err := rows.Scan(&id, &title, &image, &link, &position, &start, &end, &order, &active, &subtitle, &video, &cta1, &cta2, &cta2URL)
	if err != nil {
		return nil, err
	}
	status := "TAYANG"
	now := time.Now()
	if !active {
		status = "NONAKTIF"
	} else if start.Valid && start.Time.After(now) {
		status = "TERJADWAL"
	} else if end.Valid && end.Time.Before(now) {
		status = "BERAKHIR"
	}
	return map[string]any{"id": id, "title": title, "image_url": image, "link_url": extReadString(link), "position": position, "start_at": extReadTime(start), "end_at": extReadTime(end), "start_date": extDateValue(start), "end_date": extDateValue(end), "sort_order": order, "is_active": active, "status": status, "subtitle": extReadString(subtitle), "video_url": extReadString(video), "cta_label": extReadString(cta1), "cta2_label": extReadString(cta2), "cta2_url": extReadString(cta2URL)}, nil
}
func extDateValue(t sql.NullTime) any {
	if !t.Valid {
		return nil
	}
	return t.Time.Format("2006-01-02")
}

const extBannerFields = `id,title,image_url,link_url,position,start_at,end_at,sort_order,is_active,subtitle,video_url,cta_label,cta2_label,cta2_url`

func extBannerByID(r *http.Request, s *svc.ServiceContext, id int64) (map[string]any, error) {
	db, _ := extDB(s)
	return extBannerMap(db.QueryRowContext(r.Context(), `SELECT `+extBannerFields+` FROM crm_schema.banners WHERE id=$1`, id))
}
func ExtBannerListHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, _ := extDB(s)
		rows, e := db.QueryContext(r.Context(), `SELECT `+extBannerFields+` FROM crm_schema.banners ORDER BY CASE position WHEN 'HOME_HERO' THEN 0 WHEN 'HOME_SLIDER' THEN 1 WHEN 'HOME_SIDE' THEN 2 ELSE 3 END,sort_order,id`)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		defer rows.Close()
		items := []any{}
		for rows.Next() {
			v, e := extBannerMap(rows)
			if e != nil {
				extServerError(w, r, e)
				return
			}
			items = append(items, v)
		}
		if e = rows.Err(); e != nil {
			extServerError(w, r, e)
			return
		}
		extJSON(w, r, 200, map[string]any{"data": items})
	}
}
func ExtBannerGetHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, e := extID(r)
		if e != nil {
			extFail(w, r, 400, "ID tidak valid.")
			return
		}
		v, e := extBannerByID(r, s, id)
		if errors.Is(e, sql.ErrNoRows) {
			extFail(w, r, 404, "Banner tidak ditemukan.")
			return
		}
		if e != nil {
			extServerError(w, r, e)
			return
		}
		extJSON(w, r, 200, map[string]any{"data": v})
	}
}

type extBannerInput struct {
	Title     string  `json:"title"`
	ImageURL  string  `json:"image_url"`
	LinkURL   *string `json:"link_url"`
	Position  string  `json:"position"`
	StartDate string  `json:"start_date"`
	EndDate   string  `json:"end_date"`
	IsActive  *bool   `json:"is_active"`
	Subtitle  *string `json:"subtitle"`
	VideoURL  *string `json:"video_url"`
	CTALabel  *string `json:"cta_label"`
	CTA2Label *string `json:"cta2_label"`
	CTA2URL   *string `json:"cta2_url"`
}

func extBannerTime(s string, end bool) (any, error) {
	if s == "" {
		return nil, nil
	}
	t, e := extDate(s)
	if e != nil {
		return nil, e
	}
	if end {
		t = t.Add(24*time.Hour - time.Second)
	}
	return t, nil
}
func extBannerValidate(in extBannerInput) error {
	if in.Title == "" || len(in.Title) > 255 || in.ImageURL == "" || len(in.ImageURL) > 255 {
		return errors.New("Judul dan gambar wajib diisi.")
	}
	if in.Position != "HOME_HERO" && in.Position != "HOME_SLIDER" && in.Position != "HOME_SIDE" {
		return errors.New("Posisi banner tidak valid.")
	}
	if in.StartDate == "" {
		return errors.New("Tanggal mulai wajib diisi.")
	}
	a, e := extDate(in.StartDate)
	if e != nil {
		return errors.New("Tanggal mulai tidak valid.")
	}
	if in.EndDate != "" {
		b, e := extDate(in.EndDate)
		if e != nil || b.Before(a) {
			return errors.New("Tanggal selesai tidak valid.")
		}
	}
	for _, link := range []*string{in.LinkURL, in.CTA2URL, in.VideoURL} {
		if link != nil && *link != "" && !strings.HasPrefix(*link, "/") && !strings.HasPrefix(*link, "https://") {
			return errors.New("URL harus HTTPS atau path internal.")
		}
	}
	if in.Subtitle != nil && len(*in.Subtitle) > 160 {
		return errors.New("Subjudul terlalu panjang.")
	}
	if in.CTALabel != nil && len(*in.CTALabel) > 30 {
		return errors.New("Label CTA terlalu panjang.")
	}
	if in.CTA2Label != nil && len(*in.CTA2Label) > 30 {
		return errors.New("Label CTA2 terlalu panjang.")
	}
	return nil
}
func ExtBannerMutationHandler(s *svc.ServiceContext, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, _ := extDB(s)
		admin, _ := extPrincipal(r)
		id := int64(0)
		if method != "create" {
			var e error
			id, e = extID(r)
			if e != nil {
				extFail(w, r, 400, "ID tidak valid.")
				return
			}
		}
		var in extBannerInput
		if method == "patch" {
			var patch struct {
				Active    *bool `json:"is_active"`
				SortOrder *int  `json:"sort_order"`
			}
			if extBody(r, &patch) != nil || (patch.Active == nil && patch.SortOrder == nil) {
				extFail(w, r, 422, "Tidak ada perubahan.")
				return
			}
			tx, e := db.BeginTx(r.Context(), nil)
			if e != nil {
				extServerError(w, r, e)
				return
			}
			defer tx.Rollback()
			var before string
			e = tx.QueryRowContext(r.Context(), `SELECT title FROM crm_schema.banners WHERE id=$1 FOR UPDATE`, id).Scan(&before)
			if errors.Is(e, sql.ErrNoRows) {
				extFail(w, r, 404, "Banner tidak ditemukan.")
				return
			}
			if e != nil {
				extServerError(w, r, e)
				return
			}
			_, e = tx.ExecContext(r.Context(), `UPDATE crm_schema.banners SET is_active=COALESCE($2,is_active),sort_order=COALESCE($3,sort_order) WHERE id=$1`, id, patch.Active, patch.SortOrder)
			if e != nil {
				extServerError(w, r, e)
				return
			}
			if e = extAudit(r.Context(), tx, admin, "banner.update", "Banner", map[string]any{"title": before}, map[string]any{"is_active": patch.Active, "sort_order": patch.SortOrder}, requestIP(r)); e != nil {
				extServerError(w, r, e)
				return
			}
			if e = tx.Commit(); e != nil {
				extServerError(w, r, e)
				return
			}
		} else {
			if extBody(r, &in) != nil {
				extFail(w, r, 422, "JSON tidak valid.")
				return
			}
			in.Title = extTrim(in.Title)
			in.ImageURL = extTrim(in.ImageURL)
			if e := extBannerValidate(in); e != nil {
				extFail(w, r, 422, e.Error())
				return
			}
			start, e := extBannerTime(in.StartDate, false)
			if e != nil {
				extFail(w, r, 422, "Tanggal salah.")
				return
			}
			end, e := extBannerTime(in.EndDate, true)
			if e != nil {
				extFail(w, r, 422, "Tanggal salah.")
				return
			}
			active := true
			if in.IsActive != nil {
				active = *in.IsActive
			}
			tx, e := db.BeginTx(r.Context(), nil)
			if e != nil {
				extServerError(w, r, e)
				return
			}
			defer tx.Rollback()
			if method == "create" {
				e = tx.QueryRowContext(r.Context(), `INSERT INTO crm_schema.banners(title,image_url,link_url,position,start_at,end_at,sort_order,is_active,subtitle,video_url,cta_label,cta2_label,cta2_url) VALUES($1,$2,$3,$4,$5,$6,COALESCE((SELECT MAX(sort_order)+1 FROM crm_schema.banners WHERE position=$4),1),$7,$8,$9,$10,$11,$12) RETURNING id`, in.Title, in.ImageURL, in.LinkURL, in.Position, start, end, active, in.Subtitle, in.VideoURL, in.CTALabel, in.CTA2Label, in.CTA2URL).Scan(&id)
			} else {
				var exists int64
				e = tx.QueryRowContext(r.Context(), `SELECT id FROM crm_schema.banners WHERE id=$1 FOR UPDATE`, id).Scan(&exists)
				if errors.Is(e, sql.ErrNoRows) {
					extFail(w, r, 404, "Banner tidak ditemukan.")
					return
				}
				if e != nil {
					extServerError(w, r, e)
					return
				}
				_, e = tx.ExecContext(r.Context(), `UPDATE crm_schema.banners SET title=$2,image_url=$3,link_url=$4,position=$5,start_at=$6,end_at=$7,is_active=$8,subtitle=$9,video_url=$10,cta_label=$11,cta2_label=$12,cta2_url=$13 WHERE id=$1`, id, in.Title, in.ImageURL, in.LinkURL, in.Position, start, end, active, in.Subtitle, in.VideoURL, in.CTALabel, in.CTA2Label, in.CTA2URL)
			}
			if e != nil {
				extServerError(w, r, e)
				return
			}
			if e = extAudit(r.Context(), tx, admin, "banner."+method, "Banner", nil, map[string]any{"id": id, "title": in.Title}, requestIP(r)); e != nil {
				extServerError(w, r, e)
				return
			}
			if e = tx.Commit(); e != nil {
				extServerError(w, r, e)
				return
			}
		}
		v, e := extBannerByID(r, s, id)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		status := 200
		if method == "create" {
			status = 201
		}
		extJSON(w, r, status, map[string]any{"data": v})
	}
}

func extArticleMap(row interface{ Scan(...any) error }) (map[string]any, error) {
	var id int64
	var typ, title, slug, content, status string
	var thumb, preview sql.NullString
	var published sql.NullTime
	e := row.Scan(&id, &typ, &title, &slug, &thumb, &content, &status, &published, &preview)
	if e != nil {
		return nil, e
	}
	return map[string]any{"id": id, "type": typ, "title": title, "slug": slug, "thumbnail_url": extReadString(thumb), "content": content, "excerpt": publicArticleExcerpt(content), "status": status, "published_at": extReadTime(published), "preview_token": extReadString(preview)}, nil
}

const extArticleFields = `id,type,title,slug,thumbnail_url,content,status,published_at,preview_token`

func ExtArticleListHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, _ := extDB(s)
		page, size := extPage(r, 20)
		args := []any{}
		where := "TRUE"
		q := r.URL.Query()
		if text := extTrim(q.Get("q")); text != "" {
			args = append(args, text)
			where += ` AND (title ILIKE '%'||$1||'%' OR slug ILIKE '%'||$1||'%')`
		}
		if typ := q.Get("type"); typ != "" {
			args = append(args, typ)
			where += " AND type=$" + strconv.Itoa(len(args))
		}
		if status := q.Get("status"); status != "" {
			args = append(args, status)
			where += " AND status=$" + strconv.Itoa(len(args))
		}
		var total int64
		e := db.QueryRowContext(r.Context(), `SELECT COUNT(*) FROM crm_schema.articles WHERE `+where, args...).Scan(&total)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		query := `SELECT ` + extArticleFields + ` FROM crm_schema.articles WHERE ` + where + ` ORDER BY updated_at DESC,id DESC LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
		rows, e := db.QueryContext(r.Context(), query, append(args, size, (page-1)*size)...)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		defer rows.Close()
		items := []any{}
		for rows.Next() {
			v, e := extArticleMap(rows)
			if e != nil {
				extServerError(w, r, e)
				return
			}
			items = append(items, v)
		}
		if e = rows.Err(); e != nil {
			extServerError(w, r, e)
			return
		}
		extJSON(w, r, 200, map[string]any{"data": items, "meta": extMeta(page, size, total)})
	}
}
func ExtArticleGetHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, e := extID(r)
		if e != nil {
			extFail(w, r, 400, "ID salah.")
			return
		}
		db, _ := extDB(s)
		v, e := extArticleMap(db.QueryRowContext(r.Context(), `SELECT `+extArticleFields+` FROM crm_schema.articles WHERE id=$1`, id))
		if errors.Is(e, sql.ErrNoRows) {
			extFail(w, r, 404, "Artikel tidak ditemukan.")
			return
		}
		if e != nil {
			extServerError(w, r, e)
			return
		}
		extJSON(w, r, 200, map[string]any{"data": v})
	}
}

type extArticleInput struct {
	Type         string  `json:"type"`
	Title        string  `json:"title"`
	Slug         string  `json:"slug"`
	ThumbnailURL *string `json:"thumbnail_url"`
	Content      string  `json:"content"`
}

func ExtArticleMutationHandler(s *svc.ServiceContext, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, _ := extDB(s)
		admin, _ := extPrincipal(r)
		id := int64(0)
		if method != "create" {
			var e error
			id, e = extID(r)
			if e != nil {
				extFail(w, r, 400, "ID salah.")
				return
			}
		}
		tx, e := db.BeginTx(r.Context(), nil)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		defer tx.Rollback()
		if method == "status" {
			var in struct {
				Status string `json:"status"`
			}
			if extBody(r, &in) != nil || (in.Status != "DRAFT" && in.Status != "PUBLISHED") {
				extFail(w, r, 422, "Status salah.")
				return
			}
			_, e = tx.ExecContext(r.Context(), `UPDATE crm_schema.articles SET status=$2,published_at=CASE WHEN $2='PUBLISHED' THEN COALESCE(published_at,CURRENT_TIMESTAMP) ELSE NULL END WHERE id=$1`, id, in.Status)
			if e != nil {
				extServerError(w, r, e)
				return
			}
		} else {
			var in extArticleInput
			if extBody(r, &in) != nil {
				extFail(w, r, 422, "Data artikel tidak valid.")
				return
			}
			if len(extTrim(in.Title)) < 2 || len(in.Title) > 255 || !extValidSlug(in.Slug) || extTrim(in.Content) == "" {
				extFail(w, r, 422, "Judul, slug, dan isi wajib valid.")
				return
			}
			switch in.Type {
			case "ARTIKEL", "PROMO", "EVENT", "HALAMAN":
			default:
				extFail(w, r, 422, "Jenis artikel salah.")
				return
			}
			if method == "create" {
				previewBytes := make([]byte, 20)
				if _, er := rand.Read(previewBytes); er != nil {
					extServerError(w, r, er)
					return
				}
				e = tx.QueryRowContext(r.Context(), `INSERT INTO crm_schema.articles(type,title,slug,thumbnail_url,content,status,author_id,preview_token) VALUES ($1,$2,$3,$4,$5,'DRAFT',$6,$7) RETURNING id`, in.Type, extTrim(in.Title), in.Slug, in.ThumbnailURL, in.Content, admin, hex.EncodeToString(previewBytes)).Scan(&id)
			} else {
				previewBytes := make([]byte, 20)
				if _, er := rand.Read(previewBytes); er != nil {
					extServerError(w, r, er)
					return
				}
				result, er := tx.ExecContext(r.Context(), `UPDATE crm_schema.articles SET type=$2,title=$3,slug=$4,thumbnail_url=$5,content=$6,preview_token=COALESCE(preview_token,$7) WHERE id=$1`, id, in.Type, extTrim(in.Title), in.Slug, in.ThumbnailURL, in.Content, hex.EncodeToString(previewBytes))
				e = er
				if e == nil {
					var count int64
					count, _ = result.RowsAffected()
					if count == 0 {
						extFail(w, r, 404, "Artikel tidak ditemukan.")
						return
					}
				}
			}
			if e != nil {
				if strings.Contains(e.Error(), "articles_slug_key") {
					extFail(w, r, 409, "Slug artikel sudah digunakan.")
					return
				}
				extServerError(w, r, e)
				return
			}
		}
		if e = extAudit(r.Context(), tx, admin, "article."+method, "Artikel", nil, map[string]any{"id": id}, requestIP(r)); e != nil {
			extServerError(w, r, e)
			return
		}
		if e = tx.Commit(); e != nil {
			extServerError(w, r, e)
			return
		}
		v, e := extArticleMap(db.QueryRowContext(r.Context(), `SELECT `+extArticleFields+` FROM crm_schema.articles WHERE id=$1`, id))
		if e != nil {
			extServerError(w, r, e)
			return
		}
		status := 200
		if method == "create" {
			status = 201
		}
		extJSON(w, r, status, map[string]any{"data": v})
	}
}