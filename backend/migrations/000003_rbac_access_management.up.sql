INSERT INTO crm_schema.menus (name, code, route, icon, sort_order, is_active)
VALUES ('Akses & Peran', 'system.access', NULL, NULL, 0, TRUE)
ON CONFLICT (code) DO NOTHING;

INSERT INTO crm_schema.permissions (menu_id, code, action, description)
SELECT id, 'rbac.manage', 'manage', 'Manage access-control assignments'
FROM crm_schema.menus
WHERE code = 'system.access'
ON CONFLICT (code) DO NOTHING;

INSERT INTO crm_schema.role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM crm_schema.roles r
CROSS JOIN crm_schema.permissions p
WHERE r.code = 'SUPER_ADMIN'
  AND p.code = 'rbac.manage'
ON CONFLICT (role_id, permission_id) DO NOTHING;
