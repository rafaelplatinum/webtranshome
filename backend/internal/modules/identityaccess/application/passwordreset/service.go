package passwordreset

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	ChannelWhatsApp = "whatsapp"
	ChannelEmail    = "email"
	PurposePassword = "password_reset"
	PurposePhone    = "phone_verify"
	GenericMessage  = "Jika data terdaftar, kode verifikasi sudah dikirim."
)

var (
	ErrInvalidInput = errors.New("AUTH_INVALID_INPUT")
	ErrInvalidCode  = errors.New("AUTH_OTP_INVALID")
	ErrInvalidToken = errors.New("AUTH_RESET_TOKEN_INVALID")
	ErrRateLimited  = errors.New("AUTH_RATE_LIMITED")
	ErrNoPhone      = errors.New("AUTH_PHONE_REQUIRED")
	phonePattern    = regexp.MustCompile(`^628[0-9]{7,11}$`)
)

type Account struct {
	ID            int64  `db:"id"`
	Email         string `db:"email"`
	Phone         string `db:"phone_number"`
	Channel       string `db:"channel"`
	PhoneVerified bool   `db:"phone_verified"`
	EmailVerified bool   `db:"email_verified"`
	Active        bool   `db:"is_active"`
}

type Challenge struct {
	ID              int64        `db:"id"`
	UserID          int64        `db:"user_id"`
	Channel         string       `db:"channel"`
	Purpose         string       `db:"purpose"`
	DestinationHash string       `db:"destination_hash"`
	CodeHash        string       `db:"code_hash"`
	ExpiresAt       time.Time    `db:"expires_at"`
	Attempts        int          `db:"attempts"`
	MaxAttempts     int          `db:"max_attempts"`
	ConsumedAt      sql.NullTime `db:"consumed_at"`
	CreatedAt       time.Time    `db:"created_at"`
}

type RequestCounts struct {
	DestinationHour int
	DestinationDay  int
	IPHour          int
	IPDay           int
	LastSentAt      *time.Time
}

type Store interface {
	FindAccount(context.Context, string, string) (Account, error)
	PhoneForUser(context.Context, int64) (string, error)
	RequestCounts(context.Context, string, string, time.Time) (RequestCounts, error)
	VerifyFailureCounts(context.Context, string, string, time.Time) (int, int, error)
	CreateChallenge(context.Context, Challenge, string) (int64, error)
	LatestChallenge(context.Context, string, string, string, int64) (Challenge, error)
	RecordEvent(context.Context, int64, string, string, string, string) error
	IncrementFailedAttempt(context.Context, int64) error
	IssueResetToken(context.Context, int64, string, string, time.Time) error
	ResetPassword(context.Context, string, string, string, string) (Account, error)
	ChangePhone(context.Context, int64, string, string) error
	VerifyPhone(context.Context, int64, int64, string, string) error
}

type Sender interface {
	Send(context.Context, string, string) error
}

type Service struct {
	store       Store
	secret      []byte
	senders     map[string]Sender
	random      io.Reader
	now         func() time.Time
	asyncSender bool
}

func New(store Store, secret string, whatsapp, email Sender) *Service {
	return &Service{store: store, secret: []byte(secret), senders: map[string]Sender{ChannelWhatsApp: whatsapp, ChannelEmail: email}, random: rand.Reader, now: time.Now, asyncSender: true}
}

