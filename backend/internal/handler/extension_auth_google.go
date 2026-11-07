package handler

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	jwtissuer "webtranshome/internal/modules/identityaccess/infrastructure/token/jwtissuer"
	"webtranshome/internal/svc"

	"golang.org/x/crypto/bcrypt"
)

func extSession(ctx context.Context, s *svc.ServiceContext, userID int64) (map[string]any, error) {
	db, _ := extDB(s)
	u, e := extUserProfile(ctx, db, userID)
	if e != nil {
		return nil, e
	}
	perms, e := extPermissions(ctx, db, userID)
	if e != nil {
		return nil, e
	}
	var version int64
	e = db.QueryRowContext(ctx, `SELECT auth_version FROM crm_schema.users WHERE id=$1 AND is_active`, userID).Scan(&version)
	if e != nil {
		return nil, e
	}
	token, e := jwtissuer.NewIssuer(s.Config.Auth.AccessSecret, s.Config.Auth.AccessExpire).Issue(ctx, userID, version)
	if e != nil {
		return nil, e
	}
	return map[string]any{"token": token, "user": u, "permissions": perms}, nil
}
func ExtFrontendLoginHandler(s *svc.ServiceContext, adminOnly bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if extBody(r, &in) != nil || extTrim(in.Email) == "" || in.Password == "" {
			extFail(w, r, 422, "Email dan password wajib diisi.")
			return
		}
		db, _ := extDB(s)
		var id int64
		var hash sql.NullString
		var active, verified bool
		e := db.QueryRowContext(r.Context(), `SELECT id,password_hash,is_active,email_verified_at IS NOT NULL FROM crm_schema.users WHERE lower(email)=lower($1)`, extTrim(in.Email)).Scan(&id, &hash, &active, &verified)
		if e != nil {
			if !errors.Is(e, sql.ErrNoRows) {
				extServerError(w, r, e)
				return
			}
			extFail(w, r, 401, "Email atau password salah.")
			return
		}
		if !active || !hash.Valid || bcrypt.CompareHashAndPassword([]byte(hash.String), []byte(in.Password)) != nil {
			extFail(w, r, 401, "Email atau password salah.")
			return
		}
		var staff bool
		e = db.QueryRowContext(r.Context(), `SELECT EXISTS(SELECT 1 FROM crm_schema.user_roles ur JOIN crm_schema.roles role ON role.id=ur.role_id WHERE ur.user_id=$1 AND role.is_active AND role.code IN ('SUPER_ADMIN','ADMIN_KATALOG','ADMIN_MEMBERSHIP','ADMIN_KONTEN'))`, id).Scan(&staff)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		if adminOnly && !staff {
			extFail(w, r, 403, "Akun ini bukan Admin.")
			return
		}
		if !staff && !verified {
			extJSON(w, r, 403, map[string]any{"code": "EMAIL_BELUM_VERIFIKASI", "message": "Email belum diverifikasi."})
			return
		}
		sess, e := extSession(r.Context(), s, id)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		_, _ = db.ExecContext(r.Context(), `UPDATE crm_schema.users SET last_login_at=CURRENT_TIMESTAMP WHERE id=$1`, id)
		extJSON(w, r, 200, map[string]any{"data": sess})
	}
}
func ExtLogoutHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := extPrincipal(r)
		db, _ := extDB(s)
		_, e := db.ExecContext(r.Context(), `UPDATE crm_schema.users SET auth_version=auth_version+1 WHERE id=$1`, id)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		extJSON(w, r, 200, map[string]any{"message": "Keluar berhasil."})
	}
}

type extOAuthState struct {
	Nonce    string `json:"n"`
	ReturnTo string `json:"r"`
}

