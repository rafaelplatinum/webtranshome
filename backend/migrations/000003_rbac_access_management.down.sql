DELETE FROM crm_schema.role_permissions
WHERE permission_id IN (
    SELECT id FROM crm_schema.permissions WHERE code = 'rbac.manage'
);

DELETE FROM crm_schema.permissions
WHERE code = 'rbac.manage';

DELETE FROM crm_schema.menus
WHERE code = 'system.access';
