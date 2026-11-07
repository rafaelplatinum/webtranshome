package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"webtranshome/internal/svc"
)

func ExtAuditListHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, _ := extDB(s)
		page, size := extPage(r, 20)
		args := []any{}
		where := []string{"TRUE"}
		q := r.URL.Query()
		if id := q.Get("admin_id"); id != "" {
			v, e := extInt64(id)
			if e != nil {
				extFail(w, r, 422, "admin_id tidak valid.")
				return
			}
			args = append(args, v)
			where = append(where, "l.admin_id=$"+strconv.Itoa(len(args)))
		}
		if module := q.Get("module"); module != "" {
			args = append(args, module)
			where = append(where, "l.module=$"+strconv.Itoa(len(args)))
		}
		for _, f := range []struct{ name, operator string }{{"from", ">="}, {"to", "<="}} {
			if raw := q.Get(f.name); raw != "" {
				d, e := extDate(raw)
				if e != nil {
					extFail(w, r, 422, "Format tanggal harus YYYY-MM-DD.")
					return
				}
				if f.name == "to" {
					d = d.Add(24 * time.Hour)
				}
				args = append(args, d)
				op := f.operator
				if f.name == "to" {
					op = "<"
				}
				where = append(where, "l.created_at"+op+"$"+strconv.Itoa(len(args)))
			}
		}
		cond := strings.Join(where, " AND ")
		var total int64
		e := db.QueryRowContext(r.Context(), `SELECT COUNT(*) FROM crm_schema.activity_logs l WHERE `+cond, args...).Scan(&total)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		query := `SELECT l.id,l.admin_id,COALESCE(mp.full_name,u.email,'Admin'),l.action,l.module,l.before_data,l.after_data,l.ip_address,l.created_at FROM crm_schema.activity_logs l LEFT JOIN crm_schema.users u ON u.id=l.admin_id LEFT JOIN crm_schema.member_profiles mp ON mp.user_id=u.id WHERE ` + cond + ` ORDER BY l.created_at DESC,l.id DESC LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
		rows, e := db.QueryContext(r.Context(), query, append(args, size, (page-1)*size)...)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		defer rows.Close()
		items := []any{}
		for rows.Next() {
			var id, adminID int64
			var name, action, module string
			var before, after []byte
			var ip sql.NullString
			var when time.Time
			if e = rows.Scan(&id, &adminID, &name, &action, &module, &before, &after, &ip, &when); e != nil {
				extServerError(w, r, e)
				return
			}
			items = append(items, map[string]any{"id": id, "admin_id": adminID, "admin_name": name, "action": action, "module": module, "before_data": extAuditJSON(before), "after_data": extAuditJSON(after), "ip_address": extReadString(ip), "created_at": when.Format(time.RFC3339)})
		}
		if e = rows.Err(); e != nil {
			extServerError(w, r, e)
			return
		}
		extJSON(w, r, 200, map[string]any{"data": items, "meta": extMeta(page, size, total)})
	}
}
func extAuditJSON(raw []byte) any {
	if len(raw) == 0 {
		return nil
	}
	var v any
	if json.Unmarshal(raw, &v) != nil {
		return nil
	}
	return v
}
func ExtAuditOptionsHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, _ := extDB(s)
		admins := []map[string]any{}
		modules := []string{}
		rows, e := db.QueryContext(r.Context(), `SELECT DISTINCT l.admin_id,COALESCE(mp.full_name,u.email,'Admin') FROM crm_schema.activity_logs l JOIN crm_schema.users u ON u.id=l.admin_id LEFT JOIN crm_schema.member_profiles mp ON mp.user_id=u.id ORDER BY 2`)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		for rows.Next() {
			var id int64
			var name string
			if e = rows.Scan(&id, &name); e != nil {
				rows.Close()
				extServerError(w, r, e)
				return
			}
			admins = append(admins, map[string]any{"id": id, "name": name})
		}
		e = rows.Err()
		rows.Close()
		if e != nil {
			extServerError(w, r, e)
			return
		}
		rows, e = db.QueryContext(r.Context(), `SELECT DISTINCT module FROM crm_schema.activity_logs ORDER BY module`)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var module string
			if e = rows.Scan(&module); e != nil {
				extServerError(w, r, e)
				return
			}
			modules = append(modules, module)
		}
		if e = rows.Err(); e != nil {
			extServerError(w, r, e)
			return
		}
		extJSON(w, r, 200, map[string]any{"data": map[string]any{"admins": admins, "modules": modules}})
	}
}

// Legacy catalog handlers are not transactional with the activity_logs table.
// Capture successful writes for the audit screen; failures in audit recording
// must be surfaced to server logs, not as a false failure of a completed write.
func extLegacyCatalogAudit(s *svc.ServiceContext, action, module string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recorder := httptest.NewRecorder()
		next(recorder, r)
		status := recorder.Code
		if status >= 200 && status < 300 {
			if id, ok := extPrincipal(r); ok {
				db, _ := extDB(s)
				if err := extAudit(r.Context(), db, id, action, module, nil, map[string]any{"path": r.URL.EscapedPath()}, requestIP(r)); err != nil {
					logx.WithContext(r.Context()).Error(err)
				}
			}
		}
		for k, values := range recorder.Header() {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(status)
		_, _ = w.Write(recorder.Body.Bytes())
	}
}