package postgres

import (
	"errors"
	"testing"

	"webtranshome/internal/modules/catalog/domain/category"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func TestMapCategoryError(t *testing.T) {
	databaseError := errors.New("database unavailable")
	tests := []struct {
		name string
		err  error
		want error
	}{
		{name: "maps missing row", err: sqlx.ErrNotFound, want: category.ErrCategoryNotFound},
		{name: "preserves infrastructure error", err: databaseError, want: databaseError},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := mapCategoryError(test.err); !errors.Is(got, test.want) {
				t.Fatalf("mapCategoryError() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestNullableCategoryFields(t *testing.T) {
	parentID := int64(7)
	imageURL := "/category.png"
	tests := []struct {
		name       string
		item       category.Category
		wantParent bool
		wantImage  bool
	}{
		{name: "maps null fields"},
		{name: "maps populated fields", item: category.Category{ParentID: &parentID, ImageURL: &imageURL}, wantParent: true, wantImage: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parent, image := nullableCategoryFields(test.item)
			if parent.Valid != test.wantParent || image.Valid != test.wantImage {
				t.Fatalf("nullableCategoryFields() valid flags = %t, %t", parent.Valid, image.Valid)
			}
		})
	}
}