func extRandom(bytes int) (string, error) {
	raw := make([]byte, bytes)
	if _, e := rand.Read(raw); e != nil {
		return "", e
	}
	return hex.EncodeToString(raw), nil
}
func extAllowedReturn(raw string) bool {
	u, e := url.Parse(raw)
	if e != nil || u.Scheme == "" || u.Host == "" || u.User != nil || u.Fragment != "" {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	if u.Path != "/masuk/google" {
		return false
	}
	// Explicit allow-list. Never trust a return_to URL supplied by a browser.
	origins := strings.Split(os.Getenv("GOOGLE_FRONTEND_ORIGINS"), ",")
	for _, origin := range origins {
		if strings.TrimRight(strings.TrimSpace(origin), "/") == u.Scheme+"://"+u.Host {
			return true
		}
	}
	return false
}
func extGoogleReady() bool {
	return os.Getenv("GOOGLE_CLIENT_ID") != "" && os.Getenv("GOOGLE_CLIENT_SECRET") != "" && os.Getenv("GOOGLE_REDIRECT_URL") != "" && os.Getenv("GOOGLE_FRONTEND_ORIGINS") != ""
}
func ExtGoogleRedirectHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !extGoogleReady() {
			extFail(w, r, 503, "Login Google belum dikonfigurasi.")
			return
		}
		ret := r.URL.Query().Get("return_to")
		if !extAllowedReturn(ret) {
			extFail(w, r, 400, "Alamat kembali tidak diizinkan.")
			return
		}
		nonce, e := extRandom(24)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		state, _ := json.Marshal(extOAuthState{Nonce: nonce, ReturnTo: ret})
		http.SetCookie(w, &http.Cookie{Name: "transhome_oauth_state", Value: nonce, Path: "/api/v1/auth/google", MaxAge: 300, HttpOnly: true, Secure: strings.HasPrefix(os.Getenv("GOOGLE_REDIRECT_URL"), "https://"), SameSite: http.SameSiteLaxMode})
		q := url.Values{}
		q.Set("client_id", os.Getenv("GOOGLE_CLIENT_ID"))
		q.Set("redirect_uri", os.Getenv("GOOGLE_REDIRECT_URL"))
		q.Set("response_type", "code")
		q.Set("scope", "openid email profile")
		q.Set("state", base64.RawURLEncoding.EncodeToString(state))
		q.Set("prompt", "select_account")
		w.Header().Set("Cache-Control", "no-store")
		http.Redirect(w, r, "https://accounts.google.com/o/oauth2/v2/auth?"+q.Encode(), http.StatusFound)
	}
}

// This performs code exchange and server-side token validation against Google's
// tokeninfo endpoint. For high-throughput deployment, use a cached JWKS verifier.
type extTokenInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Name          string `json:"name"`
	Audience      string `json:"aud"`
	Issuer        string `json:"iss"`
	Expires       string `json:"exp"`
}

