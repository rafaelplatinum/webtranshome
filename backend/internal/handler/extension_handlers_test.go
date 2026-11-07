package handler

import (
	"os"
	"testing"
)

func TestGoogleReturnToAllowlist(t *testing.T) {
	old := os.Getenv("GOOGLE_FRONTEND_ORIGINS")
	defer os.Setenv("GOOGLE_FRONTEND_ORIGINS", old)
	os.Setenv("GOOGLE_FRONTEND_ORIGINS", "http://localhost:5173,http://localhost:5174")
	good := []string{"http://localhost:5173/masuk/google", "http://localhost:5174/masuk/google?asal=%2Fdaftar"}
	bad := []string{"https://evil.example/masuk/google", "http://localhost:5173.evil.example/masuk/google", "javascript:alert(1)", "http://localhost:5173//evil.example", "http://localhost:5173/admin", "http://localhost:5173/masuk/google#token=abc"}
	for _, v := range good {
		if !extAllowedReturn(v) {
			t.Errorf("rejected safe return URL %q", v)
		}
	}
	for _, v := range bad {
		if extAllowedReturn(v) {
			t.Errorf("accepted unsafe return URL %q", v)
		}
	}
}
func TestPhoneValidation(t *testing.T) {
	for _, v := range []string{"628123456789", "08123456789"} {
		if !extPhoneValid(v) {
			t.Errorf("rejected %s", v)
		}
	}
	for _, v := range []string{"+628123456789", "628123", "abc", ""} {
		if extPhoneValid(v) {
			t.Errorf("accepted %s", v)
		}
	}
}