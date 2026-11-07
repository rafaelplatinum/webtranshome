package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"webtranshome/internal/svc"
)

const extServiceFields = `id,name,slug,summary,description,image_url,service_area,sort_order,is_active`

func extServiceMap(row interface{ Scan(...any) error }) (map[string]any, error) {
	var id int64
	var name, slug, summary string
	var desc, img, area sql.NullString
	var sort int
	var active bool
	e := row.Scan(&id, &name, &slug, &summary, &desc, &img, &area, &sort, &active)
	if e != nil {
		return nil, e
	}
	return map[string]any{"id": id, "name": name, "slug": slug, "summary": summary, "description": extReadString(desc), "image_url": extReadString(img), "service_area": extReadString(area), "sort_order": sort, "is_active": active, "category_ids": []int64{}, "categories": []any{}}, nil
}
func extServiceWithCategories(r *http.Request, s *svc.ServiceContext, id int64) (map[string]any, error) {
	db, _ := extDB(s)
	v, e := extServiceMap(db.QueryRowContext(r.Context(), `SELECT `+extServiceFields+` FROM crm_schema.services WHERE id=$1`, id))
	if e != nil {
		return nil, e
	}
	rows, e := db.QueryContext(r.Context(), `SELECT c.id,c.name,c.slug,c.parent_id FROM crm_schema.service_categories sc JOIN crm_schema.categories c ON c.id=sc.category_id WHERE sc.service_id=$1 ORDER BY c.name`, id)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	ids := []int64{}
	cats := []map[string]any{}
	for rows.Next() {
		var id int64
		var name, slug string
		var parent sql.NullInt64
		if e = rows.Scan(&id, &name, &slug, &parent); e != nil {
			return nil, e
		}
		ids = append(ids, id)
		var pid any
		if parent.Valid {
			pid = parent.Int64
		}
		cats = append(cats, map[string]any{"id": id, "name": name, "slug": slug, "parent_id": pid})
	}
	if e = rows.Err(); e != nil {
		return nil, e
	}
	v["category_ids"] = ids
	v["categories"] = cats
	return v, nil
}
func ExtServiceListHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, _ := extDB(s)
		rows, e := db.QueryContext(r.Context(), `SELECT id FROM crm_schema.services ORDER BY sort_order,id`)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		defer rows.Close()
		ids := []int64{}
		for rows.Next() {
			var id int64
			if e = rows.Scan(&id); e != nil {
				extServerError(w, r, e)
				return
			}
			ids = append(ids, id)
		}
		if e = rows.Err(); e != nil {
			extServerError(w, r, e)
			return
		}
		items := []any{}
		for _, id := range ids {
			v, e := extServiceWithCategories(r, s, id)
			if e != nil {
				extServerError(w, r, e)
				return
			}
			items = append(items, v)
		}
		extJSON(w, r, 200, map[string]any{"data": items})
	}
}
func ExtServiceGetHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, e := extID(r)
		if e != nil {
			extFail(w, r, 400, "ID salah.")
			return
		}
		v, e := extServiceWithCategories(r, s, id)
		if errors.Is(e, sql.ErrNoRows) {
			extFail(w, r, 404, "Jasa tidak ditemukan.")
			return
		}
		if e != nil {
			extServerError(w, r, e)
			return
		}
		extJSON(w, r, 200, map[string]any{"data": v})
	}
}

type extServiceInput struct {
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Summary     string  `json:"summary"`
	Description *string `json:"description"`
	ImageURL    *string `json:"image_url"`
	ServiceArea *string `json:"service_area"`
	CategoryIDs []int64 `json:"category_ids"`
	IsActive    *bool   `json:"is_active"`
}

