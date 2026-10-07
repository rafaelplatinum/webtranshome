package command

import (
	"webtranshome/internal/modules/catalog/application/dto"
	"webtranshome/internal/modules/catalog/domain/brand"
)

func brandDTO(item brand.Brand) dto.Brand {
	logoURL := ""
	if item.LogoURL != nil {
		logoURL = *item.LogoURL
	}
	return dto.Brand{
		ID: item.ID, Name: item.Name, Slug: item.Slug, LogoURL: logoURL, IsActive: item.IsActive,
	}
}
