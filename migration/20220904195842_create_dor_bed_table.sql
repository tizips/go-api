-- +goose Up
-- +goose StatementBegin

create table dor_bed
(
    `id`          int unsigned     not null auto_increment,
    `building_id` int unsigned     not null default 0 comment '楼栋ID',
    `floor_id`    int unsigned     not null default 0 comment '楼层ID',
    `room_id`     int unsigned     not null default 0 comment '房间ID',
    `type_id`     int unsigned     not null default 0 comment '房型ID',
    `bed_id`      int unsigned     not null default 0 comment '房型配置ID',
    `name`        varchar(20)      not null default '' comment '名称',
    `order`       int unsigned     not null default 0 comment '序号：正序',
    `is_enable`   tinyint unsigned not null default 0 comment '是否启用：1=是；2=否；',
    `is_public`   tinyint unsigned not null default 0 comment '是否公共区域：1=是；2=否；',
    `created_at`  timestamp        not null default CURRENT_TIMESTAMP,
    `updated_at`  timestamp        not null default CURRENT_TIMESTAMP,
    `deleted_at`  timestamp                 default null,
    primary key (`id`),
    key (`building_id`),
    key (`floor_id`),
    key (`room_id`),
    key (`type_id`),
    key (`bed_id`)
) auto_increment = 10000
  default collate = utf8mb4_unicode_ci comment '宿舍床位表';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists dor_bed;

-- +goose StatementEnd
