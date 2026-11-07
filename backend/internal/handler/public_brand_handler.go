package handler

import (
	"encoding/json"
	"net/http"

	"webtranshome/internal/modules/catalog/application/dto"
	"webtranshome/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/postgres"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type publicBrandCategory struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type publicBrandItem struct {
	ID           int64                 `json:"id"`
	Name         string                `json:"name"`
	Slug         string                `json:"slug"`
	LogoURL      *string               `json:"logo_url"`
	ProductCount int64                 `json:"product_count"`
	Categories   []publicBrandCategory `json:"categories"`
	ShowOnHome   bool                  `json:"show_on_home"`
	HomeOrder    *int                  `json:"home_order"`
}

type publicBrandListResponse struct {
	Data []publicBrandItem `json:"data"`
}

type publicBrandDetailResponse struct {
	Data publicBrandItem `json:"data"`
}

func PublicBrandsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	conn := postgres.New(svcCtx.Config.Database.DataSource)
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := loadPublicBrands(r, svcCtx, conn)
		if err != nil {
			logx.WithContext(r.Context()).Errorf("list public brands: %v", err)
			httpx.WriteJsonCtx(r.Context(), w, http.StatusInternalServerError, map[string]string{"message": "Gagal memuat brand."})
			return
		}
		httpx.OkJsonCtx(r.Context(), w, publicBrandListResponse{Data: items})
	}
}

func PublicBrandHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	conn := postgres.New(svcCtx.Config.Database.DataSource)
	return func(w http.ResponseWriter, r *http.Request) {
		var path struct {
			Slug string `path:"slug"`
		}
		if err := httpx.ParsePath(r, &path); err != nil || path.Slug == "" {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusBadRequest, map[string]string{"message": "Slug brand tidak valid."})
			return
		}

		items, err := loadPublicBrands(r, svcCtx, conn)
		if err != nil {
			logx.WithContext(r.Context()).Errorf("get public brand: %v", err)
			httpx.WriteJsonCtx(r.Context(), w, http.StatusInternalServerError, map[string]string{"message": "Gagal memuat brand."})
			return
		}
		for _, item := range items {
			if item.Slug == path.Slug {
				httpx.OkJsonCtx(r.Context(), w, publicBrandDetailResponse{Data: item})
				return
			}
		}
		httpx.WriteJsonCtx(r.Context(), w, http.StatusNotFound, map[string]string{"message": "Brand tidak ditemukan."})
	}
}

// Existing query use cases are reused; frontend-specific enrichment stays here.
func loadPublicBrands(r *http.Request, svcCtx *svc.ServiceContext, conn sqlx.SqlConn) ([]publicBrandItem, error) {
	ctx := r.Context()
	brands, err := svcCtx.ListBrands.Execute(ctx)
	if err != nil {
		return nil, err
	}
	categories, err := svcCtx.ListCategories.Execute(ctx)
	if err != nil {
		return nil, err
	}
	products, err := svcCtx.ListProducts.Execute(ctx)
	if err != nil {
		return nil, err
	}

	// Migration 000005 stores the selected homepage brand IDs in site_settings.
	var idsJSON string
	if err := conn.QueryRowCtx(ctx, &idsJSON,
		`SELECT COALESCE((SELECT setting_value FROM crm_schema.site_settings
           WHERE setting_key='home_featured_brand_ids'), '[]')`); err != nil {
		return nil, err
	}
	var selectedIDs []int64
	if err := json.Unmarshal([]byte(idsJSON), &selectedIDs); err != nil {
		return nil, err
	}
	homeOrder := make(map[int64]int)
	for _, id := range selectedIDs {
		if _, exists := homeOrder[id]; !exists && len(homeOrder) < 16 {
			homeOrder[id] = len(homeOrder) + 1
		}
	}

	// Parent categories are returned in the order provided by ListCategories.
	allCategories := make(map[int64]dto.Category, len(categories))
	parents := make([]dto.Category, 0)
	for _, c := range categories {
		allCategories[c.ID] = c
		if c.IsActive && c.ParentID == nil {
			parents = append(parents, c)
		}
	}
	counts := make(map[int64]int64)
	usedCategories := make(map[int64]map[int64]bool)
	for _, p := range products {
		if !p.IsActive || p.BrandID == nil {
			continue
		}
		brandID := *p.BrandID
		counts[brandID]++
		cat, found := allCategories[p.CategoryID]
		if !found || !cat.IsActive {
			continue
		}
		parentID := cat.ID
		if cat.ParentID != nil {
			parentID = *cat.ParentID
		}
		if root, ok := allCategories[parentID]; !ok || !root.IsActive || root.ParentID != nil {
			continue
		}
		if usedCategories[brandID] == nil {
			usedCategories[brandID] = make(map[int64]bool)
		}
		usedCategories[brandID][parentID] = true
	}

	result := make([]publicBrandItem, 0, len(brands))
	for _, b := range brands {
		if !b.IsActive {
			continue
		}
		item := publicBrandItem{
			ID: b.ID, Name: b.Name, Slug: b.Slug,
			LogoURL: publicOptionalImage(b.LogoURL), ProductCount: counts[b.ID],
			Categories: make([]publicBrandCategory, 0),
		}
		if order, ok := homeOrder[b.ID]; ok {
			item.ShowOnHome = true
			item.HomeOrder = &order
		}
		for _, parent := range parents {
			if usedCategories[b.ID][parent.ID] {
				item.Categories = append(item.Categories,
					publicBrandCategory{Slug: parent.Slug, Name: parent.Name})
			}
		}
		result = append(result, item)
	}
	return result, nil
}