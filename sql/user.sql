create database `user` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
use `user`;

-- 用户基础信息表
CREATE TABLE `userinfo` (
  `user_id` bigint(20) unsigned NOT NULL DEFAULT '0' comment '用户id',
  `nickname` varchar(255) not null DEFAULT '' comment '昵称',
  `gender` tinyint(3) unsigned not null DEFAULT '0',
  `avatar` varchar(255) not null DEFAULT '',
  `birthday` int(11) unsigned not null DEFAULT '0',
  `brief` varchar(500) not null default '' comment '简介',
  `created_on` int(11) unsigned not null default '0',
  `updated_on` int(11) unsigned not null default '0',
  `add_friend_type` tinyint(3) unsigned not null default '1' comment '添加好友的类型,默认是需要验证添加',
  PRIMARY KEY (`user_id`),
  key `idx_created_on` (`created_on`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 登陆的信息表
create table `lastactive_0` (
  `id` bigint(20) unsigned NOT NULL auto_increment,
  `user_id` bigint(20) unsigned NOT NULL DEFAULT '0' comment '用户id',
  `lat` varchar(64) not null default '' comment '经纬度',
  `lng` varchar(64) not null default '' comment '经纬度',
  `ip` varchar(64) not null default '' comment '登陆ip',
  `device` tinyint(3) unsigned not null default 0 comment '设备 ',
  `device_id` varchar(100) not null default '' comment '设备id',
  `app_ver` varchar(32) not null default '' comment 'app版本',
  `created_on` int(11) unsigned not null default '0',
  PRIMARY KEY (`id`),
  key `idx_uid` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

create table `friends` (
  `user_id` bigint(20) NOT NULL DEFAULT '0' comment '用户id',
  `friend_id` bigint(20) NOT NULL DEFAULT '0' comment '好友id',
  `created_on` int(11) not null default '0',
  `markname`  varchar(64) not null default '' comment '备注名称',
  PRIMARY KEY (`user_id`,`friend_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 新申请
create table `new_apply` (
  `id` bigint(20) unsigned NOT NULL auto_increment,
  `a_id` bigint(20) unsigned NOT NULL DEFAULT '0' comment '用户id/群组id',
  `apply_id` bigint(20) unsigned NOT NULL DEFAULT '0' comment '申请用户id',
  `recver_id` bigint(20) unsigned NOT NULL DEFAULT '0' comment '被申请用户id/群组id',
  `apply_type` tinyint(3) unsigned not null default '0' comment '申请类型，群组还是用户',
  `ct` int(5) unsigned not null default 0 comment '申请的次数',
  `ignore` tinyint(1) unsigned not null DEFAULT '0' comment '或略此申请,1就忽略',
  `created_on` int(11) unsigned not null default '0',
  `status` tinyint(3) unsigned not null DEFAULT '0' comment '申请状态',
  `updated_on` int(11) unsigned not null default '0' comment '更新时间',
  `deleted_on` int(11) unsigned not null default '0' comment '删除时间',
  `reason`  varchar(200) not null default '' comment '申请原因',
  `operator_id` bigint(20) unsigned NOT NULL DEFAULT '0' comment '操作用户id',
  `operator_reason` bigint(20) unsigned NOT NULL DEFAULT '0' comment '操作原因',
  PRIMARY KEY (`id`),
  unique KEY (`a_id`,`apply_id`,`recver_id`),
  key idx_uf(`apply_id`,`recver_id`),
  key idx_t(`a_id`,`updated_on`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

create table `blacklist` (
  `user_id` bigint(20) NOT NULL DEFAULT '0' comment '用户id',
  `black_user_id` bigint(20) NOT NULL DEFAULT '0' comment '黑名单用户id',
  `created_on` int(11) not null default '0',
  PRIMARY KEY (`user_id`,`black_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 聊天群组
create table `group` (
  `id` bigint(20) unsigned NOT NULL default 0,
  `t` tinyint(3) unsigned not null default 2 comment '2群组，3聊天室',
  `name` varchar(50) not null default '' comment '群名称',
  `notice` varchar(256) not null default '' comment '群公告',
  `created_on` int(11) unsigned not null default '0',
  `deleted_on` int(11) unsigned not null default '0',
  `max_member_ct` int(11) unsigned not null default '0' comment '最大群人数',
  `current_ct` int(11) unsigned not null default '0' comment '当前人数',
  `avatar` varchar(1000) unsigned not null default '' comment '头像,自定义群头像,默认就是9宫格',
  `join_type` tinyint(3) unsigned not null default 0 comment '加入的枚举类型',
  `mute_util` int(11) unsigned not null default 0 comment '开启禁言到多久结束',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 聊天成员，分表group_id*64
create table `group_member_0` (
  `group_id` bigint(20) unsigned NOT NULL default 0,
  `user_id` bigint(20) unsigned NOT NULL default 0 comment '成员id',
  `user_role` tinyint(3) unsigned not null default 1 comment '用户角色',
  `mark_name` varchar(255) not null default '' comment '群备注',
  `not_disturb` tinyint(1) not null default 0 comment '是否设置了消息免打扰',
  `mute_until` tinyint(3) unsigned not null default 1 comment '被禁止聊天到多久结束',
  `blocked` tinyint(1) unsigned not null default 0 comment '被拉黑',
  `created_on` int(11) unsigned not null default '0',
  PRIMARY KEY (`group_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 群组部分
create table `user_group` (
  `user_id` bigint(20) unsigned NOT NULL default 0 comment '成员id',
  `group_id` bigint(20) unsigned NOT NULL default 0,
  PRIMARY KEY (`user_id`, `group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