func ExtServiceMutationHandler(s *svc.ServiceContext, method string) http.HandlerFunc {
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
		if method == "patch" {
			var in struct {
				Active *bool `json:"is_active"`
				Order  *int  `json:"sort_order"`
			}
			if extBody(r, &in) != nil || (in.Active == nil && in.Order == nil) {
				extFail(w, r, 422, "Perubahan tidak valid.")
				return
			}
			tx, e := db.BeginTx(r.Context(), nil)
			if e != nil {
				extServerError(w, r, e)
				return
			}
			defer tx.Rollback()
			res, e := tx.ExecContext(r.Context(), `UPDATE crm_schema.services SET is_active=COALESCE($2,is_active),sort_order=COALESCE($3,sort_order) WHERE id=$1`, id, in.Active, in.Order)
			if e != nil {
				extServerError(w, r, e)
				return
			}
			n, _ := res.RowsAffected()
			if n == 0 {
				extFail(w, r, 404, "Jasa tidak ditemukan.")
				return
			}
			if e = extAudit(r.Context(), tx, admin, "service.update", "Jasa", nil, map[string]any{"id": id, "is_active": in.Active, "sort_order": in.Order}, requestIP(r)); e != nil {
				extServerError(w, r, e)
				return
			}
			if e = tx.Commit(); e != nil {
				extServerError(w, r, e)
				return
			}
		} else {
			var in extServiceInput
			if extBody(r, &in) != nil || len(extTrim(in.Name)) < 2 || len(in.Name) > 100 || !extValidSlug(in.Slug) || len(extTrim(in.Summary)) < 2 || len(in.Summary) > 160 {
				extFail(w, r, 422, "Nama, slug, ringkasan jasa tidak valid.")
				return
			}
			if in.ServiceArea != nil && len(*in.ServiceArea) > 100 {
				extFail(w, r, 422, "Area layanan terlalu panjang.")
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
				e = tx.QueryRowContext(r.Context(), `INSERT INTO crm_schema.services(name,slug,summary,description,image_url,service_area,is_active,sort_order) VALUES ($1,$2,$3,$4,$5,$6,$7,COALESCE((SELECT max(sort_order)+1 FROM crm_schema.services),1)) RETURNING id`, in.Name, in.Slug, in.Summary, in.Description, in.ImageURL, in.ServiceArea, active).Scan(&id)
			} else {
				var res sql.Result
				res, e = tx.ExecContext(r.Context(), `UPDATE crm_schema.services SET name=$2,slug=$3,summary=$4,description=$5,image_url=$6,service_area=$7,is_active=$8 WHERE id=$1`, id, in.Name, in.Slug, in.Summary, in.Description, in.ImageURL, in.ServiceArea, active)
				if e == nil {
					n, _ := res.RowsAffected()
					if n == 0 {
						extFail(w, r, 404, "Jasa tidak ditemukan.")
						return
					}
				}
			}
			if e != nil {
				if strings.Contains(e.Error(), "services_slug_key") {
					extFail(w, r, 409, "Slug jasa sudah dipakai.")
					return
				}
				extServerError(w, r, e)
				return
			}
			if _, e = tx.ExecContext(r.Context(), `DELETE FROM crm_schema.service_categories WHERE service_id=$1`, id); e != nil {
				extServerError(w, r, e)
				return
			}
			for _, cat := range in.CategoryIDs {
				if cat <= 0 {
					extFail(w, r, 422, "Kategori tidak valid.")
					return
				}
				_, e = tx.ExecContext(r.Context(), `INSERT INTO crm_schema.service_categories(service_id,category_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, id, cat)
				if e != nil {
					extFail(w, r, 422, "Kategori tidak ditemukan.")
					return
				}
			}
			if e = extAudit(r.Context(), tx, admin, "service."+method, "Jasa", nil, map[string]any{"id": id, "name": in.Name}, requestIP(r)); e != nil {
				extServerError(w, r, e)
				return
			}
			if e = tx.Commit(); e != nil {
				extServerError(w, r, e)
				return
			}
		}
		v, e := extServiceWithCategories(r, s, id)
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

var extSettingsKeys = []string{"wa_number", "address", "opening_hours", "maps_url", "point_ratio_rupiah"}
var extTemplateKeys = []string{"wa_template_product", "wa_template_cart", "wa_template_cart_item"}

func extSettingMap(r *http.Request, db *sql.DB, keys []string) (map[string]any, error) {
	data := map[string]any{}
	for _, key := range keys {
		data[key] = nil
	}
	rows, e := db.QueryContext(r.Context(), `SELECT setting_key,setting_value FROM crm_schema.site_settings WHERE setting_key=ANY($1)`, keys)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	for rows.Next() {
		var key string
		var val sql.NullString
		if e = rows.Scan(&key, &val); e != nil {
			return nil, e
		}
		data[key] = extReadString(val)
	}
	return data, rows.Err()
}
func ExtSettingsHandler(s *svc.ServiceContext, templates bool, write bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, _ := extDB(s)
		keys := extSettingsKeys
		if templates {
			keys = extTemplateKeys
		}
		if !write {
			v, e := extSettingMap(r, db, keys)
			if e != nil {
				extServerError(w, r, e)
				return
			}
			extJSON(w, r, 200, map[string]any{"data": v})
			return
		}
		var in map[string]any
		if extBody(r, &in) != nil {
			extFail(w, r, 422, "JSON tidak valid.")
			return
		}
		for key := range in {
			ok := false
			for _, allowed := range keys {
				if key == allowed {
					ok = true
					break
				}
			}
			if !ok {
				extFail(w, r, 422, "Nama pengaturan tidak dikenal.")
				return
			}
		}
		if templates {
			for _, key := range keys {
				v, ok := in[key].(string)
				if !ok || len(v) > 2000 || len(v) < 5 {
					extFail(w, r, 422, "Template WhatsApp wajib diisi dengan benar.")
					return
				}
			}
		}
		if v, ok := in["point_ratio_rupiah"]; ok {
			admin, _ := extPrincipal(r)
			var allowed bool
			pe := db.QueryRowContext(r.Context(), `SELECT EXISTS (SELECT 1 FROM crm_schema.user_roles ur JOIN crm_schema.role_permissions rp ON rp.role_id=ur.role_id JOIN crm_schema.permissions p ON p.id=rp.permission_id WHERE ur.user_id=$1 AND p.code='setting.point_ratio')`, admin).Scan(&allowed)
			if pe != nil {
				extServerError(w, r, pe)
				return
			}
			if !allowed {
				extFail(w, r, 403, "Izin setting.point_ratio diperlukan.")
				return
			}
			_ = v
			n, e := strconv.ParseFloat(strings.TrimSpace(toString(v)), 64)
			if e != nil || n < 1000 || n > 1e12 {
				extFail(w, r, 422, "Rasio poin minimal 1000.")
				return
			}
		}
		if v, ok := in["maps_url"].(string); ok && v != "" && !strings.HasPrefix(v, "https://") {
			extFail(w, r, 422, "Maps URL harus HTTPS.")
			return
		}
		tx, e := db.BeginTx(r.Context(), nil)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		defer tx.Rollback()
		for _, key := range keys {
			value, exists := in[key]
			if !exists {
				continue
			}
			v := toString(value)
			if len(v) > 2000 {
				extFail(w, r, 422, "Isi pengaturan terlalu panjang.")
				return
			}
			_, e = tx.ExecContext(r.Context(), `INSERT INTO crm_schema.site_settings(setting_key,setting_value) VALUES($1,$2) ON CONFLICT(setting_key) DO UPDATE SET setting_value=excluded.setting_value,updated_at=NOW()`, key, v)
			if e != nil {
				extServerError(w, r, e)
				return
			}
		}
		admin, _ := extPrincipal(r)
		action := "setting.update"
		if templates {
			action = "whatsapp_template.update"
		}
		// Log the keys only, not the message contents or customer data.
		changed := []string{}
		for key := range in {
			changed = append(changed, key)
		}
		if e = extAudit(r.Context(), tx, admin, action, "Pengaturan", nil, map[string]any{"fields": changed}, requestIP(r)); e != nil {
			extServerError(w, r, e)
			return
		}
		if e = tx.Commit(); e != nil {
			extServerError(w, r, e)
			return
		}
		v, e := extSettingMap(r, db, keys)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		extJSON(w, r, 200, map[string]any{"data": v})
	}
}
func toString(a any) string {
	switch v := a.(type) {
	case string:
		return strings.TrimSpace(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case nil:
		return ""
	default:
		return ""
	}
}