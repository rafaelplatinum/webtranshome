DELETE FROM crm_schema.role_permissions
WHERE permission_id IN (
	SELECT id FROM crm_schema.permissions WHERE code = 'catalog.room.manage'
);

DELETE FROM crm_schema.permissions
WHERE code = 'catalog.room.manage';

DELETE FROM crm_schema.menus
WHERE code = 'catalog.rooms';

DELETE FROM crm_schema.role_permissions
WHERE permission_id IN (
	SELECT id FROM crm_schema.permissions WHERE code = 'catalog.brand.manage'
);

DELETE FROM crm_schema.permissions
WHERE code = 'catalog.brand.manage';

DELETE FROM crm_schema.menus
WHERE code = 'catalog.brands';

DELETE FROM crm_schema.role_permissions
WHERE permission_id IN (
	SELECT id FROM crm_schema.permissions WHERE code = 'catalog.category.manage'
);

DELETE FROM crm_schema.permissions
WHERE code = 'catalog.category.manage';

DELETE FROM crm_schema.menus
WHERE code = 'catalog.categories';
