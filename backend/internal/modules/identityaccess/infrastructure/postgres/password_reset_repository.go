package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"webtranshome/internal/modules/identityaccess/application/passwordreset"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type PasswordResetRepository struct{ connection sqlx.SqlConn }

func NewPasswordResetRepository(connection sqlx.SqlConn) *PasswordResetRepository {
	return &PasswordResetRepository{connection: connection}
}

func (r *PasswordResetRepository) FindAccount(ctx context.Context, channel, destination string) (passwordreset.Account, error) {
	column := "email"
	if channel == passwordreset.ChannelWhatsApp {
		column = "phone_number"
	}
	var row struct {
		ID            int64          `db:"id"`
		Email         string         `db:"email"`
		Phone         sql.NullString `db:"phone_number"`
		PhoneVerified sql.NullTime   `db:"phone_verified_at"`
		EmailVerified sql.NullTime   `db:"email_verified_at"`
		Active        bool           `db:"is_active"`
	}
	err := r.connection.QueryRowCtx(ctx, &row,
		`SELECT id, email, phone_number, phone_verified_at, email_verified_at, is_active FROM crm_schema.users WHERE `+column+` = $1 LIMIT 1`, destination)
	if errors.Is(err, sqlx.ErrNotFound) {
		return passwordreset.Account{}, sql.ErrNoRows
	}
	if err != nil {
		return passwordreset.Account{}, err
	}
	account := passwordreset.Account{ID: row.ID, Email: row.Email, Active: row.Active, PhoneVerified: row.PhoneVerified.Valid, EmailVerified: row.EmailVerified.Valid, Channel: channel}
	if row.Phone.Valid {
		account.Phone = row.Phone.String
	}
	return account, nil
}

func (r *PasswordResetRepository) PhoneForUser(ctx context.Context, userID int64) (string, error) {
	var phone sql.NullString
	err := r.connection.QueryRowCtx(ctx, &phone,
		`SELECT phone_number FROM crm_schema.users WHERE id = $1 AND is_active = TRUE`, userID)
	if errors.Is(err, sqlx.ErrNotFound) {
		return "", sql.ErrNoRows
	}
	if err != nil {
		return "", err
	}
	return phone.String, nil
}

func (r *PasswordResetRepository) RequestCounts(ctx context.Context, destinationHash, ip string, now time.Time) (passwordreset.RequestCounts, error) {
	var row struct {
		DestinationHour int          `db:"destination_hour"`
		DestinationDay  int          `db:"destination_day"`
		IPHour          int          `db:"ip_hour"`
		IPDay           int          `db:"ip_day"`
		LastSentAt      sql.NullTime `db:"last_sent_at"`
	}
	err := r.connection.QueryRowCtx(ctx, &row, `
		SELECT
			COUNT(*) FILTER (WHERE destination_hash = $1 AND created_at >= $3 - INTERVAL '1 hour')::int AS destination_hour,
			COUNT(*) FILTER (WHERE destination_hash = $1 AND created_at >= $3 - INTERVAL '24 hours')::int AS destination_day,
			COUNT(*) FILTER (WHERE request_ip = $2 AND created_at >= $3 - INTERVAL '1 hour')::int AS ip_hour,
			COUNT(*) FILTER (WHERE request_ip = $2 AND created_at >= $3 - INTERVAL '24 hours')::int AS ip_day,
			MAX(created_at) FILTER (WHERE destination_hash = $1) AS last_sent_at
		FROM crm_schema.otp_challenge
		WHERE created_at >= $3 - INTERVAL '24 hours' AND (destination_hash = $1 OR request_ip = $2)`, destinationHash, ip, now)
	if err != nil {
		return passwordreset.RequestCounts{}, err
	}
	counts := passwordreset.RequestCounts{DestinationHour: row.DestinationHour, DestinationDay: row.DestinationDay, IPHour: row.IPHour, IPDay: row.IPDay}
	if row.LastSentAt.Valid {
		counts.LastSentAt = &row.LastSentAt.Time
	}
	return counts, nil
}

func (r *PasswordResetRepository) VerifyFailureCounts(ctx context.Context, destinationHash, ip string, since time.Time) (int, int, error) {
	var row struct {
		Destination int `db:"destination_count"`
		IP          int `db:"ip_count"`
	}
	err := r.connection.QueryRowCtx(ctx, &row, `
		SELECT
			COUNT(*) FILTER (WHERE destination_hash = $1)::int AS destination_count,
			COUNT(*) FILTER (WHERE request_ip = $2)::int AS ip_count
		FROM crm_schema.auth_security_events
		WHERE event_type = 'otp.verify.failed' AND created_at >= $3 AND (destination_hash = $1 OR request_ip = $2)`, destinationHash, ip, since)
	return row.Destination, row.IP, err
}

