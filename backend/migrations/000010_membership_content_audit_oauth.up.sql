-- Apply ONLY after 000001..000009, with a tested local backup.
-- All permissions are additive and scoped to a specific role.
-- Google registration creates an email-only user; phone_number must allow NULL.
ALTER TABLE crm_schema.users ALTER COLUMN phone_number DROP NOT NULL;
ALTER TABLE crm_schema.articles ADD COLUMN IF NOT EXISTS preview_token VARCHAR(64);
CREATE UNIQUE INDEX IF NOT EXISTS uq_articles_preview_token ON crm_schema.articles(preview_token) WHERE preview_token IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_point_transactions_member_created ON crm_schema.point_transactions(member_id,created_at DESC);
CREATE INDEX IF NOT EXISTS idx_activity_logs_created ON crm_schema.activity_logs(created_at DESC,id DESC);
CREATE UNIQUE INDEX IF NOT EXISTS uq_users_google_id_present ON crm_schema.users(google_id) WHERE google_id IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS uq_users_email_lower ON crm_schema.users(lower(email));

INSERT INTO crm_schema.menus(name,code,route,icon,sort_order,is_active) VALUES
('Produk','catalog.products','/admin/produk',NULL,25,TRUE),
('Member','membership.members','/admin/member',NULL,35,TRUE),
('Poin','membership.points','/admin/poin/input',NULL,36,TRUE),
('Banner','content.banners','/admin/banner',NULL,40,TRUE),
('Artikel','content.articles','/admin/artikel',NULL,41,TRUE),
('Jasa','content.services','/admin/jasa',NULL,42,TRUE),
('Pengaturan','content.settings','/admin/pengaturan',NULL,43,TRUE),
('Log Aktivitas','system.logs','/admin/log-aktivitas',NULL,50,TRUE)
ON CONFLICT (code) DO NOTHING;

WITH defs(menu_code,code,action,description) AS (VALUES
('catalog.products','catalog.product.manage','MANAGE','Kelola produk'),
('catalog.products','product.view','READ','Lihat produk'),
('catalog.products','product.create','CREATE','Buat produk'),
('catalog.products','product.update','UPDATE','Ubah produk'),
('catalog.categories','category.view','READ','Lihat kategori'),
('catalog.categories','category.create','CREATE','Buat kategori'),
('catalog.categories','category.update','UPDATE','Ubah kategori'),
('catalog.brands','brand.view','READ','Lihat brand'),
('catalog.brands','brand.create','CREATE','Buat brand'),
('catalog.brands','brand.update','UPDATE','Ubah brand'),
('catalog.rooms','room.view','READ','Lihat ruangan'),
('catalog.rooms','room.create','CREATE','Buat ruangan'),
('catalog.rooms','room.update','UPDATE','Ubah ruangan'),
('membership.members','member.view','READ','Daftar member'),
('membership.members','member.update','UPDATE','Status member'),
('membership.points','point.view','READ','Riwayat poin'),
('membership.points','point.create','CREATE','Input poin'),
('membership.points','point.adjust','UPDATE','Koreksi saldo poin'),
('content.banners','banner.view','READ','Lihat banner'),
('content.banners','banner.create','CREATE','Tambah banner'),
('content.banners','banner.update','UPDATE','Ubah banner'),
('content.articles','article.view','READ','Lihat artikel'),
('content.articles','article.create','CREATE','Tambah artikel'),
('content.articles','article.update','UPDATE','Ubah artikel'),
('content.articles','article.publish','UPDATE','Publikasi artikel'),
('content.services','service.view','READ','Lihat jasa'),
('content.services','service.create','CREATE','Tambah jasa'),
('content.services','service.update','UPDATE','Ubah jasa'),
('content.settings','setting.view','READ','Lihat pengaturan'),
('content.settings','setting.update','UPDATE','Ubah pengaturan'),
('content.settings','setting.point_ratio','UPDATE','Ubah rasio poin'),
('system.logs','log.view','READ','Lihat audit log')
)
INSERT INTO crm_schema.permissions(menu_id,code,action,description)
SELECT m.id,d.code,d.action,d.description FROM defs d JOIN crm_schema.menus m ON m.code=d.menu_code
ON CONFLICT (code) DO NOTHING;

-- Separate roles. SUPER_ADMIN gets all new frontend permissions;
-- ADMIN_KATALOG gets catalog only; MEMBERSHIP gets member/point;
-- ADMIN_KONTEN gets content only. Audit viewing stays superadmin-only.
INSERT INTO crm_schema.role_permissions(role_id,permission_id)
SELECT r.id,p.id FROM crm_schema.roles r CROSS JOIN crm_schema.permissions p
WHERE r.code='SUPER_ADMIN'
   OR (r.code='ADMIN_KATALOG' AND (p.code LIKE 'catalog.%' OR p.code LIKE 'product.%' OR p.code LIKE 'category.%' OR p.code LIKE 'brand.%' OR p.code LIKE 'room.%'))
   OR (r.code='ADMIN_MEMBERSHIP' AND (p.code LIKE 'member.%' OR p.code LIKE 'point.%'))
   OR (r.code='ADMIN_KONTEN' AND (p.code LIKE 'banner.%' OR p.code LIKE 'article.%' OR p.code LIKE 'service.%' OR p.code IN ('setting.view','setting.update')))
ON CONFLICT(role_id,permission_id) DO NOTHING;