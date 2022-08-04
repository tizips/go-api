INSERT INTO sys_role
    (`id`, `name`, `summary`)
VALUES (666, '开发组', '系统开发账号所属角色，无需单独权限授权，即可拥有系统所有权限。'),
       (10000, '超级管理员', '');

insert into sys_module
(`id`, `slug`, `name`, `is_enable`, `order`)
values (1, 'site', '站点', 1, 50);


INSERT INTO `sys_permission` (`id`, `module_id`, `parent_i1`, `parent_i2`, `name`, `slug`, `method`, `path`)
VALUES (1, 1, 0, 0, '管理', 'management', '', ''),
       (2, 1, 1, 0, '权限', 'management.permission', '', ''),
       (3, 1, 1, 2, '创建', 'management.permission.create', 'POST', '/admin/site/management/permission'),
       (4, 1, 1, 2, '修改', 'management.permission.update', 'PUT', '/admin/site/management/permissions/:id'),
       (5, 1, 1, 2, '删除', 'management.permission.delete', 'DELETE', '/admin/site/management/permissions/:id'),
       (6, 1, 1, 2, '列表', 'management.permission.tree', 'GET', '/admin/site/management/permissions'),
       (7, 1, 1, 0, '权限', 'management.role', '', ''),
       (8, 1, 1, 7, '创建', 'management.role.create', 'POST', '/admin/site/management/role'),
       (9, 1, 1, 7, '修改', 'management.role.update', 'PUT', '/admin/site/management/roles/:id'),
       (10, 1, 1, 7, '删除', 'management.role.delete', 'DELETE', '/admin/site/management/roles/:id'),
       (11, 1, 1, 7, '列表', 'management.role.paginate', 'GET', '/admin/site/management/roles'),
       (12, 1, 1, 0, '账号', 'management.admin', '', ''),
       (13, 1, 1, 12, '创建', 'management.admin.create', 'POST', '/admin/site/management/admin'),
       (14, 1, 1, 12, '修改', 'management.admin.update', 'PUT', '/admin/site/management/admins/:id'),
       (15, 1, 1, 12, '删除', 'management.admin.delete', 'DELETE', '/admin/site/management/admins/:id'),
       (16, 1, 1, 12, '启禁', 'management.admin.enable', 'PUT', '/admin/site/management/admin/enable'),
       (17, 1, 1, 12, '列表', 'management.admin.paginate', 'GET', '/admin/site/management/admins'),
       (18, 1, 0, 0, '架构', 'architecture', '', ''),
       (19, 1, 18, 0, '模块', 'architecture.module', '', ''),
       (20, 1, 18, 19, '创建', 'architecture.module.create', 'POST', '/admin/site/architecture/module'),
       (21, 1, 18, 19, '修改', 'architecture.module.update', 'PUT', '/admin/site/architecture/modules/:id'),
       (22, 1, 18, 19, '启禁', 'architecture.module.enable', 'PUT', '/admin/site/architecture/module/enable'),
       (23, 1, 18, 19, '删除', 'architecture.module.delete', 'DELETE', '/admin/site/architecture/modules/:id'),
       (24, 1, 18, 19, '列表', 'architecture.module.list', 'GET', '/admin/site/architecture/modules');