func (r *PasswordResetRepository) CreateChallenge(ctx context.Context, challenge passwordreset.Challenge, ip string) (int64, error) {
	var challengeID int64
	err := r.connection.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		lockKey := challenge.DestinationHash + ":" + challenge.Channel + ":" + challenge.Purpose
		var lockResult int64
		if err := session.QueryRowCtx(ctx, &lockResult, `SELECT 1 FROM (SELECT pg_advisory_xact_lock(hashtextextended($1, 0))) AS advisory_lock`, lockKey); err != nil {
			return err
		}
		if _, err := session.ExecCtx(ctx, `UPDATE crm_schema.otp_challenge SET consumed_at = CURRENT_TIMESTAMP
			WHERE destination_hash = $1 AND channel = $2 AND purpose = $3 AND consumed_at IS NULL`,
			challenge.DestinationHash, challenge.Channel, challenge.Purpose); err != nil {
			return err
		}
		var userID any
		if challenge.UserID > 0 {
			userID = challenge.UserID
		}
		return session.QueryRowCtx(ctx, &challengeID, `INSERT INTO crm_schema.otp_challenge
			(user_id, purpose, channel, destination_hash, code_hash, expires_at, max_attempts, request_ip)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`, userID, challenge.Purpose, challenge.Channel,
			challenge.DestinationHash, challenge.CodeHash, challenge.ExpiresAt, challenge.MaxAttempts, ip)
	})
	return challengeID, err
}

func (r *PasswordResetRepository) LatestChallenge(ctx context.Context, destinationHash, channel, purpose string, userID int64) (passwordreset.Challenge, error) {
	query := `SELECT id, COALESCE(user_id, 0) AS user_id, channel, purpose, destination_hash, code_hash,
		expires_at, attempts, max_attempts, consumed_at, created_at
		FROM crm_schema.otp_challenge WHERE destination_hash = $1 AND channel = $2 AND purpose = $3 `
	args := []any{destinationHash, channel, purpose}
	if userID > 0 {
		query += `AND user_id = $4 `
		args = append(args, userID)
	}
	query += `ORDER BY created_at DESC, id DESC LIMIT 1`
	var challenge passwordreset.Challenge
	err := r.connection.QueryRowCtx(ctx, &challenge, query, args...)
	if errors.Is(err, sqlx.ErrNotFound) {
		return passwordreset.Challenge{}, sql.ErrNoRows
	}
	return challenge, err
}

func (r *PasswordResetRepository) RecordEvent(ctx context.Context, userID int64, eventType, channel, destinationHash, ip string) error {
	var id any
	if userID > 0 {
		id = userID
	}
	var eventChannel any
	if channel != "" {
		eventChannel = channel
	}
	var hash any
	if destinationHash != "" {
		hash = destinationHash
	}
	_, err := r.connection.ExecCtx(ctx, `INSERT INTO crm_schema.auth_security_events
		(user_id, event_type, channel, destination_hash, request_ip) VALUES ($1,$2,$3,$4,$5)`, id, eventType, eventChannel, hash, ip)
	return err
}

func (r *PasswordResetRepository) IncrementFailedAttempt(ctx context.Context, challengeID int64) error {
	_, err := r.connection.ExecCtx(ctx, `UPDATE crm_schema.otp_challenge
		SET attempts = attempts + 1,
			consumed_at = CASE WHEN attempts + 1 >= max_attempts THEN CURRENT_TIMESTAMP ELSE consumed_at END
		WHERE id = $1 AND consumed_at IS NULL`, challengeID)
	return err
}

func (r *PasswordResetRepository) IssueResetToken(ctx context.Context, challengeID int64, tokenHash, channel string, expiresAt time.Time) error {
	return r.connection.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var userID int64
		err := session.QueryRowCtx(ctx, &userID, `UPDATE crm_schema.otp_challenge
			SET consumed_at = CURRENT_TIMESTAMP
			WHERE id = $1 AND consumed_at IS NULL AND expires_at > CURRENT_TIMESTAMP AND attempts < max_attempts
				AND user_id IS NOT NULL AND channel = $2 AND purpose = 'password_reset'
			RETURNING user_id`, challengeID, channel)
		if errors.Is(err, sqlx.ErrNotFound) {
			return sql.ErrNoRows
		}
		if err != nil {
			return err
		}
		_, err = session.ExecCtx(ctx, `INSERT INTO crm_schema.password_reset_token (user_id, channel, token_hash, expires_at)
			VALUES ($1,$2,$3,$4)`, userID, channel, tokenHash, expiresAt)
		return err
	})
}

