-- ----------------------------
-- 用户表
-- ----------------------------
create table `user`
(
    `id`             bigint       not null auto_increment comment '用户ID',
    `student_id`     varchar(32)  not null comment '学号',
    `password`       varchar(255) not null comment '精弘总密码',
    `phone_num`      varchar(32)  not null default '' comment '手机号(易校园绑定)',

    `user_type`      varchar(32)  not null default '' comment '用户类型',
    `email`          varchar(128) not null default '' comment '邮箱',

    `device_id`      varchar(64)  not null default '' comment '易校园DeviceID',
    `yxy_uid`        varchar(64)  not null default '' comment '易校园UID',

    `zf_password`    varchar(255) not null default '' comment '正方密码',
    `oauth_password` varchar(255) not null default '' comment '统一身份认证密码',

    `create_time`    timestamp(3)  not null default current_timestamp(3) comment '创建时间',
    `update_time`    timestamp(3)  not null default current_timestamp(3) on update current_timestamp(3) comment '更新时间',
    primary key (`id`),
    unique key `uk_student_id` (`student_id`)
) comment = '用户表';

-- ----------------------------
-- 学生信息表
-- 存储学生的身份信息，用于激活验证
-- ----------------------------
create table `student`
(
    `id`          bigint      not null auto_increment comment '记录ID',
    `student_id`  varchar(32) not null comment '学号',
    `id_card`     varchar(64)          default '' comment '身份证号',
    `create_time` timestamp(3) not null default current_timestamp(3) comment '创建时间',
    `update_time` timestamp(3) not null default current_timestamp(3) on update current_timestamp(3) comment '更新时间',
    primary key (`id`),
    unique key `uk_student_id` (`student_id`),
    key `idx_iid` (`id_card`)
) comment = '学生信息表';

create table `college`
(
    `id`           bigint       not null auto_increment comment '学院ID',
    `college_code` varchar(16)  not null comment '学院代码',
    `college_name` varchar(128) not null comment '学院名称',
    `create_time`  timestamp(3)  not null default current_timestamp(3) comment '创建时间',
    `update_time`  timestamp(3)  not null default current_timestamp(3) on update current_timestamp(3) comment '更新时间',
    primary key (`id`),
    unique key `uk_college_code` (`college_code`)
) comment = '学院表';

-- ----------------------------
-- 初始化学院数据
-- ----------------------------
insert into `college` (`college_code`, `college_name`) values
                                                           ('010000', '健行学院'),
                                                           ('020000', '化学工程学院'),
                                                           ('030000', '生物工程学院'),
                                                           ('040000', '药学院、绿色制药协同创新中心'),
                                                           ('050000', '环境学院'),
                                                           ('060000', '材料科学与工程学院'),
                                                           ('070000', '食品科学与工程学院'),
                                                           ('080000', '机械工程学院'),
                                                           ('090000', '信息工程学院'),
                                                           ('100000', '计算机科学与技术学院、软件学院'),
                                                           ('110000', '土木工程学院'),
                                                           ('120000', '能源与碳中和科教融合学院'),
                                                           ('130000', '地理信息学院'),
                                                           ('140000', '物理学院'),
                                                           ('150000', '数学科学学院'),
                                                           ('160000', '管理学院'),
                                                           ('170000', '经济学院'),
                                                           ('180000', '教育学院(职业技术教育学院)'),
                                                           ('190000', '外国语学院'),
                                                           ('200000', '人文学院'),
                                                           ('210000', '设计与建筑学院'),
                                                           ('220000', '法学院'),
                                                           ('230000', '公共管理学院');

-- ----------------------------
-- 小程序快捷登录表，用于提供微信小程序、钉钉小程序等快捷登录方式
-- ----------------------------
create table `mini_program_user`
(
    `id`            bigint       not null auto_increment comment '用户ID',
    `student_id`    varchar(32)  not null comment '学号',
    `app_type`      varchar(32)  not null default '' comment '应用类型',
    `open_id`       varchar(64)  not null default '' comment 'open_id',

    `create_time`    timestamp(3)  not null default current_timestamp(3) comment '创建时间',
    `update_time`    timestamp(3)  not null default current_timestamp(3) on update current_timestamp(3) comment '更新时间',
    primary key (`id`),
    unique key `uk_student_id` (`student_id`)
) comment = '快捷登录表';