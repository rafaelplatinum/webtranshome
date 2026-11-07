package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"webtranshome/internal/middleware"
	"webtranshome/internal/svc"
)

// Shared helpers for the incremental frontend-compatible API. Never log secrets,
// authentication tokens, passwords, or personal details in error responses.
func extDB(s *svc.ServiceContext) (*sql.DB, error) {
	if s.SQLDB == nil {
		return nil, errors.New("SQL connection unavailable")
	}
	return s.SQLDB, nil
}
func extJSON(w http.ResponseWriter, r *http.Request, code int, payload any) {
	httpx.WriteJsonCtx(r.Context(), w, code, payload)
}
func extFail(w http.ResponseWriter, r *http.Request, code int, message string) {
	extJSON(w, r, code, map[string]any{"message": message})
}
func extServerError(w http.ResponseWriter, r *http.Request, err error) {
	logx.WithContext(r.Context()).Error(err)
	extFail(w, r, 500, "Terjadi kendala server.")
}
func extBody(r *http.Request, out any) error {
	d := json.NewDecoder(io.LimitReader(r.Body, (1<<20)+1))
	d.DisallowUnknownFields()
	if err := d.Decode(out); err != nil {
		return err
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		return errors.New("additional JSON not allowed")
	}
	return nil
}
func extID(r *http.Request) (int64, error) {
	var path struct {
		ID int64 `path:"id"`
	}
	if err := httpx.ParsePath(r, &path); err != nil || path.ID <= 0 {
		return 0, errors.New("ID tidak valid")
	}
	return path.ID, nil
}
func extPrincipal(r *http.Request) (int64, bool) { return middleware.UserIDFromContext(r.Context()) }
func extPage(r *http.Request, defaultSize int) (page, size int) {
	page = 1
	size = defaultSize
	if n, _ := strconv.Atoi(r.URL.Query().Get("page")); n > 0 && n <= 100000 {
		page = n
	}
	if n, _ := strconv.Atoi(r.URL.Query().Get("per_page")); n > 0 && n <= 100 {
		size = n
	}
	return
}
func extMeta(page, size int, total int64) map[string]any {
	pages := int(math.Ceil(float64(total) / float64(size)))
	return map[string]any{"page": page, "per_page": size, "total": total, "total_pages": pages}
}
func extReadString(v sql.NullString) any {
	if !v.Valid {
		return nil
	}
	return v.String
}
func extReadTime(v sql.NullTime) any {
	if !v.Valid {
		return nil
	}
	return v.Time.Format(time.RFC3339)
}
func extTrim(s string) string { return strings.TrimSpace(s) }

var extSlugRE = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func extValidSlug(s string) bool { return len(s) <= 255 && extSlugRE.MatchString(s) }
func extAudit(ctx context.Context, q interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}, admin int64, action, module string, before, after any, ip string) error {
	var b, a []byte
	if before != nil {
		b, _ = json.Marshal(before)
	}
	if after != nil {
		a, _ = json.Marshal(after)
	}
	if len(ip) > 45 {
		ip = ip[:45]
	}
	_, err := q.ExecContext(ctx, `INSERT INTO crm_schema.activity_logs(admin_id,action,module,before_data,after_data,ip_address)
 VALUES ($1,$2,$3,$4::jsonb,$5::jsonb,$6)`, admin, action, module, sqlNullJSON(b), sqlNullJSON(a), ip)
	return err
}
func sqlNullJSON(raw []byte) any {
	if raw == nil {
		return nil
	}
	return string(raw)
}
func extBoolFlag(raw string) (bool, bool) {
	if raw == "1" {
		return true, true
	}
	if raw == "0" {
		return false, true
	}
	return false, false
}
func extDate(raw string) (time.Time, error) { return time.Parse("2006-01-02", raw) }
func extInt64(raw string) (int64, error) {
	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || v < 1 {
		return 0, fmt.Errorf("invalid id")
	}
	return v, nil
}