func (r *PasswordResetRepository) ResetPassword(ctx context.Context, tokenHash, passwordHash, ip, _ string) (passwordreset.Account, error) {
	var account passwordreset.Account
	err := r.connection.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var row struct {
			UserID  int64  `db:"user_id"`
			Channel string `db:"channel"`
		}
		if err := session.QueryRowCtx(ctx, &row, `UPDATE crm_schema.password_reset_token
			SET used_at = CURRENT_TIMESTAMP
			WHERE token_hash = $1 AND used_at IS NULL AND expires_at > CURRENT_TIMESTAMP
			RETURNING user_id, channel`, tokenHash); err != nil {
			if errors.Is(err, sqlx.ErrNotFound) {
				return sql.ErrNoRows
			}
			return err
		}
		result, err := session.ExecCtx(ctx, `UPDATE crm_schema.users SET password_hash = $2,
			auth_version = auth_version + 1, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND is_active = TRUE`, row.UserID, passwordHash)
		if err != nil {
			return err
		}
		if changed, err := result.RowsAffected(); err != nil || changed != 1 {
			if err != nil {
				return err
			}
			return sql.ErrNoRows
		}
		if _, err := session.ExecCtx(ctx, `UPDATE crm_schema.admin_refresh_sessions SET revoked_at = CURRENT_TIMESTAMP
			WHERE user_id = $1 AND revoked_at IS NULL`, row.UserID); err != nil {
			return err
		}
		if err := session.QueryRowCtx(ctx, &account, `SELECT id, email, COALESCE(phone_number, '') AS phone_number,
			phone_verified_at IS NOT NULL AS phone_verified, is_active, $2::varchar AS channel
			FROM crm_schema.users WHERE id = $1`, row.UserID, row.Channel); err != nil {
			return err
		}
		_, err = session.ExecCtx(ctx, `INSERT INTO crm_schema.auth_security_events (user_id,event_type,channel,request_ip)
			VALUES ($1,'password_reset.succeeded',$2,$3)`, row.UserID, row.Channel, ip)
		return err
	})
	return account, err
}

func (r *PasswordResetRepository) ChangePhone(ctx context.Context, userID int64, phone, ip string) error {
	return r.connection.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if _, err := session.ExecCtx(ctx, `UPDATE crm_schema.users SET phone_number = NULLIF($2,''),
			phone_verified_at = NULL, updated_at = CURRENT_TIMESTAMP WHERE id = $1`, userID, phone); err != nil {
			return err
		}
		if _, err := session.ExecCtx(ctx, `UPDATE crm_schema.otp_challenge SET consumed_at = CURRENT_TIMESTAMP
			WHERE user_id = $1 AND purpose = 'phone_verify' AND consumed_at IS NULL`, userID); err != nil {
			return err
		}
		_, err := session.ExecCtx(ctx, `INSERT INTO crm_schema.auth_security_events (user_id,event_type,request_ip)
			VALUES ($1,'phone.changed',$2)`, userID, ip)
		return err
	})
}

func (r *PasswordResetRepository) VerifyPhone(ctx context.Context, userID, challengeID int64, destinationHash, phone string) error {
	return r.connection.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var claimedID int64
		err := session.QueryRowCtx(ctx, &claimedID, `UPDATE crm_schema.otp_challenge SET consumed_at = CURRENT_TIMESTAMP
			WHERE id=$1 AND user_id=$2 AND destination_hash=$3 AND channel='whatsapp' AND purpose='phone_verify'
				AND consumed_at IS NULL AND expires_at > CURRENT_TIMESTAMP AND attempts < max_attempts
			RETURNING user_id`, challengeID, userID, destinationHash)
		if errors.Is(err, sqlx.ErrNotFound) {
			return sql.ErrNoRows
		}
		if err != nil {
			return err
		}
		result, err := session.ExecCtx(ctx, `UPDATE crm_schema.users SET phone_verified_at=CURRENT_TIMESTAMP,
			updated_at=CURRENT_TIMESTAMP WHERE id=$1 AND phone_number=$2 AND is_active=TRUE`, claimedID, phone)
		if err != nil {
			return err
		}
		count, err := result.RowsAffected()
		if err != nil || count != 1 {
			return sql.ErrNoRows
		}
		return nil
	})
}

var _ passwordreset.Store = (*PasswordResetRepository)(nil)
