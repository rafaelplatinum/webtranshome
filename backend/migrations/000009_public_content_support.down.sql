-- HATI-HATI: rollback menghapus SELURUH data jasa dan field hero.
-- Jangan jalankan pada database berisi data penting tanpa backup.
DROP TABLE IF EXISTS crm_schema.service_categories;
DROP TABLE IF EXISTS crm_schema.services;
ALTER TABLE crm_schema.banners
    DROP COLUMN IF EXISTS subtitle,
    DROP COLUMN IF EXISTS video_url,
    DROP COLUMN IF EXISTS cta_label,
    DROP COLUMN IF EXISTS cta2_label,
    DROP COLUMN IF EXISTS cta2_url;