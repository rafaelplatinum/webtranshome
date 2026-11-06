-- Nomor WhatsApp situs ditetapkan pemilik proyek dan tidak dapat diubah lewat panel admin.
INSERT INTO crm_schema.site_settings (setting_key, setting_value)
VALUES ('wa_number', '087708777989')
ON CONFLICT (setting_key) DO UPDATE
SET setting_value = EXCLUDED.setting_value,
    updated_at = CURRENT_TIMESTAMP;
