-- DOWN intentionally conserves existing users, points, articles and activity logs.
-- WARNING: revoke grants for new permissions only, then remove them.
DELETE FROM crm_schema.role_permissions WHERE permission_id IN (
 SELECT id FROM crm_schema.permissions WHERE code IN
 ('catalog.product.manage','product.view','product.create','product.update','category.view','category.create','category.update',
 'brand.view','brand.create','brand.update','room.view','room.create','room.update','member.view','member.update',
 'point.view','point.create','point.adjust','banner.view','banner.create','banner.update','article.view','article.create',
 'article.update','article.publish','service.view','service.create','service.update','setting.view','setting.update',
 'setting.point_ratio','log.view'));
DELETE FROM crm_schema.permissions WHERE code IN
 ('catalog.product.manage','product.view','product.create','product.update','category.view','category.create','category.update',
 'brand.view','brand.create','brand.update','room.view','room.create','room.update','member.view','member.update',
 'point.view','point.create','point.adjust','banner.view','banner.create','banner.update','article.view','article.create',
 'article.update','article.publish','service.view','service.create','service.update','setting.view','setting.update',
 'setting.point_ratio','log.view');
-- Keep preview_token column and audit/index data for safe rollback.