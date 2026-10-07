INSERT INTO crm_schema.menus (name, code, route, icon, sort_order, is_active)
VALUES ('Kategori', 'catalog.categories', '/admin/catalog/categories', 'category', 10, TRUE)
ON CONFLICT (code) DO NOTHING;

INSERT INTO crm_schema.permissions (menu_id, code, action, description)
SELECT id, 'catalog.category.manage', 'manage', 'Manage catalog categories'
FROM crm_schema.menus
WHERE code = 'catalog.categories'
ON CONFLICT (code) DO NOTHING;

INSERT INTO crm_schema.role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM crm_schema.roles r
CROSS JOIN crm_schema.permissions p
WHERE r.code IN ('SUPER_ADMIN', 'ADMIN_KATALOG')
  AND p.code = 'catalog.category.manage'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO crm_schema.menus (name, code, route, icon, sort_order, is_active)
VALUES ('Merek', 'catalog.brands', '/admin/catalog/brands', 'tag', 20, TRUE)
ON CONFLICT (code) DO NOTHING;

INSERT INTO crm_schema.permissions (menu_id, code, action, description)
SELECT id, 'catalog.brand.manage', 'manage', 'Manage catalog brands'
FROM crm_schema.menus
WHERE code = 'catalog.brands'
ON CONFLICT (code) DO NOTHING;

INSERT INTO crm_schema.role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM crm_schema.roles r
CROSS JOIN crm_schema.permissions p
WHERE r.code IN ('SUPER_ADMIN', 'ADMIN_KATALOG')
  AND p.code = 'catalog.brand.manage'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO crm_schema.menus (name, code, route, icon, sort_order, is_active)
VALUES ('Ruangan', 'catalog.rooms', '/admin/catalog/rooms', 'home', 30, TRUE)
ON CONFLICT (code) DO NOTHING;

INSERT INTO crm_schema.permissions (menu_id, code, action, description)
SELECT id, 'catalog.room.manage', 'manage', 'Manage catalog rooms'
FROM crm_schema.menus
WHERE code = 'catalog.rooms'
ON CONFLICT (code) DO NOTHING;

INSERT INTO crm_schema.role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM crm_schema.roles r
CROSS JOIN crm_schema.permissions p
WHERE r.code IN ('SUPER_ADMIN', 'ADMIN_KATALOG')
  AND p.code = 'catalog.room.manage'
ON CONFLICT (role_id, permission_id) DO NOTHING;
