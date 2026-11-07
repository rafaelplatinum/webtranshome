package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"webtranshome/internal/svc"
)

func ExtCheckSKUHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sku := strings.TrimSpace(r.URL.Query().Get("sku"))
		if sku == "" {
			extFail(w, r, 422, "SKU wajib diisi.")
			return
		}
		exclude, _ := strconv.ParseInt(r.URL.Query().Get("exclude"), 10, 64)
		db, _ := extDB(s)
		var exists bool
		e := db.QueryRowContext(r.Context(), `SELECT EXISTS(SELECT 1 FROM crm_schema.products WHERE lower(sku)=lower($1) AND id<>$2)`, sku, exclude).Scan(&exists)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		extJSON(w, r, 200, map[string]any{"data": map[string]any{"available": !exists}})
	}
}
func ExtBrandHomeHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			IDs []int64 `json:"brand_ids"`
		}
		if extBody(r, &in) != nil || len(in.IDs) > 16 {
			extFail(w, r, 422, "Maksimal 16 brand.")
			return
		}
		db, _ := extDB(s)
		seen := map[int64]bool{}
		for _, id := range in.IDs {
			if id < 1 || seen[id] {
				extFail(w, r, 422, "Daftar brand tidak valid.")
				return
			}
			seen[id] = true
		}
		tx, e := db.BeginTx(r.Context(), nil)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		defer tx.Rollback()
		for _, id := range in.IDs {
			var active bool
			e = tx.QueryRowContext(r.Context(), `SELECT is_active FROM crm_schema.brands WHERE id=$1`, id).Scan(&active)
			if errors.Is(e, sql.ErrNoRows) || !active {
				extFail(w, r, 422, "Brand tidak ditemukan atau tidak aktif.")
				return
			}
			if e != nil {
				extServerError(w, r, e)
				return
			}
		}
		data, _ := json.Marshal(in.IDs)
		_, e = tx.ExecContext(r.Context(), `INSERT INTO crm_schema.site_settings(setting_key,setting_value) VALUES('home_featured_brand_ids',$1) ON CONFLICT(setting_key) DO UPDATE SET setting_value=EXCLUDED.setting_value,updated_at=NOW()`, string(data))
		if e != nil {
			extServerError(w, r, e)
			return
		}
		admin, _ := extPrincipal(r)
		if e = extAudit(r.Context(), tx, admin, "brand.home.update", "Brand", nil, map[string]any{"brand_ids": in.IDs}, requestIP(r)); e != nil {
			extServerError(w, r, e)
			return
		}
		if e = tx.Commit(); e != nil {
			extServerError(w, r, e)
			return
		}
		extJSON(w, r, 200, map[string]any{"data": map[string]any{"brand_ids": in.IDs}})
	}
}
func ExtDeleteBrandHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, e := extID(r)
		if e != nil {
			extFail(w, r, 400, "ID salah.")
			return
		}
		db, _ := extDB(s)
		admin, _ := extPrincipal(r)
		tx, e := db.BeginTx(r.Context(), nil)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		defer tx.Rollback()
		var name string
		var active bool
		e = tx.QueryRowContext(r.Context(), `SELECT name,is_active FROM crm_schema.brands WHERE id=$1 FOR UPDATE`, id).Scan(&name, &active)
		if errors.Is(e, sql.ErrNoRows) {
			extFail(w, r, 404, "Brand tidak ditemukan.")
			return
		}
		if e != nil {
			extServerError(w, r, e)
			return
		}
		if active {
			extFail(w, r, 422, "Nonaktifkan brand sebelum menghapus.")
			return
		}
		var used bool
		e = tx.QueryRowContext(r.Context(), `SELECT EXISTS(SELECT 1 FROM crm_schema.products WHERE brand_id=$1)`, id).Scan(&used)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		if used {
			extFail(w, r, 409, "Brand masih dipakai produk.")
			return
		}
		_, e = tx.ExecContext(r.Context(), `DELETE FROM crm_schema.brands WHERE id=$1`, id)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		if e = extAudit(r.Context(), tx, admin, "brand.delete", "Brand", map[string]any{"id": id, "name": name}, nil, requestIP(r)); e != nil {
			extServerError(w, r, e)
			return
		}
		if e = tx.Commit(); e != nil {
			extServerError(w, r, e)
			return
		}
		extJSON(w, r, 200, map[string]any{"success": true})
	}
}