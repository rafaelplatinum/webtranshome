package command

import (
	"context"
	"errors"
	"testing"

	"webtranshome/internal/modules/catalog/domain/brand"
)

type brandRepositoryStub struct {
	item        brand.Brand
	findErr     error
	createErr   error
	updateErr   error
	createCalls int
	updateCalls int
}

func (r *brandRepositoryStub) ListBrands(context.Context) ([]brand.Brand, error) {
	return nil, nil
}

func (r *brandRepositoryStub) FindBrand(context.Context, int64) (brand.Brand, error) {
	return r.item, r.findErr
}

func (r *brandRepositoryStub) CreateBrand(_ context.Context, item brand.Brand) (brand.Brand, error) {
	r.createCalls++
	item.ID = 6
	return item, r.createErr
}

func (r *brandRepositoryStub) UpdateBrand(_ context.Context, item brand.Brand) error {
	r.updateCalls++
	r.item = item
	return r.updateErr
}

func TestCreateBrand(t *testing.T) {
	tests := []struct {
		name      string
		input     CreateBrandInput
		createErr error
		wantErr   error
		wantCalls int
	}{
		{name: "creates brand", input: CreateBrandInput{Name: " Acme ", Slug: "acme"}, wantCalls: 1},
		{name: "rejects invalid brand", input: CreateBrandInput{Name: "Acme", Slug: "Bad Slug"}, wantErr: brand.ErrInvalidBrand},
		{name: "propagates repository error", input: CreateBrandInput{Name: "Acme", Slug: "acme"}, createErr: errors.New("database unavailable"), wantCalls: 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := &brandRepositoryStub{createErr: test.createErr}
			result, err := NewCreateBrand(repository).Execute(context.Background(), test.input)
			wantErr := test.wantErr
			if wantErr == nil {
				wantErr = test.createErr
			}
			if !errors.Is(err, wantErr) {
				t.Fatalf("Execute() error = %v, want %v", err, wantErr)
			}
			if repository.createCalls != test.wantCalls {
				t.Fatalf("CreateBrand() calls = %d, want %d", repository.createCalls, test.wantCalls)
			}
			if wantErr == nil && (result.ID != 6 || result.Name != "Acme") {
				t.Fatalf("created brand = %+v", result)
			}
		})
	}
}

func TestUpdateBrand(t *testing.T) {
	disabled := false
	invalidSlug := "Bad Slug"
	tests := []struct {
		name    string
		id      int64
		changes UpdateBrandInput
		findErr error
		wantErr error
	}{
		{name: "deactivates brand", id: 6, changes: UpdateBrandInput{IsActive: &disabled}},
		{name: "rejects invalid id", wantErr: ErrInvalidBrandID},
		{name: "propagates missing brand", id: 6, findErr: brand.ErrBrandNotFound, wantErr: brand.ErrBrandNotFound},
		{name: "rejects invalid slug", id: 6, changes: UpdateBrandInput{Slug: &invalidSlug}, wantErr: brand.ErrInvalidBrand},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := &brandRepositoryStub{
				item:    brand.Brand{ID: 6, Name: "Acme", Slug: "acme", IsActive: true},
				findErr: test.findErr,
			}
			test.changes.ID = test.id
			_, err := NewUpdateBrand(repository, repository).Execute(context.Background(), test.changes)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("Execute() error = %v, want %v", err, test.wantErr)
			}
			if test.wantErr == nil && repository.updateCalls != 1 {
				t.Fatalf("UpdateBrand() calls = %d, want 1", repository.updateCalls)
			}
			if test.name == "deactivates brand" && repository.item.IsActive {
				t.Fatal("UpdateBrand() did not deactivate the brand")
			}
		})
	}
}

var (
	_ brand.Reader = (*brandRepositoryStub)(nil)
	_ brand.Writer = (*brandRepositoryStub)(nil)
)
