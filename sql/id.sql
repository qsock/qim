create database db_a CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
use db_a;
CREATE TABLE `id` (
  `k` varchar(128) NOT NULL DEFAULT '',
  `id` bigint(20) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`k`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='id表';

CREATE TABLE `keys` (
  `id` int(11) unsigned not null auto_increment,
  `k` varchar(128) NOT NULL DEFAULT '',
  `offset` int(11) unsigned not null default 0,
  `size` int(10) unsigned not null default 0,
  primary key (`id`),
  key idx_k(`k`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='key列表';

create database db_b CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
use db_b;
CREATE TABLE `id` (
  `k` varchar(128) NOT NULL DEFAULT '',
  `id` bigint(20) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`k`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='id表';

create database db_c CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
use db_c;
CREATE TABLE `id` (
  `k` varchar(128) NOT NULL DEFAULT '',
  `id` bigint(20) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`k`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='id表';

create database db_d CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
use db_d;
CREATE TABLE `id` (
  `k` varchar(128) NOT NULL DEFAULT '',
  `id` bigint(20) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`k`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='id表';