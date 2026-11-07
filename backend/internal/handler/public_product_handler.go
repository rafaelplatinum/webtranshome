package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"webtranshome/internal/modules/catalog/application/dto"
	"webtranshome/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/postgres"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Public product response mirrors frontend/fase-0-fondasi/src/services/mock/routes.js.
// This module reuses existing catalog read use cases; no admin API is exposed.
type publicProductRef struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type publicProductCategory struct {
	ID       int64             `json:"id"`
	Name     string            `json:"name"`
	Slug     string            `json:"slug"`
	ParentID *int64            `json:"parent_id"`
	Parent   *publicProductRef `json:"parent,omitempty"`
}

type publicProductImage struct {
	ID        int64  `json:"id" db:"id"`
	ProductID int64  `json:"product_id" db:"product_id"`
	ImageURL  string `json:"image_url" db:"image_url"`
	IsPrimary bool   `json:"is_primary" db:"is_primary"`
	SortOrder int    `json:"sort_order" db:"sort_order"`
}

type publicProductItem struct {
	ID              int64                  `json:"id"`
	SKU             string                 `json:"sku"`
	Name            string                 `json:"name"`
	Slug            string                 `json:"slug"`
	PriceGeneral    float64                `json:"price_general"`
	ShowPrice       bool                   `json:"show_price"`
	UnitSale        string                 `json:"unit_sale"`
	MinOrder        int64                  `json:"min_order"`
	StockStatus     string                 `json:"stock_status"`
	StockQtyLabel   *string                `json:"stock_qty_label"`
	IsFeatured      bool                   `json:"is_featured"`
	Brand           *publicProductRef      `json:"brand"`
	Category        *publicProductCategory `json:"category"`
	PrimaryImageURL *string                `json:"primary_image_url"`
}

type publicProductDetail struct {
	publicProductItem
	Description     string               `json:"description"`
	Specifications  map[string]string    `json:"specifications"`
	DatasheetPDFURL *string              `json:"datasheet_pdf_url"`
	MetaTitle       *string              `json:"meta_title"`
	MetaDescription *string              `json:"meta_description"`
	Images          []publicProductImage `json:"images"`
	Rooms           []publicProductRef   `json:"rooms"`
}