func (s *Service) RequestForgot(ctx context.Context, channel, destination, ip string) error {
	normalized, err := normalizeDestination(channel, destination)
	if err != nil || len(s.secret) < 32 {
		return ErrInvalidInput
	}
	destinationHash := s.destinationHash(normalized)
	counts, err := s.store.RequestCounts(ctx, destinationHash, ip, s.now())
	if err != nil {
		return err
	}
	limited := counts.DestinationHour >= 3 || counts.DestinationDay >= 10 || counts.IPHour >= 3 || counts.IPDay >= 10
	if counts.LastSentAt != nil && s.now().Sub(*counts.LastSentAt) < time.Minute {
		limited = true
	}
	account, lookupErr := s.store.FindAccount(ctx, channel, normalized)
	if lookupErr != nil && !errors.Is(lookupErr, sql.ErrNoRows) {
		return lookupErr
	}
	if errors.Is(lookupErr, sql.ErrNoRows) || !account.Active ||
		(channel == ChannelWhatsApp && !account.PhoneVerified) ||
		(channel == ChannelEmail && !account.EmailVerified) {
		account = Account{}
	}
	code, err := s.newCode()
	if err != nil {
		return err
	}
	challenge := Challenge{
		UserID: account.ID, Channel: channel, Purpose: PurposePassword,
		DestinationHash: destinationHash, CodeHash: s.codeHash(destinationHash, channel, PurposePassword, code),
		ExpiresAt: s.now().Add(5 * time.Minute), MaxAttempts: 5,
	}
	if limited {
		challenge.UserID = 0
	}
	id, err := s.store.CreateChallenge(ctx, challenge, ip)
	if err != nil {
		return err
	}
	if err := s.store.RecordEvent(ctx, challenge.UserID, "password_reset.requested", channel, destinationHash, ip); err != nil {
		return err
	}
	if !limited && account.ID > 0 {
		dest := normalized
		if s.asyncSender {
			go s.deliver(id, channel, dest, code)
		} else if err := s.send(ctx, channel, dest, code); err != nil {
			_ = s.store.RecordEvent(ctx, account.ID, "otp.delivery_failed", channel, destinationHash, ip)
		}
	}
	return nil
}

func (s *Service) VerifyForgot(ctx context.Context, channel, destination, code, ip string) (string, error) {
	normalized, err := normalizeDestination(channel, destination)
	if err != nil || !sixDigits(code) {
		return "", ErrInvalidCode
	}
	destinationHash := s.destinationHash(normalized)
	failedDest, failedIP, err := s.store.VerifyFailureCounts(ctx, destinationHash, ip, s.now().Add(-time.Minute))
	if err != nil {
		return "", err
	}
	if failedDest >= 10 || failedIP >= 20 {
		return "", ErrRateLimited
	}
	challenge, err := s.store.LatestChallenge(ctx, destinationHash, channel, PurposePassword, 0)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}
	if err != nil || challenge.UserID <= 0 || challenge.ConsumedAt.Valid || challenge.ExpiresAt.Before(s.now()) || challenge.Attempts >= challenge.MaxAttempts {
		if auditErr := s.store.RecordEvent(ctx, 0, "otp.verify.failed", channel, destinationHash, ip); auditErr != nil {
			return "", auditErr
		}
		return "", ErrInvalidCode
	}
	provided := s.codeHash(destinationHash, channel, PurposePassword, code)
	if subtle.ConstantTimeCompare([]byte(provided), []byte(challenge.CodeHash)) != 1 {
		if err := s.store.IncrementFailedAttempt(ctx, challenge.ID); err != nil {
			return "", err
		}
		if err := s.store.RecordEvent(ctx, challenge.UserID, "otp.verify.failed", channel, destinationHash, ip); err != nil {
			return "", err
		}
		return "", ErrInvalidCode
	}
	tokenBytes := make([]byte, 32)
	if _, err := io.ReadFull(s.random, tokenBytes); err != nil {
		return "", err
	}
	resetToken := hex.EncodeToString(tokenBytes)
	tokenHash := sha256.Sum256([]byte(resetToken))
	if err := s.store.IssueResetToken(ctx, challenge.ID, hex.EncodeToString(tokenHash[:]), channel, s.now().Add(10*time.Minute)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrInvalidCode
		}
		return "", err
	}
	if err := s.store.RecordEvent(ctx, challenge.UserID, "password_reset.otp_verified", channel, destinationHash, ip); err != nil {
		return "", err
	}
	return resetToken, nil
}

