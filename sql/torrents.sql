/*
Navicat MySQL Data Transfer

Source Server         : localhost
Source Server Version : 50505
Source Host           : 127.0.0.1:3306
Source Database       : torrents

Target Server Type    : MYSQL
Target Server Version : 50505
File Encoding         : 65001

Date: 2020-10-14 10:07:41
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for infohashes
-- ----------------------------
DROP TABLE IF EXISTS `infohashes`;
CREATE TABLE `infohashes` (
  `info_hash` varchar(20) NOT NULL DEFAULT '' COMMENT '信息哈希',
  `info_name` varbinary(4096) NOT NULL DEFAULT '' COMMENT '名称',
  `total_length` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '所有文件总大小（字节）',
  `total_files` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '文件总数量',
  `files` longblob NOT NULL COMMENT '文件列表',
  `hot` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '热度',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录添加日期',
  `updated_at` datetime NOT NULL COMMENT '记录最后更新日期',
  UNIQUE KEY `info_hash` (`info_hash`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
