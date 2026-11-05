package passwordreset

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"
)

type memoryStore struct {
	account      Account
	findErr      error
	counts       RequestCounts
	challenges   []Challenge
	events       []string
	tokenHash    string
	tokenExpires time.Time
	resetUsed    bool
	passwordHash string
}

func (m *memoryStore) FindAccount(context.Context, string, string) (Account, error) {
	return m.account, m.findErr
}
func (m *memoryStore) PhoneForUser(context.Context, int64) (string, error) {
	return m.account.Phone, nil
}
func (m *memoryStore) RequestCounts(context.Context, string, string, time.Time) (RequestCounts, error) {
	return m.counts, nil
}
func (m *memoryStore) VerifyFailureCounts(context.Context, string, string, time.Time) (int, int, error) {
	return 0, 0, nil
}
func (m *memoryStore) CreateChallenge(_ context.Context, c Challenge, _ string) (int64, error) {
	for i := range m.challenges {
		if m.challenges[i].Purpose == c.Purpose && m.challenges[i].Channel == c.Channel {
			m.challenges[i].ConsumedAt = sql.NullTime{Time: time.Now(), Valid: true}
		}
	}
	c.ID = int64(len(m.challenges) + 1)
	m.challenges = append(m.challenges, c)
	return c.ID, nil
}
func (m *memoryStore) LatestChallenge(_ context.Context, hash, channel, purpose string, userID int64) (Challenge, error) {
	for i := len(m.challenges) - 1; i >= 0; i-- {
		c := m.challenges[i]
		if c.DestinationHash == hash && c.Channel == channel && c.Purpose == purpose && (userID == 0 || c.UserID == userID) {
			return c, nil
		}
	}
	return Challenge{}, sql.ErrNoRows
}
func (m *memoryStore) RecordEvent(_ context.Context, _ int64, eventType, _, _, _ string) error {
	m.events = append(m.events, eventType)
	return nil
}
func (m *memoryStore) IncrementFailedAttempt(_ context.Context, id int64) error {
	for i := range m.challenges {
		if m.challenges[i].ID == id {
			m.challenges[i].Attempts++
			if m.challenges[i].Attempts >= m.challenges[i].MaxAttempts {
				m.challenges[i].ConsumedAt = sql.NullTime{Time: time.Now(), Valid: true}
			}
		}
	}
	return nil
}
func (m *memoryStore) IssueResetToken(_ context.Context, challengeID int64, hash, _ string, expiry time.Time) error {
	for i := range m.challenges {
		if m.challenges[i].ID == challengeID && !m.challenges[i].ConsumedAt.Valid {
			m.challenges[i].ConsumedAt = sql.NullTime{Time: time.Now(), Valid: true}
			m.tokenHash = hash
			m.tokenExpires = expiry
			return nil
		}
	}
	return sql.ErrNoRows
}
func (m *memoryStore) ResetPassword(_ context.Context, tokenHash, passwordHash, _, nowStr string) (Account, error) {
	now := time.Now()
	if parsed, err := time.Parse(time.RFC3339Nano, nowStr); err == nil {
		now = parsed
	}
	if m.resetUsed || tokenHash != m.tokenHash || now.After(m.tokenExpires) {
		return Account{}, sql.ErrNoRows
	}
	m.resetUsed = true
	m.passwordHash = passwordHash
	return m.account, nil
}
func (m *memoryStore) ChangePhone(_ context.Context, _ int64, phone, _ string) error {
	m.account.Phone = phone
	m.account.PhoneVerified = false
	return nil
}
func (m *memoryStore) VerifyPhone(_ context.Context, _, _ int64, _, _ string) error {
	m.account.PhoneVerified = true
	return nil
}

type captureSender struct{ messages []string }

func (s *captureSender) Send(_ context.Context, _, message string) error {
	s.messages = append(s.messages, message)
	return nil
}

func newTestService(store Store, sender Sender, now *time.Time) *Service {
	s := New(store, strings.Repeat("s", 32), sender, sender)
	s.asyncSender = false
	s.random = strings.NewReader(strings.Repeat("\x00", 256))
	s.now = func() time.Time { return *now }
	return s
}

