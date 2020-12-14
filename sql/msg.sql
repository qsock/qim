create database msg_shard0 CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
use msg_shard0;

-- 消息原始库
create table `im_msg_0` (
  `msg_id` bigint(20) unsigned not null comment '消息id',
  `chat_id` varchar(32) not null comment '会话id，分shard的依据',
  `sender_id` bigint(20) unsigned not null comment '发送者id',
  `recver_id` bigint(20) unsigned not null comment '接收者id',
  `chat_type` tinyint(3) unsigned not null default 0 comment '会话类型',
  `msg_type` tinyint(3) unsigned not null default 0 comment '消息类型',
  `status` tinyint(3) unsigned not null default 0 comment '消息状态',
  `created_on` int(11) not null default 0,
  `content` text not null comment '消息内容',
  primary key (`msg_id`),
  key idx_chat_id (`chat_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 系统消息表
create table `sys_msg` (
  `msg_id` bigint(20) unsigned not null comment '消息id',
  `sender_id` bigint(20) unsigned not null comment '发送者id',
  `recver_id` bigint(20) unsigned not null comment '接收者id',
  `created_on` int(11) not null default 0,
  `msg_type` tinyint(3) unsigned not null default 0 comment '消息类型',
  `need_push` tinyint(1) unsigned not null default 1 comment '是否需要推送',
  `send_on` int(11) not null default 0 comment '发送开始时间',
  `ahead_end` int(11) not null default 0 comment '置顶结束时间',
  `status` tinyint(3) unsigned not null default 0 comment '消息状态',
  `content` text not null comment '消息内容',
  primary key (`msg_id`),
  key idx_user_id(`recver_id`),
  key idx_send_on(`send_on`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 会话表
create table `chat_list` (
  `user_id` bigint(20) not null default 0 comment '用户id，也是分shard的依据 user_id%shardCount',
  `chat_id` varchar(32) not null comment '会话id',
  `chat_type` tinyint(3) unsigned not null default 0 comment '会话类型',
  `ahead_on` int(11) not null default 0 comment '置顶时间',
  `created_on` int(11) not null default 0 comment '创建时间',
  `updated_on` int(11) not null default 0 comment '更新时间',
  `deleted_on` int(11) not null default 0 comment '删除时间',
  `unread_ct` int(5) not null default 0 comment '未读数量',
  `is_mute` tinyint(1) not null default 0 comment '是否被静音',
  `read_last_msg_id` bigint(20) not null default 0 comment '已读的最后一条会话信息',
  `last_msg_id` bigint(20) not null default 0 comment '最后一条会话信息',
  primary key(`user_id`,`chat_id`),
  key idx_chat(`chat_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;