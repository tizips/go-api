-- +goose Up
-- +goose StatementBegin

create table dor_grant_detail
(
    `id`          int unsigned     not null auto_increment,
    `grant_id`    int unsigned     not null default 0 comment '发放ID',
    `package_id`  int unsigned     not null default 0 comment '打包ID',
    `position_id` int unsigned     not null default 0 comment '位置ID',
    `type_id`     int unsigned     not null default 0 comment '房型ID',
    `building_id` int unsigned     not null default 0 comment '楼栋ID',
    `floor_id`    int unsigned     not null default 0 comment '楼层ID',
    `room_id`     int unsigned     not null default 0 comment '房间ID',
    `bed_id`      int unsigned     not null default 0 comment '床位ID',
    `people_id`   int unsigned     not null default 0 comment '入住人员ID',
    `member_id`   varchar(64)               default null comment '用户ID',
    `device_id`   int unsigned     not null default 0 comment '设备ID',
    `number`      int unsigned     not null default 0 comment '数量',
    `is_public`   tinyint unsigned not null default 0 comment '是否公共设备：1=是；2=否；',
    `created_at`  timestamp        not null default CURRENT_TIMESTAMP,
    `updated_at`  timestamp        not null default CURRENT_TIMESTAMP,
    `deleted_at`  timestamp                 default null,
    primary key (`id`),
    key (`package_id`),
    key (`type_id`),
    key (`building_id`),
    key (`floor_id`),
    key (`room_id`),
    key (`bed_id`),
    key (`device_id`),
    key (`people_id`),
    key (`member_id`)
) auto_increment = 10000
  default collate = utf8mb4_unicode_ci comment '宿舍发放明细表';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists dor_grant_detail;

-- +goose StatementEnd
