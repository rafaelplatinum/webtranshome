package catalog

import (
	"context"

	"webtranshome/internal/modules/catalog/application/command"
	"webtranshome/internal/modules/catalog/application/dto"
	"webtranshome/internal/modules/catalog/domain/brand"
	"webtranshome/internal/svc"
	"webtranshome/internal/types"
)

type BrandLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBrandLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BrandLogic {
	return &BrandLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *BrandLogic) List() (*types.CatalogBrandListResponse, error) {
	items, err := l.svcCtx.ListBrands.Execute(l.ctx)
	if err != nil {
		return nil, err
	}
	result := &types.CatalogBrandListResponse{Brands: make([]types.CatalogBrand, 0, len(items))}
	for _, item := range items {
		result.Brands = append(result.Brands, catalogBrand(item))
	}
	return result, nil
}

func (l *BrandLogic) Get(brandID int64) (*types.CatalogBrand, error) {
	item, err := l.svcCtx.GetBrand.Execute(l.ctx, brandID)
	if err != nil {
		return nil, err
	}
	result := catalogBrand(item)
	return &result, nil
}

func (l *BrandLogic) Create(req *types.CreateBrandRequest) (*types.CatalogBrand, error) {
	if req == nil {
		return nil, brand.ErrInvalidBrand
	}
	item, err := l.svcCtx.CreateBrand.Execute(l.ctx, command.CreateBrandInput{
		Name: req.Name, Slug: req.Slug, LogoURL: req.LogoURL,
	})
	if err != nil {
		return nil, err
	}
	result := catalogBrand(item)
	return &result, nil
}

func (l *BrandLogic) Update(req *types.UpdateBrandRequest, brandID int64) (*types.CatalogBrand, error) {
	if req == nil {
		return nil, brand.ErrInvalidBrand
	}
	item, err := l.svcCtx.UpdateBrand.Execute(l.ctx, command.UpdateBrandInput{
		ID: brandID, Name: req.Name, Slug: req.Slug, LogoURL: req.LogoURL, IsActive: req.IsActive,
	})
	if err != nil {
		return nil, err
	}
	result := catalogBrand(item)
	return &result, nil
}

func catalogBrand(item dto.Brand) types.CatalogBrand {
	return types.CatalogBrand{
		ID: item.ID, Name: item.Name, Slug: item.Slug, LogoURL: item.LogoURL, IsActive: item.IsActive,
	}
}
