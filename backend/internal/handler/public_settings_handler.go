package handler

import (
	"net/http"
	"strconv"
	"strings"

	"webtranshome/internal/svc"

	"github.com/zeromicro/go-zero/core/stores/postgres"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type publicSettingRow struct {
	Key   string `db:"setting_key"`
	Value string `db:"setting_value"`
}

type publicSiteSettings struct {
	WANumber           *string `json:"wa_number"`
	Address            *string `json:"address"`
	OpeningHours       *string `json:"opening_hours"`
	MapsURL            *string `json:"maps_url"`
	PointRatioRupiah   *int64  `json:"point_ratio_rupiah"`
	WATemplateProduct  *string `json:"wa_template_product"`
	WATemplateCart     *string `json:"wa_template_cart"`
	WATemplateCartItem *string `json:"wa_template_cart_item"`
}

type publicSiteSettingsResponse struct {
	Data publicSiteSettings `json:"data"`
}

func publicSettingString(values map[string]string, key string) *string {
	value := strings.TrimSpace(values[key])
	if value == "" {
		return nil
	}
	return &value
}

func publicSettingWA(raw string) *string {
	// Frontend uses a number without + or separators for wa.me links.
	digits := strings.Builder{}
	for _, ch := range raw {
		if ch >= '0' && ch <= '9' {
			digits.WriteRune(ch)
		}
	}
	value := digits.String()
	if strings.HasPrefix(value, "0") {
		value = "62" + strings.TrimPrefix(value, "0")
	}
	if value == "" {
		return nil
	}
	return &value
}

func PublicSettingsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	db := postgres.New(svcCtx.Config.Database.DataSource)
	return func(w http.ResponseWriter, r *http.Request) {
		// Explicit allowlist prevents accidentally exposing future private
		// secrets stored in site_settings to anonymous visitors.
		var rows []publicSettingRow
		err := db.QueryRowsCtx(r.Context(), &rows, `
            SELECT setting_key,COALESCE(setting_value,'') AS setting_value
            FROM crm_schema.site_settings
            WHERE setting_key IN ('wa_number','address','opening_hours','maps_url',
             'point_ratio_rupiah','wa_template_product','wa_template_cart','wa_template_cart_item')`)
		if err != nil {
			publicContentError(w, r, "pengaturan publik", err)
			return
		}
		values := make(map[string]string, len(rows))
		for _, row := range rows {
			values[row.Key] = row.Value
		}
		var ratio *int64
		if raw := strings.TrimSpace(values["point_ratio_rupiah"]); raw != "" {
			parsed, e := strconv.ParseInt(raw, 10, 64)
			if e == nil && parsed >= 1000 {
				ratio = &parsed
			}
		}
		result := publicSiteSettings{
			WANumber:           publicSettingWA(values["wa_number"]),
			Address:            publicSettingString(values, "address"),
			OpeningHours:       publicSettingString(values, "opening_hours"),
			MapsURL:            publicSettingString(values, "maps_url"),
			PointRatioRupiah:   ratio,
			WATemplateProduct:  publicSettingString(values, "wa_template_product"),
			WATemplateCart:     publicSettingString(values, "wa_template_cart"),
			WATemplateCartItem: publicSettingString(values, "wa_template_cart_item"),
		}
		httpx.OkJsonCtx(r.Context(), w, publicSiteSettingsResponse{Data: result})
	}
}