func (s *Service) Reset(ctx context.Context, token, newPassword, ip string) error {
	if len(newPassword) < 8 || len(newPassword) > 72 || len(token) < 64 {
		return ErrInvalidInput
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	tokenDigest := sha256.Sum256([]byte(token))
	account, err := s.store.ResetPassword(ctx, hex.EncodeToString(tokenDigest[:]), string(hash), ip, s.now().Format(time.RFC3339Nano))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidToken
		}
		return err
	}
	destination := account.Email
	if account.Channel == ChannelWhatsApp {
		destination = account.Phone
	}
	if sender := s.senders[account.Channel]; sender != nil && destination != "" {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			if err := sender.Send(ctx, destination, "Password kamu baru saja diubah. Jika ini bukan kamu, segera hubungi tim Transhome."); err != nil {
				log.Printf("password-change notification failed (user=%d channel=%s)", account.ID, account.Channel)
			}
		}()
	}
	return nil
}

func (s *Service) RequestPhone(ctx context.Context, userID int64, ip string) error {
	phone, err := s.store.PhoneForUser(ctx, userID)
	if err != nil || phone == "" {
		return ErrNoPhone
	}
	normalized, err := normalizeDestination(ChannelWhatsApp, phone)
	if err != nil {
		return ErrNoPhone
	}
	destinationHash := s.destinationHash(normalized)
	counts, err := s.store.RequestCounts(ctx, destinationHash, ip, s.now())
	if err != nil {
		return err
	}
	if counts.DestinationHour >= 3 || counts.DestinationDay >= 10 || counts.IPHour >= 3 || counts.IPDay >= 10 || (counts.LastSentAt != nil && s.now().Sub(*counts.LastSentAt) < time.Minute) {
		return ErrRateLimited
	}
	code, err := s.newCode()
	if err != nil {
		return err
	}
	challenge := Challenge{UserID: userID, Channel: ChannelWhatsApp, Purpose: PurposePhone, DestinationHash: destinationHash,
		CodeHash: s.codeHash(destinationHash, ChannelWhatsApp, PurposePhone, code), ExpiresAt: s.now().Add(5 * time.Minute), MaxAttempts: 5}
	id, err := s.store.CreateChallenge(ctx, challenge, ip)
	if err != nil {
		return err
	}
	if err := s.store.RecordEvent(ctx, userID, "phone_verification.requested", ChannelWhatsApp, destinationHash, ip); err != nil {
		return err
	}
	if s.asyncSender {
		go s.deliver(id, ChannelWhatsApp, normalized, code)
	} else if err := s.send(ctx, ChannelWhatsApp, normalized, code); err != nil {
		return err
	}
	return nil
}

func (s *Service) VerifyPhone(ctx context.Context, userID int64, code, ip string) error {
	phone, err := s.store.PhoneForUser(ctx, userID)
	if err != nil || phone == "" || !sixDigits(code) {
		return ErrInvalidCode
	}
	destinationHash := s.destinationHash(phone)
	failedDest, failedIP, err := s.store.VerifyFailureCounts(ctx, destinationHash, ip, s.now().Add(-time.Minute))
	if err != nil {
		return err
	}
	if failedDest >= 10 || failedIP >= 20 {
		return ErrRateLimited
	}
	challenge, err := s.store.LatestChallenge(ctx, destinationHash, ChannelWhatsApp, PurposePhone, userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if err != nil || challenge.ConsumedAt.Valid || challenge.ExpiresAt.Before(s.now()) || challenge.Attempts >= challenge.MaxAttempts {
		if auditErr := s.store.RecordEvent(ctx, userID, "otp.verify.failed", ChannelWhatsApp, destinationHash, ip); auditErr != nil {
			return auditErr
		}
		return ErrInvalidCode
	}
	provided := s.codeHash(destinationHash, ChannelWhatsApp, PurposePhone, code)
	if subtle.ConstantTimeCompare([]byte(provided), []byte(challenge.CodeHash)) != 1 {
		if err := s.store.IncrementFailedAttempt(ctx, challenge.ID); err != nil {
			return err
		}
		if err := s.store.RecordEvent(ctx, userID, "otp.verify.failed", ChannelWhatsApp, destinationHash, ip); err != nil {
			return err
		}
		return ErrInvalidCode
	}
	if err := s.store.VerifyPhone(ctx, userID, challenge.ID, destinationHash, phone); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidCode
		}
		return err
	}
	if err := s.store.RecordEvent(ctx, userID, "phone_verification.succeeded", ChannelWhatsApp, destinationHash, ip); err != nil {
		return err
	}
	return nil
}

