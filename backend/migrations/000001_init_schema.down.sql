-- Rollback migration 000001: menghapus semua tabel Muka 1.
-- PERHATIAN: semua data di tabel-tabel ini ikut terhapus.
-- Schema crm_schema sendiri tidak dihapus.

SET search_path TO crm_schema;

DROP TABLE IF EXISTS
    sync_logs,
    qontak_contacts,
    activity_logs,
    site_settings,
    articles,
    banners,
    product_rooms,
    product_images,
    products,
    rooms,
    brands,
    categories,
    point_transactions,
    member_profiles,
    role_permissions,
    user_roles,
    permissions,
    menus,
    roles,
    password_resets,
    users
CASCADE;

DROP FUNCTION IF EXISTS crm_schema.set_updated_at();
