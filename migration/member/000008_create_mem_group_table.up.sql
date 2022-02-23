create table mem_group
(
    `id`         int unsigned     not null auto_increment,
    `code`       varchar(32)               default null comment '代码',
    `name`       varchar(32)      not null default '' comment '名称',
    `is_default` tinyint unsigned not null default 0 comment '是否默认：0=否；1=是',
    `is_enable`  tinyint unsigned not null default 0 comment '是否启用：0=否；1=是',
    `created_at` timestamp        not null default CURRENT_TIMESTAMP,
    `updated_at` timestamp        not null default CURRENT_TIMESTAMP,
    `deleted_at` timestamp                 default null,
    primary key (`id`)
) auto_increment = 10000
  default collate = utf8mb4_unicode_ci comment '会员级别表'