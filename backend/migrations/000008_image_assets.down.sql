DELETE FROM crm_schema.role_permissions
WHERE permission_id IN (SELECT id FROM crm_schema.permissions WHERE code = 'image.manage');
DELETE FROM crm_schema.permissions WHERE code = 'image.manage';
DELETE FROM crm_schema.menus WHERE code = 'content.images';
DROP TABLE IF EXISTS crm_schema.image_assets;
