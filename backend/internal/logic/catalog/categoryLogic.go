package catalog

import (
	"context"

	"webtranshome/internal/modules/catalog/application/command"
	"webtranshome/internal/modules/catalog/application/dto"
	"webtranshome/internal/svc"
	"webtranshome/internal/types"
)

type CategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryLogic {
	return &CategoryLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *CategoryLogic) List() (*types.CatalogCategoryListResponse, error) {
	items, err := l.svcCtx.ListCategories.Execute(l.ctx)
	if err != nil {
		return nil, err
	}
	result := &types.CatalogCategoryListResponse{Categories: make([]types.CatalogCategory, 0, len(items))}
	for _, item := range items {
		result.Categories = append(result.Categories, catalogCategory(item))
	}
	return result, nil
}

func (l *CategoryLogic) Get(categoryID int64) (*types.CatalogCategory, error) {
	item, err := l.svcCtx.GetCategory.Execute(l.ctx, categoryID)
	if err != nil {
		return nil, err
	}
	result := catalogCategory(item)
	return &result, nil
}

func (l *CategoryLogic) Create(req *types.CreateCategoryRequest) (*types.CatalogCategory, error) {
	if req == nil {
		return nil, command.ErrInvalidCategory
	}
	item, err := l.svcCtx.CreateCategory.Execute(l.ctx, command.CreateCategoryInput{
		Name:      req.Name,
		Slug:      req.Slug,
		ParentID:  req.ParentID,
		ImageURL:  req.ImageURL,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		return nil, err
	}
	result := catalogCategory(item)
	return &result, nil
}

func (l *CategoryLogic) Update(req *types.UpdateCategoryRequest, categoryID int64) (*types.CatalogCategory, error) {
	if req == nil {
		return nil, command.ErrInvalidCategory
	}
	item, err := l.svcCtx.UpdateCategory.Execute(l.ctx, command.UpdateCategoryInput{
		ID:        categoryID,
		Name:      req.Name,
		Slug:      req.Slug,
		ParentID:  req.ParentID,
		ImageURL:  req.ImageURL,
		SortOrder: req.SortOrder,
		IsActive:  req.IsActive,
	})
	if err != nil {
		return nil, err
	}
	result := catalogCategory(item)
	return &result, nil
}

func catalogCategory(item dto.Category) types.CatalogCategory {
	return types.CatalogCategory{
		ID:        item.ID,
		Name:      item.Name,
		Slug:      item.Slug,
		ParentID:  item.ParentID,
		ImageURL:  item.ImageURL,
		SortOrder: item.SortOrder,
		IsActive:  item.IsActive,
	}
}
