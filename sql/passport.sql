create database passport CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
use passport;


-- https://github.com/kakuilan/china_area_mysql.git
-- 暂时把地区库放到passport里面

create table `mobile` (
  `user_id` bigint(20) unsigned NOT NULL,
  `tel` varchar(50) not null default '',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `idx_tel` (`tel`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

create table `wx` (
  `user_id` bigint(20) unsigned NOT NULL,
  `openid` varchar(100) not null default '',
  `unionid` varchar(100) not null default '',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `idx_openid` (`openid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

create table `qq` (
  `user_id` bigint(20) unsigned NOT NULL,
  `openid` varchar(100) not null default '',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `idx_openid` (`openid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

create table `seq` (
  `user_id` bigint(20) unsigned NOT NULL,
  `seq_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 禁止登陆的表
create table `ban` (
  `user_id` bigint(20) unsigned NOT NULL,
  `created_on` int(11) unsigned not null default '0',
  `end_on` int(11) unsigned not null default '0' comment '结束时间',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `qqinfo`(
  `id` int(11) unsigned not null auto_increment,
  `openid` varchar(100) not null default '',
  `nickname`           VARCHAR(50) NOT NULL default '',
  `figureurl`          VARCHAR(150) NOT NULL default '',
  `figureurl_1`        VARCHAR(150) NOT NULL default '',
  `figureurl_2`        VARCHAR(150) NOT NULL default '',
  `figureurl_qq_1`     VARCHAR(150) NOT NULL default '',
  `figureurl_qq_2`     VARCHAR(150) NOT NULL default '',
  `gender`             tinyint(3) unsigned  NOT NULL default '0',
  `is_yellow_vip`      tinyint(3) unsigned   NOT NULL default '0',
  `vip`                tinyint(3) unsigned   NOT NULL default '0',
  `yellow_vip_level`   tinyint(3) unsigned   NOT NULL default '0',
  `level`              tinyint(3) unsigned   NOT NULL default '0',
  `is_yellow_year_vip` tinyint(3) unsigned   NOT NULL default '0',
  `created_on` int(11) unsigned not null default 0,
  primary key (`id`),
  unique key (`openid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `wxinfo` (
  `id` int(11) unsigned not null auto_increment,
  `openid` varchar(100) not null default '',
  `nickname` VARCHAR(50) not null default '',
  `sex` tinyint(3) unsigned  NOT NULL default '0',
  `province` VARCHAR(50) not null default '',
  `city` VARCHAR(50) not null default '',
  `country` VARCHAR(50) not null default '',
  `headimgurl` VARCHAR(200) not null default '',
  `unionid` VARCHAR(100) not null default '',
  `created_on` int(11) unsigned not null default 0,
  `privilege` varchar(500) not null default '',
  primary key (`id`),
  unique key (`openid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
