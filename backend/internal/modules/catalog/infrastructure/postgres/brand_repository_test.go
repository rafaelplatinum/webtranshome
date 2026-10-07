package postgres

import (
	"errors"
	"testing"

	"webtranshome/internal/modules/catalog/domain/brand"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func TestMapBrandError(t *testing.T) {
	databaseError := errors.New("database unavailable")
	tests := []struct {
		name string
		err  error
		want error
	}{
		{name: "maps missing row", err: sqlx.ErrNotFound, want: brand.ErrBrandNotFound},
		{name: "preserves infrastructure error", err: databaseError, want: databaseError},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := mapBrandError(test.err); !errors.Is(got, test.want) {
				t.Fatalf("mapBrandError() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestNullableBrandLogo(t *testing.T) {
	logo := "/logo.png"
	tests := []struct {
		name  string
		value *string
		want  bool
	}{
		{name: "maps null logo"},
		{name: "maps logo URL", value: &logo, want: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := nullableBrandLogo(test.value); got.Valid != test.want {
				t.Fatalf("nullableBrandLogo().Valid = %t, want %t", got.Valid, test.want)
			}
		})
	}
}
