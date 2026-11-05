-- Brand pilihan di Beranda: simpan 18 brand rekomendasi awal sebagai daftar ID JSON.
INSERT INTO crm_schema.site_settings (setting_key, setting_value)
SELECT
    'home_featured_brand_ids',
    COALESCE(
        (
            SELECT json_agg(id ORDER BY logo_rank, product_count DESC, name)::text
            FROM (
                SELECT
                    b.id,
                    b.name,
                    CASE WHEN b.logo_url IS NULL THEN 1 ELSE 0 END AS logo_rank,
                    COUNT(p.id)::INTEGER AS product_count
                FROM crm_schema.brands b
                LEFT JOIN crm_schema.products p
                    ON p.brand_id = b.id AND p.is_active = TRUE
                WHERE b.is_active = TRUE
                GROUP BY b.id
                ORDER BY logo_rank, product_count DESC, b.name
                LIMIT 18
            ) AS pilihan
        ),
        '[]'
    )
ON CONFLICT (setting_key) DO NOTHING;