type publicFacetCount struct {
	Slug  string `json:"slug"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type publicStockFacet struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

type publicAttrValue struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

type publicAttrFacet struct {
	Key    string            `json:"key"`
	Values []publicAttrValue `json:"values"`
}

type publicPriceRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

type publicProductFacets struct {
	Brands     []publicFacetCount `json:"brands"`
	Stock      []publicStockFacet `json:"stock"`
	Rooms      []publicFacetCount `json:"rooms"`
	Attributes []publicAttrFacet  `json:"attributes"`
	Price      *publicPriceRange  `json:"price"`
}

type publicProductMeta struct {
	Page       int                 `json:"page"`
	PerPage    int                 `json:"per_page"`
	Total      int                 `json:"total"`
	TotalPages int                 `json:"total_pages"`
	Facets     publicProductFacets `json:"facets"`
}

type publicProductsResponse struct {
	Data []publicProductItem `json:"data"`
	Meta publicProductMeta   `json:"meta"`
}

type publicProductResponse struct {
	Data publicProductDetail `json:"data"`
}

type publicSuggestionsResponse struct {
	Data struct {
		Products   []publicProductRef `json:"products"`
		Categories []publicProductRef `json:"categories"`
		Brands     []publicProductRef `json:"brands"`
	} `json:"data"`
}

type publicProductRoomRow struct {
	ProductID int64 `db:"product_id"`
	RoomID    int64 `db:"room_id"`
}

type publicProductDateRow struct {
	ID        int64     `db:"id"`
	UpdatedAt time.Time `db:"updated_at"`
}

type publicProductSnapshot struct {
	products       []dto.Product
	categories     []dto.Category
	brands         []dto.Brand
	rooms          []dto.Room
	categoriesByID map[int64]dto.Category
	brandsByID     map[int64]dto.Brand
	roomsByID      map[int64]dto.Room
	images         map[int64][]publicProductImage
	roomIDs        map[int64][]int64
	updated        map[int64]time.Time
}

func loadPublicProducts(ctx context.Context, s *svc.ServiceContext, db sqlx.SqlConn) (publicProductSnapshot, error) {
	var shot publicProductSnapshot
	var err error
	if shot.products, err = s.ListProducts.Execute(ctx); err != nil {
		return shot, err
	}
	if shot.categories, err = s.ListCategories.Execute(ctx); err != nil {
		return shot, err
	}
	if shot.brands, err = s.ListBrands.Execute(ctx); err != nil {
		return shot, err
	}
	if shot.rooms, err = s.ListRooms.Execute(ctx); err != nil {
		return shot, err
	}
	shot.categoriesByID = make(map[int64]dto.Category, len(shot.categories))
	shot.brandsByID = make(map[int64]dto.Brand, len(shot.brands))
	shot.roomsByID = make(map[int64]dto.Room, len(shot.rooms))
	for _, c := range shot.categories {
		shot.categoriesByID[c.ID] = c
	}
	for _, b := range shot.brands {
		shot.brandsByID[b.ID] = b
	}
	for _, rm := range shot.rooms {
		shot.roomsByID[rm.ID] = rm
	}

	var images []publicProductImage
	if err = db.QueryRowsCtx(ctx, &images,
		`SELECT id, product_id, image_url, is_primary, sort_order FROM crm_schema.product_images
          ORDER BY product_id, is_primary DESC, sort_order, id`); err != nil {
		return shot, err
	}
	shot.images = make(map[int64][]publicProductImage)
	for _, im := range images {
		shot.images[im.ProductID] = append(shot.images[im.ProductID], im)
	}

	var assigned []publicProductRoomRow
	if err = db.QueryRowsCtx(ctx, &assigned, `SELECT product_id, room_id FROM crm_schema.product_rooms`); err != nil {
		return shot, err
	}
	shot.roomIDs = make(map[int64][]int64)
	for _, row := range assigned {
		shot.roomIDs[row.ProductID] = append(shot.roomIDs[row.ProductID], row.RoomID)
	}

	var dates []publicProductDateRow
	if err = db.QueryRowsCtx(ctx, &dates, `SELECT id, COALESCE(updated_at, created_at, CURRENT_TIMESTAMP) AS updated_at FROM crm_schema.products`); err != nil {
		return shot, err
	}
	shot.updated = make(map[int64]time.Time)
	for _, item := range dates {
		shot.updated[item.ID] = item.UpdatedAt
	}
	return shot, nil
}

// Public URLs stored as /media/<filename> must resolve to the Go server rather
// than to the Vite dev server. Absolute external media URLs are left untouched.
func publicProductMediaURL(r *http.Request, value string) string {
	if value == "" {
		return ""
	}
	parsed, err := url.Parse(value)
	if err != nil || parsed.IsAbs() || strings.HasPrefix(value, "//") {
		return value
	}
	if !strings.HasPrefix(value, "/") {
		return value
	}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + r.Host + value
}

func publicProductOptional(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

// Restrict internal Master Barang keys BEFORE sending public JSON, not only in Vue.
var publicProductPrivateKey = regexp.MustCompile(`(?i)^(hpp|hnet|harga_(beli|pokok|net|proyek)|disk\d*|diskon.*|stok_.*|lokasi_(stok|gudang).*|gudang.*|konsinyasi.*|karantina.*|supplier.*|modal.*|margin.*|cost.*|biaya_.*)$`)

func publicProductSpecs(raw string) map[string]string {
	result := make(map[string]string)
	if raw == "" {
		return result
	}
	var decoded map[string]any
	if err := json.Unmarshal([]byte(raw), &decoded); err != nil {
		return result
	}
	for key, value := range decoded {
		if publicProductPrivateKey.MatchString(key) || value == nil {
			continue
		}
		result[key] = fmt.Sprint(value)
	}
	return result
}

func (shot publicProductSnapshot) publicCategoryVisible(id int64) bool {
	// Walk every ancestor; avoid cycles in a corrupted hierarchy.
	for visited := 0; visited <= len(shot.categories); visited++ {
		cat, ok := shot.categoriesByID[id]
		if !ok || !cat.IsActive {
			return false
		}
		if cat.ParentID == nil {
			return true
		}
		id = *cat.ParentID
	}
	return false
}

func (shot publicProductSnapshot) visible(p dto.Product) bool {
	if !p.IsActive || !shot.publicCategoryVisible(p.CategoryID) {
		return false
	}
	if p.BrandID != nil {
		b, ok := shot.brandsByID[*p.BrandID]
		if ok && !b.IsActive {
			return false
		}
	}
	return true
}

func (shot publicProductSnapshot) summary(r *http.Request, p dto.Product) publicProductItem {
	var brand *publicProductRef
	if p.BrandID != nil {
		if b, ok := shot.brandsByID[*p.BrandID]; ok && b.IsActive {
			brand = &publicProductRef{ID: b.ID, Name: b.Name, Slug: b.Slug}
		}
	}
	var category *publicProductCategory
	if c, ok := shot.categoriesByID[p.CategoryID]; ok {
		category = &publicProductCategory{ID: c.ID, Name: c.Name, Slug: c.Slug, ParentID: c.ParentID}
	}
	var primary *string
	for _, im := range shot.images[p.ID] {
		if im.IsPrimary {
			primary = publicProductOptional(publicProductMediaURL(r, im.ImageURL))
			break
		}
	}
	if primary == nil && len(shot.images[p.ID]) > 0 {
		primary = publicProductOptional(publicProductMediaURL(r, shot.images[p.ID][0].ImageURL))
	}
	return publicProductItem{
		ID: p.ID, SKU: p.SKU, Name: p.Name, Slug: p.Slug,
		PriceGeneral: p.PriceGeneral, ShowPrice: true, UnitSale: p.UnitSale, MinOrder: p.MinOrder,
		StockStatus: p.StockStatus, StockQtyLabel: publicProductOptional(p.StockQtyLabel), IsFeatured: p.IsFeatured,
		Brand: brand, Category: category, PrimaryImageURL: primary,
	}
}

func (shot publicProductSnapshot) detail(r *http.Request, p dto.Product) publicProductDetail {
	item := shot.summary(r, p)
	if item.Category != nil && item.Category.ParentID != nil {
		if parent, ok := shot.categoriesByID[*item.Category.ParentID]; ok {
			item.Category.Parent = &publicProductRef{ID: parent.ID, Name: parent.Name, Slug: parent.Slug}
		}
	}
	images := make([]publicProductImage, 0, len(shot.images[p.ID]))
	for _, im := range shot.images[p.ID] {
		im.ImageURL = publicProductMediaURL(r, im.ImageURL)
		images = append(images, im)
	}
	rooms := make([]publicProductRef, 0)
	for _, id := range shot.roomIDs[p.ID] {
		if rm, ok := shot.roomsByID[id]; ok && rm.IsActive {
			rooms = append(rooms, publicProductRef{ID: rm.ID, Name: rm.Name, Slug: rm.Slug})
		}
	}
	sort.Slice(rooms, func(i, j int) bool { return rooms[i].Name < rooms[j].Name })
	return publicProductDetail{
		publicProductItem: item, Description: p.Description, Specifications: publicProductSpecs(p.Specifications),
		DatasheetPDFURL: publicProductOptional(publicProductMediaURL(r, p.DatasheetPDFURL)),
		MetaTitle:       publicProductOptional(p.MetaTitle), MetaDescription: publicProductOptional(p.MetaDescription),
		Images: images, Rooms: rooms,
	}
}

type publicProductFilters struct {
	query    string
	category string
	brand    map[string]bool
	stock    map[string]bool
	room     map[string]bool
	attr     map[string]map[string]bool
	minimum  *float64
	maximum  *float64
	exclude  int64
	ids      map[int64]bool
	onlyIDs  bool
	featured bool
	sort     string
	page     int
	perPage  int
}

func publicProductCSV(raw string) map[string]bool {
	values := make(map[string]bool)
	for _, item := range strings.Split(raw, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			values[item] = true
		}
	}
	return values
}

func publicProductNumber(query url.Values, key string) (*float64, error) {
	text := strings.TrimSpace(query.Get(key))
	if text == "" {
		return nil, nil
	}
	value, err := strconv.ParseFloat(text, 64)
	if err != nil || math.IsNaN(value) || math.IsInf(value, 0) || value < 0 {
		return nil, fmt.Errorf("parameter %s harus angka nonnegatif", key)
	}
	return &value, nil
}

func parsePublicProductFilters(query url.Values) (publicProductFilters, error) {
	f := publicProductFilters{
		query: strings.ToLower(strings.TrimSpace(query.Get("q"))), category: strings.TrimSpace(query.Get("category")),
		brand: publicProductCSV(query.Get("brand")), stock: publicProductCSV(query.Get("stock")),
		room: publicProductCSV(query.Get("room")), attr: make(map[string]map[string]bool),
		sort: query.Get("sort"), page: 1, perPage: 24,
	}
	if f.sort == "" {
		f.sort = "terbaru"
	}
	if f.sort != "terbaru" && f.sort != "harga_terendah" && f.sort != "harga_tertinggi" && f.sort != "nama" {
		return f, errors.New("sort tidak valid")
	}
	if text := query.Get("page"); text != "" {
		n, err := strconv.Atoi(text)
		if err != nil || n < 1 || n > 100000 {
			return f, errors.New("page tidak valid")
		}
		f.page = n
	}
	if text := query.Get("per_page"); text != "" {
		n, err := strconv.Atoi(text)
		if err != nil || n < 1 {
			return f, errors.New("per_page tidak valid")
		}
		f.perPage = n
	}
	if f.perPage > 60 {
		f.perPage = 60
	}
	if text := query.Get("exclude"); text != "" {
		n, err := strconv.ParseInt(text, 10, 64)
		if err != nil || n < 1 {
			return f, errors.New("exclude tidak valid")
		}
		f.exclude = n
	}
	if text := query.Get("ids"); text != "" {
		f.onlyIDs = true
		f.ids = make(map[int64]bool)
		for _, field := range strings.Split(text, ",") {
			n, err := strconv.ParseInt(strings.TrimSpace(field), 10, 64)
			if err != nil || n < 1 {
				return f, errors.New("ids tidak valid")
			}
			f.ids[n] = true
			if len(f.ids) > 30 {
				return f, errors.New("maksimal 30 ids")
			}
		}
	}
	f.featured = query.Get("featured") == "1"
	for key, values := range query {
		if !strings.HasPrefix(key, "attr[") || !strings.HasSuffix(key, "]") {
			continue
		}
		name := strings.TrimSuffix(strings.TrimPrefix(key, "attr["), "]")
		if name == "" || len(name) > 80 || publicProductPrivateKey.MatchString(name) {
			continue
		}
		for _, value := range values {
			for selected := range publicProductCSV(value) {
				if f.attr[name] == nil {
					f.attr[name] = make(map[string]bool)
				}
				f.attr[name][selected] = true
			}
		}
	}
	var err error
	if f.minimum, err = publicProductNumber(query, "min"); err != nil {
		return f, err
	}
	if f.maximum, err = publicProductNumber(query, "max"); err != nil {
		return f, err
	}
	if f.minimum != nil && f.maximum != nil && *f.minimum > *f.maximum {
		return f, errors.New("min lebih besar dari max")
	}
	return f, nil
}

func (shot publicProductSnapshot) selectedCategoryIDs(slug string) map[int64]bool {
	if slug == "" {
		return nil
	}
	var root int64
	for _, c := range shot.categories {
		if c.IsActive && c.Slug == slug {
			root = c.ID
			break
		}
	}
	selected := make(map[int64]bool)
	if root == 0 {
		return selected
	}
	selected[root] = true
	// Allow any number of category levels; migration currently permits self-reference.
	for n := 0; n < len(shot.categories); n++ {
		changed := false
		for _, c := range shot.categories {
			if c.IsActive && c.ParentID != nil && selected[*c.ParentID] && !selected[c.ID] {
				selected[c.ID] = true
				changed = true
			}
		}
		if !changed {
			break
		}
	}
	return selected
}

func (shot publicProductSnapshot) baseMatch(p dto.Product, f publicProductFilters, categories map[int64]bool) bool {
	if !shot.visible(p) {
		return false
	}
	if categories != nil && !categories[p.CategoryID] {
		return false
	}
	if f.onlyIDs && !f.ids[p.ID] {
		return false
	}
	if f.exclude > 0 && p.ID == f.exclude {
		return false
	}
	if f.featured && !p.IsFeatured {
		return false
	}
	if f.query != "" {
		text := strings.ToLower(p.Name + " " + p.SKU)
		if p.BrandID != nil {
			if b, ok := shot.brandsByID[*p.BrandID]; ok {
				text += " " + strings.ToLower(b.Name)
			}
		}
		if !strings.Contains(text, f.query) {
			return false
		}
	}
	return true
}

func (shot publicProductSnapshot) match(p dto.Product, f publicProductFilters, omit string) bool {
	if omit != "brand" && len(f.brand) > 0 {
		if p.BrandID == nil {
			return false
		}
		b, ok := shot.brandsByID[*p.BrandID]
		if !ok || !f.brand[b.Slug] {
			return false
		}
	}
	if omit != "stock" && len(f.stock) > 0 && !f.stock[p.StockStatus] {
		return false
	}
	if omit != "room" && len(f.room) > 0 {
		found := false
		for _, roomID := range shot.roomIDs[p.ID] {
			rm, ok := shot.roomsByID[roomID]
			if ok && rm.IsActive && f.room[rm.Slug] {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	if len(f.attr) > 0 {
		specs := publicProductSpecs(p.Specifications)
		for name, selected := range f.attr {
			if omit == "attr:"+name || len(selected) == 0 {
				continue
			}
			value, ok := specs[name]
			if !ok || !selected[value] {
				return false
			}
		}
	}
	return true
}

func (shot publicProductSnapshot) facet(base []dto.Product, withoutPrice []dto.Product, f publicProductFilters) publicProductFacets {
	facets := publicProductFacets{
		Brands: make([]publicFacetCount, 0), Stock: make([]publicStockFacet, 0),
		Rooms: make([]publicFacetCount, 0), Attributes: make([]publicAttrFacet, 0),
	}
	count := func(omit string, test func(dto.Product) bool) int {
		result := 0
		for _, p := range base {
			if shot.match(p, f, omit) && test(p) {
				result++
			}
		}
		return result
	}
	for _, b := range shot.brands {
		if !b.IsActive {
			continue
		}
		bc := b
		n := count("brand", func(p dto.Product) bool { return p.BrandID != nil && *p.BrandID == bc.ID })
		if n > 0 || f.brand[b.Slug] {
			facets.Brands = append(facets.Brands, publicFacetCount{Slug: b.Slug, Name: b.Name, Count: n})
		}
	}
	for _, status := range []string{"TERSEDIA", "SISA_STOK", "PRE_ORDER", "HABIS"} {
		st := status
		n := count("stock", func(p dto.Product) bool { return p.StockStatus == st })
		facets.Stock = append(facets.Stock, publicStockFacet{Value: status, Count: n})
	}
	for _, rm := range shot.rooms {
		if !rm.IsActive {
			continue
		}
		room := rm
		n := count("room", func(p dto.Product) bool {
			for _, id := range shot.roomIDs[p.ID] {
				if id == room.ID {
					return true
				}
			}
			return false
		})
		if n > 0 || f.room[rm.Slug] {
			facets.Rooms = append(facets.Rooms, publicFacetCount{Slug: rm.Slug, Name: rm.Name, Count: n})
		}
	}
	// Public filter attributes match the categories documented by the frontend.
	for _, key := range []string{"ukuran", "finishing", "warna", "bahan", "daya"} {
		values := make(map[string]bool)
		for _, p := range base {
			if value, ok := publicProductSpecs(p.Specifications)[key]; ok && value != "" {
				values[value] = true
			}
		}
		for selected := range f.attr[key] {
			values[selected] = true
		}
		keys := make([]string, 0, len(values))
		for key := range values {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		group := publicAttrFacet{Key: key, Values: make([]publicAttrValue, 0)}
		for _, value := range keys {
			selectedValue := value
			n := count("attr:"+key, func(p dto.Product) bool {
				raw, ok := publicProductSpecs(p.Specifications)[key]
				return ok && raw == selectedValue
			})
			if n > 0 || f.attr[key][value] {
				group.Values = append(group.Values, publicAttrValue{Value: value, Count: n})
			}
		}
		if len(group.Values) >= 2 || len(f.attr[key]) > 0 {
			facets.Attributes = append(facets.Attributes, group)
		}
	}
	for _, p := range withoutPrice {
		if !shot.match(p, f, "") {
			continue
		}
		if facets.Price == nil {
			facets.Price = &publicPriceRange{Min: p.PriceGeneral, Max: p.PriceGeneral}
			continue
		}
		facets.Price.Min = math.Min(facets.Price.Min, p.PriceGeneral)
		facets.Price.Max = math.Max(facets.Price.Max, p.PriceGeneral)
	}
	return facets
}

func (shot publicProductSnapshot) ordered(items []dto.Product, order string) {
	sort.SliceStable(items, func(i, j int) bool {
		a, b := items[i], items[j]
		switch order {
		case "nama":
			if a.Name != b.Name {
				return strings.ToLower(a.Name) < strings.ToLower(b.Name)
			}
		case "harga_terendah":
			if a.PriceGeneral != b.PriceGeneral {
				return a.PriceGeneral < b.PriceGeneral
			}
		case "harga_tertinggi":
			if a.PriceGeneral != b.PriceGeneral {
				return a.PriceGeneral > b.PriceGeneral
			}
		default:
			if !shot.updated[a.ID].Equal(shot.updated[b.ID]) {
				return shot.updated[a.ID].After(shot.updated[b.ID])
			}
		}
		return a.ID > b.ID
	})
}

func PublicProductsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	db := postgres.New(svcCtx.Config.Database.DataSource)
	return func(w http.ResponseWriter, r *http.Request) {
		filters, err := parsePublicProductFilters(r.URL.Query())
		if err != nil {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusBadRequest, map[string]string{"message": err.Error()})
			return
		}
		shot, err := loadPublicProducts(r.Context(), svcCtx, db)
		if err != nil {
			publicProductsFail(w, r, err)
			return
		}
		catIDs := shot.selectedCategoryIDs(filters.category)
		withoutPrice := make([]dto.Product, 0)
		base := make([]dto.Product, 0)
		matched := make([]dto.Product, 0)
		for _, p := range shot.products {
			if !shot.baseMatch(p, filters, catIDs) {
				continue
			}
			withoutPrice = append(withoutPrice, p)
			if filters.minimum != nil && p.PriceGeneral < *filters.minimum {
				continue
			}
			if filters.maximum != nil && p.PriceGeneral > *filters.maximum {
				continue
			}
			base = append(base, p)
			if shot.match(p, filters, "") {
				matched = append(matched, p)
			}
		}
		facets := shot.facet(base, withoutPrice, filters)
		shot.ordered(matched, filters.sort)
		total := len(matched)
		pages := int(math.Max(1, math.Ceil(float64(total)/float64(filters.perPage))))
		if filters.page > pages {
			filters.page = pages
		}
		start := (filters.page - 1) * filters.perPage
		end := start + filters.perPage
		if end > total {
			end = total
		}
		result := make([]publicProductItem, 0, end-start)
		for _, p := range matched[start:end] {
			result = append(result, shot.summary(r, p))
		}
		httpx.OkJsonCtx(r.Context(), w, publicProductsResponse{
			Data: result, Meta: publicProductMeta{
				Page: filters.page, PerPage: filters.perPage, Total: total, TotalPages: pages, Facets: facets,
			},
		})
	}
}

func PublicProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	db := postgres.New(svcCtx.Config.Database.DataSource)
	return func(w http.ResponseWriter, r *http.Request) {
		var path struct {
			Slug string `path:"slug"`
		}
		if err := httpx.ParsePath(r, &path); err != nil || strings.TrimSpace(path.Slug) == "" {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusBadRequest, map[string]string{"message": "Slug produk tidak valid."})
			return
		}
		shot, err := loadPublicProducts(r.Context(), svcCtx, db)
		if err != nil {
			publicProductsFail(w, r, err)
			return
		}
		for _, p := range shot.products {
			if p.Slug == path.Slug && shot.visible(p) {
				httpx.OkJsonCtx(r.Context(), w, publicProductResponse{Data: shot.detail(r, p)})
				return
			}
		}
		httpx.WriteJsonCtx(r.Context(), w, http.StatusNotFound, map[string]string{"message": "Produk tidak ditemukan."})
	}
}

func PublicProductSuggestHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	db := postgres.New(svcCtx.Config.Database.DataSource)
	return func(w http.ResponseWriter, r *http.Request) {
		term := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
		response := publicSuggestionsResponse{}
		response.Data.Products = make([]publicProductRef, 0)
		response.Data.Categories = make([]publicProductRef, 0)
		response.Data.Brands = make([]publicProductRef, 0)
		if len([]rune(term)) < 2 {
			httpx.OkJsonCtx(r.Context(), w, response)
			return
		}
		shot, err := loadPublicProducts(r.Context(), svcCtx, db)
		if err != nil {
			publicProductsFail(w, r, err)
			return
		}
		for _, p := range shot.products {
			if len(response.Data.Products) >= 5 {
				break
			}
			if shot.visible(p) && (strings.Contains(strings.ToLower(p.Name), term) || strings.Contains(strings.ToLower(p.SKU), term)) {
				response.Data.Products = append(response.Data.Products, publicProductRef{ID: p.ID, Name: p.Name, Slug: p.Slug})
			}
		}
		for _, c := range shot.categories {
			if len(response.Data.Categories) >= 3 {
				break
			}
			if shot.publicCategoryVisible(c.ID) && strings.Contains(strings.ToLower(c.Name), term) {
				response.Data.Categories = append(response.Data.Categories, publicProductRef{ID: c.ID, Name: c.Name, Slug: c.Slug})
			}
		}
		for _, b := range shot.brands {
			if len(response.Data.Brands) >= 3 {
				break
			}
			if b.IsActive && strings.Contains(strings.ToLower(b.Name), term) {
				response.Data.Brands = append(response.Data.Brands, publicProductRef{ID: b.ID, Name: b.Name, Slug: b.Slug})
			}
		}
		httpx.OkJsonCtx(r.Context(), w, response)
	}
}

func publicProductsFail(w http.ResponseWriter, r *http.Request, err error) {
	// Deliberately do not reveal SQL or database details in public HTTP responses.
	logx.WithContext(r.Context()).Errorf("public product query: %v", err)
	httpx.WriteJsonCtx(r.Context(), w, http.StatusInternalServerError, map[string]string{"message": "Gagal memuat katalog produk."})
}

// The product catalog schema currently stores price_general but no show_price
// flag. Therefore, ShowPrice=true for all published rows until that policy is
// added at the database and Admin API levels.