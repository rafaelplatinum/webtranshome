package query

import (
	"context"
	"errors"
	"testing"

	"webtranshome/internal/modules/catalog/domain/category"
)

type categoryReaderStub struct {
	items     []category.Category
	item      category.Category
	err       error
	listCalls int
}

func (r *categoryReaderStub) ListCategories(context.Context) ([]category.Category, error) {
	r.listCalls++
	return r.items, r.err
}

func (r *categoryReaderStub) FindCategory(context.Context, int64) (category.Category, error) {
	return r.item, r.err
}

func (r *categoryReaderStub) DescendantIDs(context.Context, int64) ([]int64, error) {
	return nil, r.err
}

func TestListCategories(t *testing.T) {
	parentID := int64(2)
	tests := []struct {
		name  string
		items []category.Category
		err   error
		want  int
	}{
		{name: "maps categories", items: []category.Category{{ID: 3, Name: "Child", Slug: "child", ParentID: &parentID, IsActive: true}}, want: 1},
		{name: "propagates repository error", err: errors.New("repository unavailable")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader := &categoryReaderStub{items: test.items, err: test.err}
			result, err := NewListCategories(reader).Execute(context.Background())
			if !errors.Is(err, test.err) {
				t.Fatalf("Execute() error = %v, want %v", err, test.err)
			}
			if len(result) != test.want || reader.listCalls != 1 {
				t.Fatalf("Execute() categories = %d, list calls = %d", len(result), reader.listCalls)
			}
			if test.want == 1 && (result[0].ID != 3 || result[0].ParentID == nil || *result[0].ParentID != parentID) {
				t.Fatalf("Execute() category = %+v", result[0])
			}
		})
	}
}

func TestGetCategory(t *testing.T) {
	tests := []struct {
		name      string
		id        int64
		item      category.Category
		err       error
		wantError error
	}{
		{name: "rejects invalid id", id: 0, wantError: ErrInvalidCategoryID},
		{name: "returns category", id: 5, item: category.Category{ID: 5, Name: "Tools", Slug: "tools", IsActive: true}},
		{name: "propagates not found", id: 9, err: category.ErrCategoryNotFound, wantError: category.ErrCategoryNotFound},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader := &categoryReaderStub{item: test.item, err: test.err}
			result, err := NewGetCategory(reader).Execute(context.Background(), test.id)
			if !errors.Is(err, test.wantError) {
				t.Fatalf("Execute() error = %v, want %v", err, test.wantError)
			}
			if test.wantError == nil && result.ID != test.item.ID {
				t.Fatalf("Execute() category = %+v", result)
			}
		})
	}
}

var _ category.Reader = (*categoryReaderStub)(nil)