func (s *Service) ChangePhone(ctx context.Context, userID int64, rawPhone, ip string) error {
	phone := ""
	if strings.TrimSpace(rawPhone) != "" {
		var err error
		phone, err = normalizePhone(rawPhone)
		if err != nil {
			return ErrInvalidInput
		}
	}
	if err := s.store.ChangePhone(ctx, userID, phone, ip); err != nil {
		return err
	}
	return nil
}

func (s *Service) destinationHash(destination string) string {
	mac := hmac.New(sha256.New, s.secret)
	_, _ = mac.Write([]byte("destination:v1:" + destination))
	return hex.EncodeToString(mac.Sum(nil))
}

func (s *Service) codeHash(destinationHash, channel, purpose, code string) string {
	mac := hmac.New(sha256.New, s.secret)
	_, _ = mac.Write([]byte("otp:v1:" + destinationHash + ":" + channel + ":" + purpose + ":" + code))
	return hex.EncodeToString(mac.Sum(nil))
}

func (s *Service) newCode() (string, error) {
	n, err := rand.Int(s.random, big.NewInt(1_000_000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func (s *Service) send(ctx context.Context, channel, destination, code string) error {
	sender := s.senders[channel]
	if sender == nil {
		return errors.New("OTP sender is not configured")
	}
	return sender.Send(ctx, destination, fmt.Sprintf("Kode verifikasi Trans Home: %s. Berlaku 5 menit. Jangan bagikan kode ini ke siapa pun.", code))
}

func (s *Service) deliver(challengeID int64, channel, destination, code string) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := s.send(ctx, channel, destination, code); err != nil {
		log.Printf("password OTP delivery failed (challenge=%d channel=%s)", challengeID, channel)
	}
}

func normalizeDestination(channel, destination string) (string, error) {
	destination = strings.TrimSpace(destination)
	switch channel {
	case ChannelEmail:
		address, err := mail.ParseAddress(destination)
		if err != nil || address.Address != destination || len(destination) > 254 {
			return "", ErrInvalidInput
		}
		return strings.ToLower(destination), nil
	case ChannelWhatsApp:
		return normalizePhone(destination)
	default:
		return "", ErrInvalidInput
	}
}

func normalizePhone(raw string) (string, error) {
	text := strings.TrimSpace(raw)
	for _, r := range text {
		if (r < '0' || r > '9') && r != '+' && r != ' ' && r != '-' && r != '(' && r != ')' {
			return "", ErrInvalidInput
		}
	}
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, text)
	switch {
	case strings.HasPrefix(text, "+62") && strings.HasPrefix(digits, "62"):
	case strings.HasPrefix(digits, "08"):
		digits = "62" + digits[1:]
	case strings.HasPrefix(digits, "62"):
	default:
		return "", ErrInvalidInput
	}
	if !phonePattern.MatchString(digits) {
		return "", ErrInvalidInput
	}
	return digits, nil
}

func sixDigits(code string) bool {
	if len(code) != 6 {
		return false
	}
	for _, r := range code {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