func extGoogleIdentity(ctx context.Context, code string) (extTokenInfo, error) {
	client := &http.Client{Timeout: 8 * time.Second}
	form := url.Values{"code": {code}, "client_id": {os.Getenv("GOOGLE_CLIENT_ID")}, "client_secret": {os.Getenv("GOOGLE_CLIENT_SECRET")}, "redirect_uri": {os.Getenv("GOOGLE_REDIRECT_URL")}, "grant_type": {"authorization_code"}}
	req, _ := http.NewRequestWithContext(ctx, "POST", "https://oauth2.googleapis.com/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, e := client.Do(req)
	if e != nil {
		return extTokenInfo{}, e
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return extTokenInfo{}, errors.New("google code exchange denied")
	}
	var tokens struct {
		IDToken string `json:"id_token"`
	}
	if e = json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&tokens); e != nil || tokens.IDToken == "" {
		return extTokenInfo{}, errors.New("no Google ID token")
	}
	req, e = http.NewRequestWithContext(ctx, "GET", "https://oauth2.googleapis.com/tokeninfo?id_token="+url.QueryEscape(tokens.IDToken), nil)
	if e != nil {
		return extTokenInfo{}, e
	}
	res, e := client.Do(req)
	if e != nil {
		return extTokenInfo{}, e
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return extTokenInfo{}, errors.New("invalid Google ID token")
	}
	var who extTokenInfo
	if e = json.NewDecoder(io.LimitReader(res.Body, 1<<20)).Decode(&who); e != nil {
		return who, e
	}
	if who.Sub == "" || who.Email == "" || who.Audience != os.Getenv("GOOGLE_CLIENT_ID") || who.EmailVerified != "true" || (who.Issuer != "accounts.google.com" && who.Issuer != "https://accounts.google.com") {
		return who, errors.New("Google claims rejected")
	}
	return who, nil
}
func extGoogleUser(ctx context.Context, db *sql.DB, who extTokenInfo) (int64, bool, error) {
	tx, e := db.BeginTx(ctx, nil)
	if e != nil {
		return 0, false, e
	}
	defer tx.Rollback()
	var userID int64
	var active bool
	e = tx.QueryRowContext(ctx, `SELECT id,is_active FROM crm_schema.users WHERE google_id=$1`, who.Sub).Scan(&userID, &active)
	if e == nil {
		if !active {
			return 0, false, errors.New("inactive account")
		}
		_, _ = tx.ExecContext(ctx, `UPDATE crm_schema.users SET last_login_at=NOW() WHERE id=$1`, userID)
		return userID, false, tx.Commit()
	}
	if !errors.Is(e, sql.ErrNoRows) {
		return 0, false, e
	}
	// Do not silently link Google to an existing password account by email.
	var existing int64
	e = tx.QueryRowContext(ctx, `SELECT id FROM crm_schema.users WHERE lower(email)=lower($1)`, who.Email).Scan(&existing)
	if e == nil {
		return 0, false, errors.New("email belongs to an existing account")
	}
	if !errors.Is(e, sql.ErrNoRows) {
		return 0, false, e
	}
	e = tx.QueryRowContext(ctx, `INSERT INTO crm_schema.users(email,google_id,is_active,email_verified_at,last_login_at) VALUES ($1,$2,TRUE,NOW(),NOW()) RETURNING id`, strings.ToLower(who.Email), who.Sub).Scan(&userID)
	if e != nil {
		return 0, false, e
	}
	code, e := extRandom(7)
	if e != nil {
		return 0, false, e
	}
	code = "TH-" + strings.ToUpper(code)
	name := extTrim(who.Name)
	if name == "" {
		name = strings.Split(who.Email, "@")[0]
	}
	if len(name) > 255 {
		name = name[:255]
	}
	_, e = tx.ExecContext(ctx, `INSERT INTO crm_schema.member_profiles(user_id,member_code,full_name) VALUES($1,$2,$3)`, userID, code, name)
	if e != nil {
		return 0, false, e
	}
	_, e = tx.ExecContext(ctx, `INSERT INTO crm_schema.user_roles(user_id,role_id) SELECT $1,id FROM crm_schema.roles WHERE code='CUSTOMER'`, userID)
	if e != nil {
		return 0, false, e
	}
	if e = tx.Commit(); e != nil {
		return 0, false, e
	}
	return userID, true, nil
}
func ExtGoogleCallbackHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Referrer-Policy", "no-referrer")
		cookie, e := r.Cookie("transhome_oauth_state")
		if e != nil {
			extFail(w, r, 400, "Sesi OAuth tidak ada.")
			return
		}
		http.SetCookie(w, &http.Cookie{Name: "transhome_oauth_state", Path: "/api/v1/auth/google", MaxAge: -1, HttpOnly: true, Secure: strings.HasPrefix(os.Getenv("GOOGLE_REDIRECT_URL"), "https://"), SameSite: http.SameSiteLaxMode})
		raw, e := base64.RawURLEncoding.DecodeString(r.URL.Query().Get("state"))
		if e != nil {
			extFail(w, r, 400, "State OAuth tidak valid.")
			return
		}
		var state extOAuthState
		if json.Unmarshal(raw, &state) != nil || !extAllowedReturn(state.ReturnTo) || len(state.Nonce) != 48 || subtle.ConstantTimeCompare([]byte(cookie.Value), []byte(state.Nonce)) != 1 {
			extFail(w, r, 400, "State OAuth tidak cocok.")
			return
		}
		redirect, _ := url.Parse(state.ReturnTo)
		if r.URL.Query().Get("error") != "" {
			q := redirect.Query()
			q.Set("error", "access_denied")
			redirect.RawQuery = q.Encode()
			http.Redirect(w, r, redirect.String(), http.StatusFound)
			return
		}
		code := r.URL.Query().Get("code")
		if code == "" {
			q := redirect.Query()
			q.Set("error", "missing_code")
			redirect.RawQuery = q.Encode()
			http.Redirect(w, r, redirect.String(), http.StatusFound)
			return
		}
		who, e := extGoogleIdentity(r.Context(), code)
		if e != nil {
			q := redirect.Query()
			q.Set("error", "google_failed")
			redirect.RawQuery = q.Encode()
			http.Redirect(w, r, redirect.String(), http.StatusFound)
			return
		}
		db, _ := extDB(s)
		userID, isNew, e := extGoogleUser(r.Context(), db, who)
		if e != nil {
			q := redirect.Query()
			q.Set("error", "account_conflict")
			redirect.RawQuery = q.Encode()
			http.Redirect(w, r, redirect.String(), http.StatusFound)
			return
		}
		sess, e := extSession(r.Context(), s, userID)
		if e != nil {
			extServerError(w, r, e)
			return
		}
		fragment := url.Values{}
		fragment.Set("token", sess["token"].(string))
		if isNew {
			fragment.Set("new", "1")
		}
		redirect.Fragment = fragment.Encode()
		http.Redirect(w, r, redirect.String(), http.StatusFound)
	}
}