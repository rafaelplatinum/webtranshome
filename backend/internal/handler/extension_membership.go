package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"webtranshome/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
	"golang.org/x/crypto/bcrypt"
)

// IDs exposed to Vue are crm_schema.users.id; point_transactions.member_id
// references member_profiles.id. Transactions always lock member_profiles.
func extUserProfile(ctx context.Context, db *sql.DB, id int64) (map[string]any, error) {
	var email, code, name, tier string
	var phone sql.NullString
	var consent, active, hasPassword, hasGoogle, emailVerified bool
	var since sql.NullTime
	var rolesRaw string
	err := db.QueryRowContext(ctx, `SELECT u.email,COALESCE(mp.member_code,''),COALESCE(mp.full_name,''),COALESCE(mp.tier,''),u.phone_number,
 COALESCE(mp.communication_consent,FALSE),u.is_active,u.password_hash IS NOT NULL,u.google_id IS NOT NULL,
 u.email_verified_at IS NOT NULL,mp.created_at,
 COALESCE((SELECT string_agg(r.code,',') FROM crm_schema.user_roles ur JOIN crm_schema.roles r ON r.id=ur.role_id WHERE ur.user_id=u.id AND r.is_active),'')
 FROM crm_schema.users u LEFT JOIN crm_schema.member_profiles mp ON mp.user_id=u.id
 WHERE u.id=$1`, id).Scan(&email, &code, &name, &tier, &phone, &consent, &active, &hasPassword, &hasGoogle, &emailVerified, &since, &rolesRaw)
	if err != nil {
		return nil, err
	}
	roles := []string{}
	if rolesRaw != "" {
		roles = strings.Split(rolesRaw, ",")
	}
	return map[string]any{"id": id, "email": email, "phone_number": extReadString(phone), "full_name": name,
		"member_code": code, "tier": tier, "communication_consent": consent, "member_since": extReadTime(since),
		"email_verified": emailVerified, "has_password": hasPassword, "google_connected": hasGoogle, "is_active": active, "roles": roles}, nil
}
func extPermissions(ctx context.Context, db *sql.DB, id int64) ([]string, error) {
	rows, e := db.QueryContext(ctx, `SELECT DISTINCT p.code FROM crm_schema.permissions p JOIN crm_schema.role_permissions rp ON rp.permission_id=p.id JOIN crm_schema.roles ro ON ro.id=rp.role_id JOIN crm_schema.user_roles ur ON ur.role_id=ro.id WHERE ur.user_id=$1 AND ro.is_active ORDER BY p.code`, id)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	result := []string{}
	for rows.Next() {
		var p string
		if e = rows.Scan(&p); e != nil {
			return nil, e
		}
		result = append(result, p)
	}
	return result, rows.Err()
}
func ExtMeHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := extPrincipal(r)
		if !ok {
			extFail(w, r, 401, "Sesi tidak valid.")
			return
		}
		db, _ := extDB(s)
		u, e := extUserProfile(r.Context(), db, id)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		perms, e := extPermissions(r.Context(), db, id)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		extJSON(w, r, 200, map[string]any{"data": map[string]any{"user": u, "permissions": perms}})
	}
}
func ExtPatchProfileHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := extPrincipal(r)
		var input struct {
			FullName    string  `json:"full_name"`
			PhoneNumber *string `json:"phone_number"`
		}
		if extBody(r, &input) != nil || len(strings.TrimSpace(input.FullName)) < 2 || len(input.FullName) > 255 {
			extFail(w, r, 422, "Nama minimal 2 karakter.")
			return
		}
		db, _ := extDB(s)
		tx, e := db.BeginTx(r.Context(), nil)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		defer tx.Rollback()
		var phone any
		if input.PhoneNumber != nil && strings.TrimSpace(*input.PhoneNumber) != "" {
			phone = strings.TrimSpace(*input.PhoneNumber)
			if strings.HasPrefix(phone.(string), "08") {
				phone = "62" + phone.(string)[1:]
			}
			if !extPhoneValid(phone.(string)) {
				extFail(w, r, 422, "Nomor WhatsApp tidak valid.")
				return
			}
		}
		if _, e = tx.ExecContext(r.Context(), `UPDATE crm_schema.users SET phone_number=$1,phone_verified_at=CASE WHEN phone_number IS NOT DISTINCT FROM $1 THEN phone_verified_at ELSE NULL END WHERE id=$2`, phone, id); e != nil {
			if strings.Contains(e.Error(), "unique") {
				extFail(w, r, 409, "Nomor sudah digunakan akun lain.")
				return
			}
			extServerError(w, r, e)
			return
		}
		if _, e = tx.ExecContext(r.Context(), `UPDATE crm_schema.member_profiles SET full_name=$1 WHERE user_id=$2`, strings.TrimSpace(input.FullName), id); e != nil {
			extServerError(w, r, e)
			return
		}
		if e = tx.Commit(); e != nil {
			extServerError(w, r, e)
			return
		}
		u, e := extUserProfile(r.Context(), db, id)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		extJSON(w, r, 200, map[string]any{"data": map[string]any{"user": u}})
	}
}
func extPhoneValid(s string) bool {
	if strings.HasPrefix(s, "08") {
		s = "62" + s[1:]
	}
	if !strings.HasPrefix(s, "628") {
		return false
	}
	if len(s) < 10 || len(s) > 15 {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
func ExtConsentHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := extPrincipal(r)
		var in struct {
			Consent bool `json:"communication_consent"`
		}
		if extBody(r, &in) != nil {
			extFail(w, r, 422, "Data tidak valid.")
			return
		}
		db, _ := extDB(s)
		_, e := db.ExecContext(r.Context(), `UPDATE crm_schema.member_profiles SET communication_consent=$1,consent_at=CASE WHEN $1 THEN CURRENT_TIMESTAMP ELSE NULL END WHERE user_id=$2`, in.Consent, id)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		u, e := extUserProfile(r.Context(), db, id)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		extJSON(w, r, 200, map[string]any{"data": map[string]any{"user": u}})
	}
}
func ExtChangePasswordHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := extPrincipal(r)
		var in struct {
			Old string `json:"old_password"`
			New string `json:"new_password"`
		}
		if extBody(r, &in) != nil || len(in.New) < 8 || len(in.New) > 72 {
			extFail(w, r, 422, "Password baru harus 8–72 karakter.")
			return
		}
		db, _ := extDB(s)
		var old sql.NullString
		if e := db.QueryRowContext(r.Context(), `SELECT password_hash FROM crm_schema.users WHERE id=$1`, id).Scan(&old); e != nil {
			extServerError(w, r, e)
			return
		}
		if old.Valid && bcrypt.CompareHashAndPassword([]byte(old.String), []byte(in.Old)) != nil {
			extFail(w, r, 422, "Password lama salah.")
			return
		}
		hash, e := bcrypt.GenerateFromPassword([]byte(in.New), bcrypt.DefaultCost)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		_, e = db.ExecContext(r.Context(), `UPDATE crm_schema.users SET password_hash=$1,auth_version=auth_version+1 WHERE id=$2`, string(hash), id)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		extJSON(w, r, 200, map[string]any{"message": "Password diperbarui. Silakan masuk kembali."})
	}
}
func extMemberQuery(ctx context.Context, db *sql.DB, where string, args []any, limit, offset int) ([]map[string]any, int64, error) {
	var total int64
	if e := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM crm_schema.member_profiles mp JOIN crm_schema.users u ON u.id=mp.user_id LEFT JOIN crm_schema.qontak_contacts qc ON qc.user_id=u.id WHERE `+where, args...).Scan(&total); e != nil {
		return nil, 0, e
	}
	query := `SELECT u.id,mp.member_code,mp.full_name,u.email,u.phone_number,mp.tier,u.is_active,(u.email_verified_at IS NOT NULL),mp.communication_consent,COALESCE(qc.sync_status,''),mp.created_at,COALESCE((SELECT SUM(points) FROM crm_schema.point_transactions pt WHERE pt.member_id=mp.id),0) FROM crm_schema.member_profiles mp JOIN crm_schema.users u ON u.id=mp.user_id LEFT JOIN crm_schema.qontak_contacts qc ON qc.user_id=u.id WHERE ` + where + ` ORDER BY mp.created_at DESC,u.id DESC LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
	rows, e := db.QueryContext(ctx, query, append(args, limit, offset)...)
	if e != nil {
		return nil, 0, e
	}
	defer rows.Close()
	result := []map[string]any{}
	for rows.Next() {
		var id, balance int64
		var code, name, email, tier, sync string
		var phone sql.NullString
		var active, verified, consent bool
		var created time.Time
		if e = rows.Scan(&id, &code, &name, &email, &phone, &tier, &active, &verified, &consent, &sync, &created, &balance); e != nil {
			return nil, 0, e
		}
		var syncVal any
		if sync != "" {
			syncVal = sync
		}
		result = append(result, map[string]any{"id": id, "member_code": code, "full_name": name, "email": email, "phone_number": extReadString(phone), "tier": tier, "is_active": active, "email_verified": verified, "communication_consent": consent, "balance": balance, "sync_status": syncVal, "created_at": created.Format(time.RFC3339)})
	}
	return result, total, rows.Err()
}
func ExtMembersHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, _ := extDB(s)
		q := r.URL.Query()
		page, size := extPage(r, 20)
		where := []string{"TRUE"}
		args := []any{}
		add := func(cond string, v any) { args = append(args, v); where = append(where, fmt.Sprintf(cond, len(args))) }
		if text := strings.TrimSpace(q.Get("q")); text != "" {
			args = append(args, text)
			v := strconv.Itoa(len(args))
			where = append(where, `(mp.member_code ILIKE '%'||$`+v+`||'%' OR mp.full_name ILIKE '%'||$`+v+`||'%' OR u.email ILIKE '%'||$`+v+`||'%' OR u.phone_number ILIKE '%'||$`+v+`||'%')`)
		}
		for _, f := range []struct{ key, column string }{{"consent", "mp.communication_consent"}, {"hp", "u.phone_number IS NOT NULL"}, {"active", "u.is_active"}} {
			if b, ok := extBoolFlag(q.Get(f.key)); ok {
				add("("+f.column+")=$%d", b)
			}
		}
		if st := q.Get("sync"); st == "PENDING" || st == "SYNCED" || st == "FAILED" {
			add("qc.sync_status=$%d", st)
		}
		// Use a single parameter for all ILIKE expressions, no untrusted SQL fragments.
		clause := strings.Join(where, " AND ")
		if strings.Contains(clause, "%%") {
			clause = strings.ReplaceAll(clause, "%%", "%")
		}
		items, total, e := extMemberQuery(r.Context(), db, clause, args, size, (page-1)*size)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		extJSON(w, r, 200, map[string]any{"data": items, "meta": extMeta(page, size, total)})
	}
}
func ExtMemberSearchHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := strings.TrimSpace(r.URL.Query().Get("q"))
		if len(q) < 2 {
			extJSON(w, r, 200, map[string]any{"data": []any{}})
			return
		}
		db, _ := extDB(s)
		where := `(mp.member_code ILIKE '%'||$1||'%' OR mp.full_name ILIKE '%'||$1||'%' OR u.email ILIKE '%'||$1||'%' OR u.phone_number ILIKE '%'||$1||'%')`
		items, _, e := extMemberQuery(r.Context(), db, where, []any{q}, 10, 0)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		extJSON(w, r, 200, map[string]any{"data": items})
	}
}
func ExtMemberDetailHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, e := extID(r)
		if e != nil {
			extFail(w, r, 400, "ID tidak valid.")
			return
		}
		db, _ := extDB(s)
		var profileID int64
		var last sql.NullTime
		var consentAt sql.NullTime
		var google bool
		var status, contactID sql.NullString
		var synced sql.NullTime
		e = db.QueryRowContext(r.Context(), `SELECT mp.id,u.last_login_at,mp.consent_at,u.google_id IS NOT NULL,qc.sync_status,qc.qontak_contact_id,qc.last_synced_at FROM crm_schema.member_profiles mp JOIN crm_schema.users u ON u.id=mp.user_id LEFT JOIN crm_schema.qontak_contacts qc ON qc.user_id=u.id WHERE u.id=$1`, id).Scan(&profileID, &last, &consentAt, &google, &status, &contactID, &synced)
		if errors.Is(e, sql.ErrNoRows) {
			extFail(w, r, 404, "Member tidak ditemukan.")
			return
		}
		if e != nil {
			extServerError(w, r, e)
			return
		}
		list, _, e := extMemberQuery(r.Context(), db, "u.id=$1", []any{id}, 1, 0)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		m := list[0]
		m["last_login_at"] = extReadTime(last)
		m["consent_at"] = extReadTime(consentAt)
		m["google_connected"] = google
		if status.Valid {
			m["qontak"] = map[string]any{"sync_status": status.String, "qontak_contact_id": extReadString(contactID), "last_synced_at": extReadTime(synced)}
		} else {
			m["qontak"] = nil
		}
		extJSON(w, r, 200, map[string]any{"data": m})
	}
}
func ExtMemberPatchHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, e := extID(r)
		if e != nil {
			extFail(w, r, 400, "ID tidak valid.")
			return
		}
		var in struct {
			Active *bool `json:"is_active"`
		}
		if extBody(r, &in) != nil || in.Active == nil {
			extFail(w, r, 422, "is_active wajib diisi.")
			return
		}
		admin, _ := extPrincipal(r)
		db, _ := extDB(s)
		tx, e := db.BeginTx(r.Context(), nil)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		defer tx.Rollback()
		var before bool
		e = tx.QueryRowContext(r.Context(), `SELECT u.is_active FROM crm_schema.users u JOIN crm_schema.member_profiles mp ON mp.user_id=u.id WHERE u.id=$1 FOR UPDATE OF u`, id).Scan(&before)
		if errors.Is(e, sql.ErrNoRows) {
			extFail(w, r, 404, "Member tidak ditemukan.")
			return
		}
		if e != nil {
			extServerError(w, r, e)
			return
		}
		_, e = tx.ExecContext(r.Context(), `UPDATE crm_schema.users SET is_active=$1,auth_version=auth_version+1 WHERE id=$2`, *in.Active, id)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		if e = extAudit(r.Context(), tx, admin, "member.update", "Member", map[string]any{"is_active": before}, map[string]any{"is_active": *in.Active}, requestIP(r)); e != nil {
			extServerError(w, r, e)
			return
		}
		if e = tx.Commit(); e != nil {
			extServerError(w, r, e)
			return
		}
		ExtMemberDetailHandler(s)(w, r)
	}
}
func extPointRows(ctx context.Context, db *sql.DB, where string, args []any, page, size int) ([]map[string]any, int64, error) {
	var total int64
	e := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM crm_schema.point_transactions t JOIN crm_schema.member_profiles mp ON mp.id=t.member_id WHERE `+where, args...).Scan(&total)
	if e != nil {
		return nil, 0, e
	}
	query := `SELECT t.id,mp.user_id,mp.member_code,t.type,t.points,t.reference_no,t.purchase_amount,t.note,COALESCE(au.email,''),t.created_at FROM crm_schema.point_transactions t JOIN crm_schema.member_profiles mp ON mp.id=t.member_id LEFT JOIN crm_schema.users au ON au.id=t.created_by WHERE ` + where + ` ORDER BY t.created_at DESC,t.id DESC LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
	rows, e := db.QueryContext(ctx, query, append(args, size, (page-1)*size)...)
	if e != nil {
		return nil, 0, e
	}
	defer rows.Close()
	result := []map[string]any{}
	for rows.Next() {
		var id, userID int64
		var typ, code, by string
		var pts int
		var ref, note sql.NullString
		var purchase sql.NullFloat64
		var created time.Time
		if e = rows.Scan(&id, &userID, &code, &typ, &pts, &ref, &purchase, &note, &by, &created); e != nil {
			return nil, 0, e
		}
		var purchaseVal any
		if purchase.Valid {
			purchaseVal = purchase.Float64
		}
		var byName any
		if by != "" {
			byName = by
		}
		result = append(result, map[string]any{"id": id, "user_id": userID, "member_code": code, "type": typ, "points": pts, "reference_no": extReadString(ref), "purchase_amount": purchaseVal, "note": extReadString(note), "created_by_name": byName, "created_at": created.Format(time.RFC3339)})
	}
	return result, total, rows.Err()
}
func ExtPointsHandler(s *svc.ServiceContext, own bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := extPrincipal(r)
		page, size := extPage(r, 20)
		q := r.URL.Query()
		args := []any{}
		conditions := []string{"TRUE"}
		userID := int64(0)
		if own {
			userID = id
		} else if q.Get("user_id") != "" {
			v, e := extInt64(q.Get("user_id"))
			if e != nil {
				extFail(w, r, 422, "user_id tidak valid.")
				return
			}
			userID = v
		}
		if userID > 0 {
			args = append(args, userID)
			conditions = append(conditions, fmt.Sprintf("mp.user_id=$%d", len(args)))
		}
		if typ := q.Get("type"); typ == "EARN" || typ == "REDEEM" || typ == "ADJUST" {
			args = append(args, typ)
			conditions = append(conditions, fmt.Sprintf("t.type=$%d", len(args)))
		}
		db, _ := extDB(s)
		items, total, e := extPointRows(r.Context(), db, strings.Join(conditions, " AND "), args, page, size)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		meta := extMeta(page, size, total)
		if userID > 0 {
			var balance int64
			e = db.QueryRowContext(r.Context(), `SELECT COALESCE(SUM(pt.points),0) FROM crm_schema.member_profiles mp LEFT JOIN crm_schema.point_transactions pt ON pt.member_id=mp.id WHERE mp.user_id=$1`, userID).Scan(&balance)
			if e != nil {
				extServerError(w, r, e)
				return
			}
			meta["balance"] = balance
		}
		extJSON(w, r, 200, map[string]any{"data": items, "meta": meta})
	}
}
func ExtCheckReferenceHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ref := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("ref")))
		if ref == "" {
			extFail(w, r, 422, "Nomor nota wajib diisi.")
			return
		}
		db, _ := extDB(s)
		var date time.Time
		var code string
		e := db.QueryRowContext(r.Context(), `SELECT t.created_at,mp.member_code FROM crm_schema.point_transactions t JOIN crm_schema.member_profiles mp ON mp.id=t.member_id WHERE t.type='EARN' AND t.reference_no=$1`, ref).Scan(&date, &code)
		if errors.Is(e, sql.ErrNoRows) {
			extJSON(w, r, 200, map[string]any{"data": map[string]any{"available": true}})
			return
		}
		if e != nil {
			extServerError(w, r, e)
			return
		}
		extJSON(w, r, 200, map[string]any{"data": map[string]any{"available": false, "used_at": date.Format(time.RFC3339), "member_code": code}})
	}
}
func ExtCreatePointHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		admin, _ := extPrincipal(r)
		var in struct {
			UserID    int64   `json:"user_id"`
			Type      string  `json:"type"`
			Points    int64   `json:"points"`
			Reference string  `json:"reference_no"`
			Purchase  float64 `json:"purchase_amount"`
			Note      string  `json:"note"`
		}
		if extBody(r, &in) != nil || in.UserID <= 0 {
			extFail(w, r, 422, "Data poin tidak valid.")
			return
		}
		if in.Type != "EARN" && in.Type != "REDEEM" && in.Type != "ADJUST" {
			extFail(w, r, 422, "Jenis transaksi tidak valid.")
			return
		}
		if in.Type == "ADJUST" {
			db, _ := extDB(s)
			var permitted bool
			e := db.QueryRowContext(r.Context(), `SELECT EXISTS(SELECT 1 FROM crm_schema.user_roles ur JOIN crm_schema.role_permissions rp ON rp.role_id=ur.role_id JOIN crm_schema.permissions p ON p.id=rp.permission_id WHERE ur.user_id=$1 AND p.code='point.adjust')`, admin).Scan(&permitted)
			if e != nil {
				extServerError(w, r, e)
				return
			}
			if !permitted {
				extFail(w, r, 403, "Tidak memiliki izin penyesuaian poin.")
				return
			}
		}
		db, _ := extDB(s)
		tx, e := db.BeginTx(r.Context(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})
		if e != nil {
			extServerError(w, r, e)
			return
		}
		defer tx.Rollback()
		var memberID int64
		var active bool
		e = tx.QueryRowContext(r.Context(), `SELECT mp.id,u.is_active FROM crm_schema.member_profiles mp JOIN crm_schema.users u ON u.id=mp.user_id WHERE mp.user_id=$1 FOR UPDATE OF mp`, in.UserID).Scan(&memberID, &active)
		if errors.Is(e, sql.ErrNoRows) {
			extFail(w, r, 404, "Member tidak ditemukan.")
			return
		}
		if e != nil {
			extServerError(w, r, e)
			return
		}
		if !active {
			extFail(w, r, 422, "Akun member nonaktif.")
			return
		}
		var balance int64
		e = tx.QueryRowContext(r.Context(), `SELECT COALESCE(SUM(points),0) FROM crm_schema.point_transactions WHERE member_id=$1`, memberID).Scan(&balance)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		pts := in.Points
		var reference, purchase any
		if in.Type == "EARN" {
			if in.Purchase <= 0 || in.Purchase > 100000000000 || strings.TrimSpace(in.Reference) == "" {
				extFail(w, r, 422, "Nomor nota dan total belanja wajib diisi.")
				return
			}
			var raw string
			e = tx.QueryRowContext(r.Context(), `SELECT setting_value FROM crm_schema.site_settings WHERE setting_key='point_ratio_rupiah'`).Scan(&raw)
			if errors.Is(e, sql.ErrNoRows) {
				extFail(w, r, 422, "Rasio poin belum dikonfigurasi.")
				return
			}
			if e != nil {
				extServerError(w, r, e)
				return
			}
			ratio, e := strconv.ParseFloat(raw, 64)
			if e != nil || ratio < 1000 {
				extFail(w, r, 422, "Rasio poin tidak valid.")
				return
			}
			purchaseCents := int64(math.Round(in.Purchase * 100))
			ratioCents := int64(math.Round(ratio * 100))
			if ratioCents <= 0 {
				extFail(w, r, 422, "Rasio poin tidak valid.")
				return
			}
			pts = purchaseCents / ratioCents
			if pts <= 0 {
				extFail(w, r, 422, "Belanja belum mencapai nilai minimum poin.")
				return
			}
			reference = strings.ToUpper(strings.TrimSpace(in.Reference))
			purchase = in.Purchase
		} else {
			if pts == 0 || pts > 100000000 || pts < -100000000 || strings.TrimSpace(in.Note) == "" {
				extFail(w, r, 422, "Jumlah poin dan catatan wajib diisi.")
				return
			}
			if in.Type == "REDEEM" {
				if pts < 0 {
					extFail(w, r, 422, "Poin harus positif.")
					return
				}
				pts = -pts
			}
		}
		if balance+pts < 0 {
			extFail(w, r, 422, "Saldo poin tidak mencukupi.")
			return
		}
		var tid int64
		var created time.Time
		e = tx.QueryRowContext(r.Context(), `INSERT INTO crm_schema.point_transactions(member_id,type,points,reference_no,purchase_amount,note,created_by)
 VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING id,created_at`, memberID, in.Type, pts, reference, purchase, strings.TrimSpace(in.Note), admin).Scan(&tid, &created)
		if e != nil {
			if strings.Contains(e.Error(), "uq_points_reference") {
				extFail(w, r, 409, "Nomor nota sudah dipakai.")
				return
			}
			extServerError(w, r, e)
			return
		}
		if e = extAudit(r.Context(), tx, admin, "point.create", "Poin", nil, map[string]any{"member_id": memberID, "type": in.Type, "points": pts, "reference_no": reference}, requestIP(r)); e != nil {
			extServerError(w, r, e)
			return
		}
		if e = tx.Commit(); e != nil {
			extServerError(w, r, e)
			return
		}
		transaction := map[string]any{"id": tid, "user_id": in.UserID, "type": in.Type, "points": pts, "reference_no": reference, "purchase_amount": purchase, "note": in.Note, "created_at": created.Format(time.RFC3339)}
		extJSON(w, r, 201, map[string]any{"data": map[string]any{"transaction": transaction, "balance": balance + pts}})
	}
}

// Guard for future member-only routes: permissions are still checked on admin endpoints.
var _ = httpx.OkJsonCtx