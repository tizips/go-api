create table oa_department
(
    `id`         int unsigned     not null auto_increment,
    `parent_id`  int unsigned     not null default 0 comment '父级ID',
    `name`       varchar(255)     not null default '' comment '名称',
    `is_enable`  tinyint unsigned not null default 0 comment '是否启用：0=否；1=是',
    `created_at` timestamp        not null default CURRENT_TIMESTAMP,
    `updated_at` timestamp        not null default CURRENT_TIMESTAMP,
    `deleted_at` timestamp                 default null,
    primary key (`id`)
) auto_increment = 10000
  default collate = utf8mb4_unicode_ci comment 'OA架构表'