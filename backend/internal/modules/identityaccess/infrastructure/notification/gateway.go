package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Sender interface {
	Send(context.Context, string, string) error
}

type GatewaySender struct {
	endpoint string
	token    string
	client   *http.Client
}

func NewGatewaySender(endpoint, token string) Sender {
	return &GatewaySender{endpoint: strings.TrimSpace(endpoint), token: token, client: &http.Client{Timeout: 12 * time.Second}}
}

func (s *GatewaySender) Send(ctx context.Context, destination, message string) error {
	parsed, err := url.Parse(s.endpoint)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "https" && parsed.Hostname() != "localhost" && parsed.Hostname() != "127.0.0.1") {
		return errors.New("notification gateway URL must use HTTPS")
	}
	if s.token == "" {
		return errors.New("notification gateway token is not configured")
	}
	body, err := json.Marshal(struct {
		To      string `json:"to"`
		Message string `json:"message"`
	}{To: destination, Message: message})
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, s.endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer "+s.token)
	request.Header.Set("Content-Type", "application/json")
	response, err := s.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("notification gateway returned HTTP %d", response.StatusCode)
	}
	return nil
}

// Local fallback intentionally omits message content and OTP digits from logs.
type RedactedDevelopmentSender struct{ Channel string }

func (s RedactedDevelopmentSender) Send(_ context.Context, destination, _ string) error {
	log.Printf("development %s message queued for %s (message content redacted)", s.Channel, maskDestination(destination))
	return nil
}

func maskDestination(value string) string {
	if len(value) < 5 {
		return "***"
	}
	return value[:2] + strings.Repeat("*", len(value)-4) + value[len(value)-2:]
}
