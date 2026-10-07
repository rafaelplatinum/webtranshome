package query

import (
	"context"
	"errors"
	"testing"

	"webtranshome/internal/modules/catalog/domain/brand"
)

type brandReaderStub struct {
	items []brand.Brand
	item  brand.Brand
	err   error
}

func (r *brandReaderStub) ListBrands(context.Context) ([]brand.Brand, error) {
	return r.items, r.err
}

func (r *brandReaderStub) FindBrand(context.Context, int64) (brand.Brand, error) {
	return r.item, r.err
}

func TestListBrands(t *testing.T) {
	tests := []struct {
		name  string
		items []brand.Brand
		err   error
		want  int
	}{
		{name: "maps brand", items: []brand.Brand{{ID: 3, Name: "Acme", Slug: "acme", IsActive: true}}, want: 1},
		{name: "propagates repository error", err: errors.New("database unavailable")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := NewListBrands(&brandReaderStub{items: test.items, err: test.err}).Execute(context.Background())
			if !errors.Is(err, test.err) {
				t.Fatalf("Execute() error = %v, want %v", err, test.err)
			}
			if len(result) != test.want {
				t.Fatalf("Execute() brands = %d, want %d", len(result), test.want)
			}
		})
	}
}

func TestGetBrand(t *testing.T) {
	tests := []struct {
		name      string
		id        int64
		item      brand.Brand
		err       error
		wantError error
	}{
		{name: "rejects invalid id", wantError: ErrInvalidBrandID},
		{name: "returns brand", id: 3, item: brand.Brand{ID: 3, Name: "Acme", Slug: "acme"}},
		{name: "propagates not found", id: 7, err: brand.ErrBrandNotFound, wantError: brand.ErrBrandNotFound},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := NewGetBrand(&brandReaderStub{item: test.item, err: test.err}).Execute(context.Background(), test.id)
			if !errors.Is(err, test.wantError) {
				t.Fatalf("Execute() error = %v, want %v", err, test.wantError)
			}
			if test.wantError == nil && result.ID != test.item.ID {
				t.Fatalf("Execute() brand = %+v", result)
			}
		})
	}
}

var _ brand.Reader = (*brandReaderStub)(nil)
