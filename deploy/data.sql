
SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


DROP TABLE IF EXISTS `log`;
CREATE TABLE IF NOT EXISTS `log` (
                                     `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `uuid` varchar(64) NOT NULL COMMENT '标识',
    `worker_id` int(10) UNSIGNED NOT NULL COMMENT '节点id',
    `worker_snapshot` varchar(2048) NOT NULL COMMENT '节点快照',
    `task_id` int(10) UNSIGNED NOT NULL COMMENT '任务id',
    `task_snapshot` varchar(2048) NOT NULL COMMENT '任务快照',
    `start_at` int(10) UNSIGNED NOT NULL COMMENT '开始时间',
    `end_at` int(10) UNSIGNED NOT NULL COMMENT '结束时间',
    `content` text NOT NULL COMMENT '执行日主',
    `status` char(32) NOT NULL COMMENT '执行结果',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uuid` (`uuid`)
    ) ENGINE=InnoDB AUTO_INCREMENT=439 DEFAULT CHARSET=utf8mb4 COMMENT='日志信息';


DROP TABLE IF EXISTS `task`;
CREATE TABLE IF NOT EXISTS `task` (
                                      `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `group_id` int(10) UNSIGNED DEFAULT NULL COMMENT '分组id',
    `name` varchar(64) NOT NULL COMMENT '名称',
    `tag` varchar(32) NOT NULL COMMENT '标签',
    `spec` varchar(32) NOT NULL COMMENT '表达式',
    `worker_type` char(32) NOT NULL COMMENT '选择类型',
    `worker_group_id` int(10) UNSIGNED DEFAULT NULL COMMENT '节点分组id',
    `worker_id` int(10) UNSIGNED DEFAULT NULL COMMENT '节点id',
    `exec_type` varchar(32) NOT NULL COMMENT '执行类型',
    `exec_value` text NOT NULL COMMENT '执行内容',
    `expect_code` int(11) NOT NULL COMMENT '预期状态码',
    `retry_count` int(11) NOT NULL COMMENT '重试次数',
    `retry_wait_time` int(11) NOT NULL COMMENT '等待时长',
    `max_exec_time` int(11) NOT NULL COMMENT '执行市场',
    `status` tinyint(1) NOT NULL COMMENT '状态',
    `version` varchar(64) NOT NULL COMMENT '版本',
    `description` varchar(256) NOT NULL COMMENT '描述',
    `created_at` bigint(20) DEFAULT NULL COMMENT '创建时间',
    `updated_at` bigint(20) DEFAULT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `version` (`version`),
    KEY `idx_task_updated_at` (`updated_at`),
    KEY `idx_task_created_at` (`created_at`),
    KEY `group_id` (`group_id`),
    KEY `worker_group_id` (`worker_group_id`),
    KEY `worker_id` (`worker_id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COMMENT='任务信息';

INSERT INTO `task` (`id`, `group_id`, `name`, `tag`, `spec`, `worker_type`, `worker_group_id`, `worker_id`, `exec_type`, `exec_value`, `expect_code`, `retry_count`, `retry_wait_time`, `max_exec_time`, `status`, `version`, `description`, `created_at`, `updated_at`) VALUES
    (3, 2, '测试任务', 'tets', '*/20 * * * * ?', 'worker', NULL, 2, 'http', '{\n    \"url\": \"https://www.baidu.com\",\n    \"params\": {\n        \"q\": 1\n    },\n    \"bodyJson\": {\n        \"params\": \"1\"\n    }\n}', 0, 1, 0, 11, 0, 'ABA12FD3D517908377520AA3C4345FC5', '测试任务', 1713409741, 1719733478);


DROP TABLE IF EXISTS `task_group`;
CREATE TABLE IF NOT EXISTS `task_group` (
                                            `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `name` varchar(32) NOT NULL COMMENT '名称',
    `description` varchar(256) NOT NULL COMMENT '描述',
    `created_at` bigint(20) DEFAULT NULL COMMENT '创建时间',
    `updated_at` bigint(20) DEFAULT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `name` (`name`),
    KEY `idx_task_group_updated_at` (`updated_at`),
    KEY `idx_task_group_created_at` (`created_at`)
    ) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT='任务分组';

INSERT INTO `task_group` (`id`, `name`, `description`, `created_at`, `updated_at`) VALUES
    (2, '测试任务分组', 'test', 1713409660, 1719679153);


DROP TABLE IF EXISTS `worker`;
CREATE TABLE IF NOT EXISTS `worker` (
                                        `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `name` varchar(32) NOT NULL COMMENT '名称',
    `ip` varchar(64) NOT NULL COMMENT 'ip',
    `group_id` int(10) UNSIGNED DEFAULT NULL COMMENT '分组',
    `status` tinyint(1) NOT NULL COMMENT '状态',
    `description` varchar(256) NOT NULL COMMENT '描述',
    `created_at` bigint(20) DEFAULT NULL COMMENT '创建时间',
    `updated_at` bigint(20) DEFAULT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `ip` (`ip`),
    KEY `idx_worker_updated_at` (`updated_at`),
    KEY `idx_worker_created_at` (`created_at`),
    KEY `group_id` (`group_id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT='节点信息';


INSERT INTO `worker` (`id`, `name`, `ip`, `group_id`, `status`, `description`, `created_at`, `updated_at`) VALUES
    (2, '测试节点', '127.0.0.1:8121', 13, 1, '测试节点12', 1713409635, 1719720853);


DROP TABLE IF EXISTS `worker_group`;
CREATE TABLE IF NOT EXISTS `worker_group` (
                                              `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `name` varchar(32) NOT NULL COMMENT '名称',
    `description` varchar(256) NOT NULL COMMENT '描述',
    `created_at` bigint(20) DEFAULT NULL COMMENT '创建时间',
    `updated_at` bigint(20) DEFAULT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `name` (`name`),
    KEY `idx_worker_group_updated_at` (`updated_at`),
    KEY `idx_worker_group_created_at` (`created_at`)
    ) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COMMENT='节点分组';


INSERT INTO `worker_group` (`id`, `name`, `description`, `created_at`, `updated_at`) VALUES
    (13, '测试节点分组', '测试分组', 1719651436, 1719679142);


ALTER TABLE `task`
    ADD CONSTRAINT `task_ibfk_1` FOREIGN KEY (`group_id`) REFERENCES `task_group` (`id`),
  ADD CONSTRAINT `task_ibfk_2` FOREIGN KEY (`worker_group_id`) REFERENCES `worker_group` (`id`),
  ADD CONSTRAINT `task_ibfk_3` FOREIGN KEY (`worker_id`) REFERENCES `worker` (`id`);


ALTER TABLE `worker`
    ADD CONSTRAINT `worker_ibfk_1` FOREIGN KEY (`group_id`) REFERENCES `worker_group` (`id`);
COMMIT;
