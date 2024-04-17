/*
 Navicat Premium Data Transfer

 Source Server         : dev
 Source Server Type    : MySQL
 Source Server Version : 80200
 Source Host           : localhost:3306
 Source Schema         : cron

 Target Server Type    : MySQL
 Target Server Version : 80200
 File Encoding         : 65001

 Date: 17/04/2024 12:45:42
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for log
-- ----------------------------
DROP TABLE IF EXISTS `log`;
CREATE TABLE `log` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `uuid` varchar(64) NOT NULL COMMENT '标识',
  `worker_id` int unsigned NOT NULL COMMENT '节点id',
  `worker_snapshot` varchar(2048) NOT NULL COMMENT '节点快照',
  `task_id` int unsigned NOT NULL COMMENT '任务id',
  `task_snapshot` varchar(2048) NOT NULL COMMENT '任务快照',
  `start` int unsigned NOT NULL COMMENT '开始时间',
  `end` int unsigned DEFAULT NULL COMMENT '结束时间',
  `content` text NOT NULL COMMENT '执行日主',
  `status` char(32) NOT NULL COMMENT '执行结果',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`)
) ENGINE=InnoDB AUTO_INCREMENT=224 DEFAULT CHARSET=utf8mb4  COMMENT='日志信息';

-- ----------------------------
-- Table structure for task
-- ----------------------------
DROP TABLE IF EXISTS `task`;
CREATE TABLE `task` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `group_id` int unsigned DEFAULT NULL COMMENT '分组id',
  `name` varchar(64) NOT NULL COMMENT '名称',
  `tag` varchar(32) NOT NULL COMMENT '标签',
  `spec` varchar(32) NOT NULL COMMENT '表达式',
  `select_type` varchar(32) NOT NULL COMMENT '选择类型',
  `worker_group_id` int unsigned DEFAULT NULL COMMENT '节点分组id',
  `worker_id` int unsigned DEFAULT NULL COMMENT '节点id',
  `exec_type` varchar(32) NOT NULL COMMENT '执行类型',
  `exec_value` text NOT NULL COMMENT '执行内容',
  `expect_code` int NOT NULL COMMENT '预期状态码',
  `retry_count` int NOT NULL COMMENT '重试次数',
  `retry_wait_time` int NOT NULL COMMENT '等待时长',
  `max_exec_time` int NOT NULL COMMENT '执行市场',
  `status` tinyint(1) NOT NULL COMMENT '状态',
  `version` varchar(64) NOT NULL COMMENT '版本',
  `description` varchar(256) NOT NULL COMMENT '描述',
  `created_at` bigint DEFAULT NULL COMMENT '创建时间',
  `updated_at` bigint DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `version` (`version`),
  KEY `idx_task_updated_at` (`updated_at`),
  KEY `idx_task_created_at` (`created_at`),
  KEY `group_id` (`group_id`),
  KEY `worker_group_id` (`worker_group_id`),
  KEY `worker_id` (`worker_id`),
  CONSTRAINT `task_ibfk_1` FOREIGN KEY (`group_id`) REFERENCES `task_group` (`id`),
  CONSTRAINT `task_ibfk_2` FOREIGN KEY (`worker_group_id`) REFERENCES `worker_group` (`id`),
  CONSTRAINT `task_ibfk_3` FOREIGN KEY (`worker_id`) REFERENCES `worker` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4  COMMENT='任务信息';

-- ----------------------------
-- Table structure for task_group
-- ----------------------------
DROP TABLE IF EXISTS `task_group`;
CREATE TABLE `task_group` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(32) NOT NULL COMMENT '名称',
  `description` varchar(256) NOT NULL COMMENT '描述',
  `created_at` bigint DEFAULT NULL COMMENT '创建时间',
  `updated_at` bigint DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  KEY `idx_task_group_updated_at` (`updated_at`),
  KEY `idx_task_group_created_at` (`created_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4  COMMENT='任务分组';

-- ----------------------------
-- Table structure for worker
-- ----------------------------
DROP TABLE IF EXISTS `worker`;
CREATE TABLE `worker` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(32) NOT NULL COMMENT '名称',
  `ip` varchar(64) NOT NULL COMMENT 'ip',
  `group_id` int unsigned DEFAULT NULL COMMENT '分组',
  `status` tinyint(1) NOT NULL COMMENT '状态',
  `description` varchar(256) NOT NULL COMMENT '描述',
  `created_at` bigint DEFAULT NULL COMMENT '创建时间',
  `updated_at` bigint DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `ip` (`ip`),
  KEY `idx_worker_updated_at` (`updated_at`),
  KEY `idx_worker_created_at` (`created_at`),
  KEY `group_id` (`group_id`),
  CONSTRAINT `worker_ibfk_1` FOREIGN KEY (`group_id`) REFERENCES `worker_group` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4  COMMENT='节点信息';

-- ----------------------------
-- Table structure for worker_group
-- ----------------------------
DROP TABLE IF EXISTS `worker_group`;
CREATE TABLE `worker_group` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(32) NOT NULL COMMENT '名称',
  `description` varchar(256) NOT NULL COMMENT '描述',
  `created_at` bigint DEFAULT NULL COMMENT '创建时间',
  `updated_at` bigint DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  KEY `idx_worker_group_updated_at` (`updated_at`),
  KEY `idx_worker_group_created_at` (`created_at`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4  COMMENT='节点分组';

SET FOREIGN_KEY_CHECKS = 1;
