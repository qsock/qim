create database `file` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
use `file`;

CREATE TABLE `file_0` (
  `id` bigint(20) unsigned not null auto_increment,
  `user_id` bigint(20) unsigned NOT NULL,
  `url` varchar(255) not null default '',
  `path` varchar(100) not null default '',
  `created_on` int(11) not null default '0',
  PRIMARY KEY (`id`),
  KEY `idx_uid` (`user_id`,`path`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;