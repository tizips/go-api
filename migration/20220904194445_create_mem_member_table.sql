-- +goose Up
-- +goose StatementBegin

create table mem_member
(
    `id`         varchar(64)       not null,
    `group_id`   int unsigned      not null default 0 comment '级别ID',
    `username`   varchar(32)                default null comment '用户名',
    `mobile`     varchar(20)                default null comment '手机号',
    `email`      varchar(64)                default null comment '邮箱',
    `name`       varchar(32)       not null default '' comment '姓名',
    `avatar`     varchar(255)      not null default '' comment '头像',
    `nickname`   varchar(32)       not null default '' comment '昵称',
    `password`   varchar(64)       not null default '' comment '密码',
    `sex`        tinyint unsigned  not null default 0 comment '性别：0=未知；1=男；2=女',
    `province`   int unsigned      not null default 0 comment '省',
    `city`       int unsigned      not null default 0 comment '市',
    `area`       int unsigned      not null default 0 comment '区',
    `year`       smallint unsigned not null default 0 comment '年',
    `month`      tinyint unsigned  not null default 0 comment '月',
    `day`        tinyint unsigned  not null default 0 comment '日',
    `is_enable`  tinyint unsigned  not null default 0 comment '是否启用：1=是；2=否；',
    `created_at` timestamp         not null default CURRENT_TIMESTAMP,
    `updated_at` timestamp         not null default CURRENT_TIMESTAMP,
    `deleted_at` timestamp                  default null,
    primary key (`id`),
    key (`username`),
    key (`mobile`),
    key (`email`),
    key (`province`),
    key (`city`),
    key (`area`),
    key (`year`)
) default collate = utf8mb4_unicode_ci comment '会员用户表';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists mem_member;

-- +goose StatementEnd
