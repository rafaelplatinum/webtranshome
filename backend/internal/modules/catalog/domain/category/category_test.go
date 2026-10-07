package category

import (
	"errors"
	"testing"
)

func TestNewValidatesAndNormalizes(t *testing.T) {
	parentID := int64(4)
	category, err := New("  Electrical  ", " electrical ", &parentID, " ", 3)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if category.Name != "Electrical" || category.Slug != "electrical" ||
		category.ParentID == nil || *category.ParentID != parentID ||
		category.ImageURL != nil || category.SortOrder != 3 || !category.IsActive {
		t.Fatalf("New() = %+v", category)
	}
}

func TestNewRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name         string
		categoryName string
		slug         string
		parentID     *int64
		sortOrder    int64
		wantError    error
	}{
		{name: "empty name", categoryName: " ", slug: "valid", wantError: ErrInvalidCategory},
		{name: "invalid slug", categoryName: "Valid", slug: "Bad Slug", wantError: ErrInvalidCategory},
		{name: "negative sort order", categoryName: "Valid", slug: "valid", sortOrder: -1, wantError: ErrInvalidCategory},
		{name: "sort order exceeds database range", categoryName: "Valid", slug: "valid", sortOrder: maxSortOrder + 1, wantError: ErrInvalidCategory},
		{name: "invalid parent", categoryName: "Valid", slug: "valid", parentID: int64Pointer(0), wantError: ErrInvalidParent},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := New(test.categoryName, test.slug, test.parentID, "", test.sortOrder)
			if !errors.Is(err, test.wantError) {
				t.Fatalf("New() error = %v, want %v", err, test.wantError)
			}
		})
	}
}

func TestApplyChanges(t *testing.T) {
	tests := []struct {
		name        string
		changes     Changes
		descendants []int64
		wantError   error
		wantActive  bool
		wantParent  int64
	}{
		{name: "partial update", changes: Changes{Name: stringPointer(" Updated ")}, wantActive: true},
		{name: "rejects empty patch", wantError: ErrInvalidCategory},
		{name: "deactivate", changes: Changes{IsActive: boolPointer(false)}},
		{name: "clear parent", changes: Changes{ParentID: int64Pointer(0)}, wantActive: true},
		{name: "reject self parent", changes: Changes{ParentID: int64Pointer(8)}, wantError: ErrInvalidParent},
		{name: "reject descendant parent", changes: Changes{ParentID: int64Pointer(9)}, descendants: []int64{9}, wantError: ErrInvalidParent},
		{name: "set parent", changes: Changes{ParentID: int64Pointer(10)}, wantActive: true, wantParent: 10},
		{name: "reject invalid slug", changes: Changes{Slug: stringPointer("Invalid slug")}, wantError: ErrInvalidCategory},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertCategoryChanges(t, test)
		})
	}
}

func assertCategoryChanges(t *testing.T, test struct {
	name        string
	changes     Changes
	descendants []int64
	wantError   error
	wantActive  bool
	wantParent  int64
}) {
	t.Helper()
	item := Category{ID: 8, Name: "Valid", Slug: "valid", IsActive: true}
	err := item.Apply(test.changes, test.descendants)
	if !errors.Is(err, test.wantError) {
		t.Fatalf("Apply() error = %v, want %v", err, test.wantError)
	}
	if test.wantError == nil && (item.IsActive != test.wantActive || !hasParent(item, test.wantParent)) {
		t.Fatalf("Apply() category = %+v", item)
	}
}

func hasParent(item Category, parentID int64) bool {
	if parentID == 0 {
		return item.ParentID == nil
	}
	return item.ParentID != nil && *item.ParentID == parentID
}

func TestApplyChangesClearsExistingParent(t *testing.T) {
	category := Category{ID: 3, Name: "Child", Slug: "child", ParentID: int64Pointer(2), IsActive: true}
	if err := category.Apply(Changes{ParentID: int64Pointer(0)}, nil); err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if category.ParentID != nil {
		t.Fatalf("Apply() parent = %v, want nil", category.ParentID)
	}
}

func int64Pointer(value int64) *int64    { return &value }
func stringPointer(value string) *string { return &value }
func boolPointer(value bool) *bool       { return &value }
