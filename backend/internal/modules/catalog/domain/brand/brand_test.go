package brand

import (
	"errors"
	"strings"
	"testing"
)

func TestNewBrand(t *testing.T) {
	tests := []struct {
		name      string
		brandName string
		slug      string
		logoURL   string
		wantLogo  string
		wantErr   error
	}{
		{name: "trims name and logo", brandName: " Acme ", slug: "acme", logoURL: " /acme.png ", wantLogo: "/acme.png"},
		{name: "rejects empty name", slug: "acme", wantErr: ErrInvalidBrand},
		{name: "rejects invalid slug", brandName: "Acme", slug: "Acme brand", wantErr: ErrInvalidBrand},
		{name: "rejects name beyond column size", brandName: strings.Repeat("a", 151), slug: "acme", wantErr: ErrInvalidBrand},
		{name: "rejects logo beyond column size", brandName: "Acme", slug: "acme", logoURL: strings.Repeat("a", 256), wantErr: ErrInvalidBrand},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := New(test.brandName, test.slug, test.logoURL)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("New() error = %v, want %v", err, test.wantErr)
			}
			if test.wantErr == nil && (got.Name != "Acme" || !got.IsActive || got.LogoURL == nil || *got.LogoURL != test.wantLogo) {
				t.Fatalf("New() brand = %+v", got)
			}
		})
	}
}

func TestBrandApply(t *testing.T) {
	tests := []struct {
		name    string
		changes Changes
		wantErr error
	}{
		{name: "clears logo", changes: Changes{LogoURL: stringPointer("")}},
		{name: "deactivates", changes: Changes{IsActive: boolPointer(false)}},
		{name: "rejects empty patch", wantErr: ErrInvalidBrand},
		{name: "rejects invalid slug", changes: Changes{Slug: stringPointer("Not valid")}, wantErr: ErrInvalidBrand},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			item := Brand{ID: 2, Name: "Acme", Slug: "acme", IsActive: true, LogoURL: stringPointer("/logo.png")}
			err := item.Apply(test.changes)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("Apply() error = %v, want %v", err, test.wantErr)
			}
			if test.wantErr == nil && test.name == "clears logo" && item.LogoURL != nil {
				t.Fatal("Apply() did not clear logo URL")
			}
			if test.wantErr == nil && test.name == "deactivates" && item.IsActive {
				t.Fatal("Apply() did not deactivate brand")
			}
		})
	}
}

func stringPointer(value string) *string { return &value }
func boolPointer(value bool) *bool       { return &value }
