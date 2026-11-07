package handler

import (
	"net/http"

	"webtranshome/internal/modules/catalog/application/dto"
	"webtranshome/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// The frontend's category tree is always at most two levels deep.
// A disabled parent also hides its children from the public catalog.
type publicCategoryNode struct {
	ID       int64                `json:"id"`
	Name     string               `json:"name"`
	Slug     string               `json:"slug"`
	ImageURL *string              `json:"image_url"`
	Children []publicCategoryNode `json:"children,omitempty"`
}

type publicCategoryListResponse struct {
	Data []publicCategoryNode `json:"data"`
}

func PublicCategoriesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := svcCtx.ListCategories.Execute(r.Context())
		if err != nil {
			logx.WithContext(r.Context()).Errorf("list public categories: %v", err)
			httpx.WriteJsonCtx(r.Context(), w, http.StatusInternalServerError, map[string]string{"message": "Gagal memuat kategori."})
			return
		}

		httpx.OkJsonCtx(r.Context(), w, publicCategoryListResponse{Data: publicCategoryTree(items)})
	}
}

func publicCategoryTree(items []dto.Category) []publicCategoryNode {
	roots := make([]publicCategoryNode, 0)
	children := make(map[int64][]publicCategoryNode)
	activeRoots := make(map[int64]bool)

	// Query ListCategories orders by sort_order, name, and id already.
	for _, item := range items {
		if !item.IsActive || item.ParentID != nil {
			continue
		}
		activeRoots[item.ID] = true
		roots = append(roots, publicCategoryNode{
			ID: item.ID, Name: item.Name, Slug: item.Slug,
			ImageURL: publicOptionalImage(item.ImageURL),
			Children: make([]publicCategoryNode, 0),
		})
	}

	for _, item := range items {
		if !item.IsActive || item.ParentID == nil || !activeRoots[*item.ParentID] {
			continue
		}
		children[*item.ParentID] = append(children[*item.ParentID], publicCategoryNode{
			ID: item.ID, Name: item.Name, Slug: item.Slug,
			ImageURL: publicOptionalImage(item.ImageURL),
		})
	}
	for i := range roots {
		if group, ok := children[roots[i].ID]; ok {
			roots[i].Children = group
		}
	}
	return roots
}

func publicOptionalImage(url string) *string {
	if url == "" {
		return nil
	}
	return &url
}