func TestRequestForgotUsesSameGenericResponseAndSendsOnlyToVerifiedAccount(t *testing.T) {
	now := time.Date(2026, 11, 5, 10, 0, 0, 0, time.UTC)
	sender := &captureSender{}
	store := &memoryStore{account: Account{ID: 4, Email: "member@example.test", EmailVerified: true, Active: true}}
	s := newTestService(store, sender, &now)
	if err := s.RequestForgot(context.Background(), ChannelEmail, "Member@example.test", "192.0.2.1"); err != nil {
		t.Fatal(err)
	}
	if len(sender.messages) != 1 || !strings.Contains(sender.messages[0], "000000") {
		t.Fatalf("message was not sent: %#v", sender.messages)
	}
	if got := GenericMessage; got != "Jika data terdaftar, kode verifikasi sudah dikirim." {
		t.Fatalf("unexpected public response: %q", got)
	}
	store.findErr = sql.ErrNoRows
	if err := s.RequestForgot(context.Background(), ChannelEmail, "unknown@example.test", "192.0.2.2"); err != nil {
		t.Fatal(err)
	}
	if len(sender.messages) != 1 {
		t.Fatalf("sent a code for unknown account: %d messages", len(sender.messages))
	}
}

func TestVerifyAndResetAreSingleUse(t *testing.T) {
	now := time.Date(2026, 11, 5, 10, 0, 0, 0, time.UTC)
	store := &memoryStore{account: Account{ID: 8, Email: "member@example.test", EmailVerified: true, Channel: ChannelEmail, Active: true}}
	s := newTestService(store, &captureSender{}, &now)
	if err := s.RequestForgot(context.Background(), ChannelEmail, store.account.Email, "192.0.2.1"); err != nil {
		t.Fatal(err)
	}
	token, err := s.VerifyForgot(context.Background(), ChannelEmail, store.account.Email, "000000", "192.0.2.1")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.VerifyForgot(context.Background(), ChannelEmail, store.account.Email, "000000", "192.0.2.1"); !errors.Is(err, ErrInvalidCode) {
		t.Fatalf("second OTP use error=%v", err)
	}
	if err := s.Reset(context.Background(), token, "new-password-123", "192.0.2.1"); err != nil {
		t.Fatal(err)
	}
	if store.passwordHash == "" || store.passwordHash == "new-password-123" {
		t.Fatal("password was not stored as a hash")
	}
	if err := s.Reset(context.Background(), token, "another-password", "192.0.2.1"); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("second reset use error=%v", err)
	}
}

func TestWrongOtpLocksAfterFiveAttemptsAndExpiryIsRejected(t *testing.T) {
	now := time.Date(2026, 11, 5, 10, 0, 0, 0, time.UTC)
	store := &memoryStore{account: Account{ID: 2, Email: "member@example.test", EmailVerified: true, Active: true}}
	s := newTestService(store, &captureSender{}, &now)
	if err := s.RequestForgot(context.Background(), ChannelEmail, store.account.Email, "192.0.2.1"); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 5; i++ {
		_, _ = s.VerifyForgot(context.Background(), ChannelEmail, store.account.Email, "123456", "192.0.2.1")
	}
	if _, err := s.VerifyForgot(context.Background(), ChannelEmail, store.account.Email, "000000", "192.0.2.1"); !errors.Is(err, ErrInvalidCode) {
		t.Fatalf("locked challenge error=%v", err)
	}
	if err := s.RequestForgot(context.Background(), ChannelEmail, store.account.Email, "192.0.2.1"); err != nil {
		t.Fatal(err)
	}
	now = now.Add(6 * time.Minute)
	if _, err := s.VerifyForgot(context.Background(), ChannelEmail, store.account.Email, "000000", "192.0.2.1"); !errors.Is(err, ErrInvalidCode) {
		t.Fatalf("expired challenge error=%v", err)
	}
}

func TestRequestRateLimitAndIndonesianPhoneNormalization(t *testing.T) {
	now := time.Date(2026, 11, 5, 10, 0, 0, 0, time.UTC)
	store := &memoryStore{account: Account{ID: 1, Phone: "6281234567890", PhoneVerified: true, Active: true}, counts: RequestCounts{IPHour: 3}}
	sender := &captureSender{}
	s := newTestService(store, sender, &now)
	if err := s.RequestForgot(context.Background(), ChannelWhatsApp, "0812 3456 7890", "192.0.2.1"); err != nil {
		t.Fatal(err)
	}
	if len(sender.messages) != 0 {
		t.Fatal("rate-limited request sent OTP")
	}
	got, err := normalizePhone("+62 (812) 3456-7890")
	if err != nil || got != "6281234567890" {
		t.Fatalf("normalizePhone() = %q, %v", got, err)
	}
}
