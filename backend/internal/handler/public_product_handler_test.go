package handler

import (
	"encoding/json"
	"net/http/httptest"
	"net/url"
	"testing"

	"webtranshome/internal/modules/catalog/application/dto"
)

func TestPublicProductSpecsHideInternalColumns(t *testing.T) {
	specs := publicProductSpecs(`{"ukuran":"60x60","daya":900,"hpp":5000,"stok_gudang":120,"diskon_proyek":15,"supplier":"internal"}`)
	if specs["ukuran"] != "60x60" || specs["daya"] != "900" {
		t.Fatalf("public specs missing: %v", specs)
	}
	for _, hidden := range []string{"hpp", "stok_gudang", "diskon_proyek", "supplier"} {
		if _, exists := specs[hidden]; exists {
			t.Fatalf("internal spec %s leaked", hidden)
		}
	}
}

func TestPublicProductFilterParsing(t *testing.T) {
	values, err := url.ParseQuery("category=keramik&brand=a,b&stock=HABIS&min=100&max=900&page=2&per_page=200&attr%5Bukuran%5D=60x60,30x30")
	if err != nil {
		t.Fatal(err)
	}
	f, err := parsePublicProductFilters(values)
	if err != nil {
		t.Fatal(err)
	}
	if f.category != "keramik" || f.page != 2 || f.perPage != 60 || !f.brand["a"] || !f.brand["b"] || !f.stock["HABIS"] || !f.attr["ukuran"]["30x30"] {
		t.Fatalf("incorrect filter parsing: %+v", f)
	}
	if f.minimum == nil || *f.minimum != 100 || f.maximum == nil || *f.maximum != 900 {
		t.Fatalf("incorrect price range: %+v", f)
	}
}

func TestPublicProductFilterRejectsBadValues(t *testing.T) {
	for _, raw := range []string{"page=0", "per_page=abc", "min=-1", "min=100&max=5", "sort=unknown", "ids=1,not-id"} {
		values, _ := url.ParseQuery(raw)
		if _, err := parsePublicProductFilters(values); err == nil {
			t.Fatalf("expected error for %q", raw)
		}
	}
}

func TestPublicProductCategoryDescendants(t *testing.T) {
	rootID := int64(1)
	childID := int64(2)
	shot := publicProductSnapshot{categories: []dto.Category{
		{ID: rootID, Slug: "keramik", IsActive: true},
		{ID: childID, Slug: "lantai", ParentID: &rootID, IsActive: true},
		{ID: 3, Slug: "detail", ParentID: &childID, IsActive: true},
		{ID: 4, Slug: "nonaktif", ParentID: &rootID, IsActive: false},
	}}
	ids := shot.selectedCategoryIDs("keramik")
	if !ids[1] || !ids[2] || !ids[3] || ids[4] {
		t.Fatalf("incorrect descendants: %v", ids)
	}
}

func TestPublicProductMediaURL(t *testing.T) {
	r := httptest.NewRequest("GET", "http://localhost:8888/api/v1/products", nil)
	got := publicProductMediaURL(r, "/media/example.webp")
	if got != "http://localhost:8888/media/example.webp" {
		t.Fatalf("got %s", got)
	}
	external := "https://cdn.example.com/a.jpg"
	if publicProductMediaURL(r, external) != external {
		t.Fatal("external URL changed")
	}
}

func TestPublicProductJSONContract(t *testing.T) {
	payload := publicProductResponse{Data: publicProductDetail{
		publicProductItem: publicProductItem{ID: 11, SKU: "KER-001", Name: "Keramik", Slug: "keramik", ShowPrice: true},
		Specifications:    map[string]string{"ukuran": "60x60"},
		Images:            make([]publicProductImage, 0), Rooms: make([]publicProductRef, 0),
	}}
	bytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	var data map[string]map[string]any
	if err := json.Unmarshal(bytes, &data); err != nil {
		t.Fatal(err)
	}
	if data["data"]["slug"] != "keramik" {
		t.Fatalf("product missing: %s", bytes)
	}
	if _, ok := data["data"]["images"]; !ok {
		t.Fatalf("images absent: %s", bytes)
	}
}