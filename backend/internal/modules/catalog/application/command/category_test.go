package command

import (
	"context"
	"errors"
	"testing"

	"webtranshome/internal/modules/catalog/domain/category"
)

type categoryRepositoryStub struct {
	found       category.Category
	findResults []category.Category
	descendants []int64
	created     category.Category
	updated     category.Category
	findErr     error
	createErr   error
	updateErr   error
	findCalls   int
	createCalls int
	updateCalls int
}

func (r *categoryRepositoryStub) ListCategories(context.Context) ([]category.Category, error) {
	return nil, nil
}

func (r *categoryRepositoryStub) FindCategory(context.Context, int64) (category.Category, error) {
	r.findCalls++
	if r.findCalls <= len(r.findResults) {
		return r.findResults[r.findCalls-1], r.findErr
	}
	return r.found, r.findErr
}

func (r *categoryRepositoryStub) DescendantIDs(context.Context, int64) ([]int64, error) {
	return r.descendants, nil
}

func (r *categoryRepositoryStub) CreateCategory(_ context.Context, item category.Category) (category.Category, error) {
	r.createCalls++
	r.created = item
	item.ID = 20
	return item, r.createErr
}

func (r *categoryRepositoryStub) UpdateCategory(_ context.Context, item category.Category) error {
	r.updateCalls++
	r.updated = item
	return r.updateErr
}

func TestCreateCategory(t *testing.T) {
	tests := []struct {
		name      string
		input     CreateCategoryInput
		parent    category.Category
		findErr   error
		wantErr   error
		wantCalls int
	}{
		{name: "creates root category", input: CreateCategoryInput{Name: " Electrical ", Slug: "electrical"}, wantCalls: 1},
		{name: "rejects invalid category", input: CreateCategoryInput{Name: "Invalid", Slug: "Bad Slug"}, wantErr: category.ErrInvalidCategory},
		{name: "rejects inactive parent", input: CreateCategoryInput{Name: "Child", Slug: "child", ParentID: int64Pointer(4)}, parent: category.Category{ID: 4}, wantErr: category.ErrInvalidParent},
		{name: "propagates parent lookup error", input: CreateCategoryInput{Name: "Child", Slug: "child", ParentID: int64Pointer(4)}, findErr: category.ErrCategoryNotFound, wantErr: category.ErrCategoryNotFound},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := &categoryRepositoryStub{found: test.parent, findErr: test.findErr}
			created, err := NewCreateCategory(repository, repository).Execute(context.Background(), test.input)
			assertCategoryError(t, err, test.wantErr)
			if repository.createCalls != test.wantCalls {
				t.Fatalf("CreateCategory() calls = %d, want %d", repository.createCalls, test.wantCalls)
			}
			if test.wantErr == nil && (created.ID != 20 || created.Name != "Electrical" || created.Slug != "electrical") {
				t.Fatalf("created category = %+v", created)
			}
		})
	}
}

func TestUpdateCategory(t *testing.T) {
	tests := []struct {
		name         string
		id           int64
		changes      category.Changes
		descendants  []int64
		wantErr      error
		wantActive   bool
		parentActive bool
	}{
		{name: "updates selected fields", id: 9, changes: category.Changes{Name: stringPointer(" Updated ")}, wantActive: true},
		{name: "deactivates category", id: 9, changes: category.Changes{IsActive: boolPointer(false)}},
		{name: "rejects invalid id", id: 0, wantErr: ErrInvalidCategoryID},
		{name: "rejects descendant parent", id: 9, changes: category.Changes{ParentID: int64Pointer(12)}, descendants: []int64{12}, wantErr: category.ErrInvalidParent, parentActive: true},
		{name: "rejects inactive new parent", id: 9, changes: category.Changes{ParentID: int64Pointer(12)}, wantErr: category.ErrInvalidParent},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertCategoryUpdate(t, test)
		})
	}
}

func assertCategoryUpdate(t *testing.T, test struct {
	name         string
	id           int64
	changes      category.Changes
	descendants  []int64
	wantErr      error
	wantActive   bool
	parentActive bool
}) {
	t.Helper()
	repository := &categoryRepositoryStub{
		findResults: []category.Category{
			{ID: 9, Name: "Existing", Slug: "existing", IsActive: true},
			{ID: 12, Name: "Parent", Slug: "parent", IsActive: test.parentActive},
		},
		descendants: test.descendants,
	}
	result, err := NewUpdateCategory(repository, repository).Execute(context.Background(), updateInput(test))
	assertCategoryError(t, err, test.wantErr)
	if test.wantErr == nil && (repository.updateCalls != 1 || result.IsActive != test.wantActive) {
		t.Fatalf("updated category = %+v, updates = %d", result, repository.updateCalls)
	}
}

func updateInput(test struct {
	name         string
	id           int64
	changes      category.Changes
	descendants  []int64
	wantErr      error
	wantActive   bool
	parentActive bool
}) UpdateCategoryInput {
	return UpdateCategoryInput{
		ID: test.id, Name: test.changes.Name, Slug: test.changes.Slug,
		ParentID: test.changes.ParentID, ImageURL: test.changes.ImageURL,
		SortOrder: test.changes.SortOrder, IsActive: test.changes.IsActive,
	}
}

func assertCategoryError(t *testing.T, err, want error) {
	t.Helper()
	if !errors.Is(err, want) {
		t.Fatalf("error = %v, want %v", err, want)
	}
}

func int64Pointer(value int64) *int64    { return &value }
func stringPointer(value string) *string { return &value }
func boolPointer(value bool) *bool       { return &value }

var (
	_ category.Reader = (*categoryRepositoryStub)(nil)
	_ category.Writer = (*categoryRepositoryStub)(nil)
)
