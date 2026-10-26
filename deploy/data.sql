
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
    `status` tinyint(1) DEFAULT NULL COMMENT '状态',
    `version` varchar(64) NOT NULL COMMENT '版本',
    `description` varchar(256) DEFAULT NULL COMMENT '描述',
    `start` bigint DEFAULT NULL COMMENT '开始时间',
    `end` bigint DEFAULT NULL COMMENT '结束时间',
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


DROP TABLE IF EXISTS `task_group`;
CREATE TABLE IF NOT EXISTS `task_group` (
                                            `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `name` varchar(32) NOT NULL COMMENT '名称',
    `description` varchar(256) DEFAULT NULL COMMENT '描述',
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
    `ak` varchar(64) NOT NULL DEFAULT '' COMMENT '签名AK',
    `sk` varchar(64) NOT NULL DEFAULT '' COMMENT '签名SK',
    `group_id` int(10) UNSIGNED DEFAULT NULL COMMENT '分组',
    `status` tinyint(1) DEFAULT NULL COMMENT '状态',
    `description` varchar(256) DEFAULT NULL COMMENT '描述',
    `created_at` bigint(20) DEFAULT NULL COMMENT '创建时间',
    `updated_at` bigint(20) DEFAULT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `ip` (`ip`),
    KEY `idx_worker_updated_at` (`updated_at`),
    KEY `idx_worker_created_at` (`created_at`),
    KEY `group_id` (`group_id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT='节点信息';


INSERT INTO `worker` (`id`, `name`, `ip`, `ak`, `sk`, `group_id`, `status`, `description`, `created_at`, `updated_at`) VALUES
    (2, '测试节点', '127.0.0.1:8121', 'client', '12345678', 13, 1, '测试节点12', 1713409635, 1719720853);


DROP TABLE IF EXISTS `worker_group`;
CREATE TABLE IF NOT EXISTS `worker_group` (
                                              `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `name` varchar(32) NOT NULL COMMENT '名称',
    `description` varchar(256) DEFAULT NULL COMMENT '描述',
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


-- ------------------------------------------------
-- 执行意图表（P1 调度内核：触发与执行解耦的状态载体）
-- ------------------------------------------------
DROP TABLE IF EXISTS `execution`;
CREATE TABLE IF NOT EXISTS `execution` (
    `id`           varchar(64)    NOT NULL COMMENT '执行意图ID（uuid）',
    `task_id`      int(10) UNSIGNED NOT NULL COMMENT '任务ID',
    `trigger_at`   bigint         NOT NULL COMMENT '触发时刻（整周期对齐，Unix秒）',
    `trigger_type` varchar(16)    NOT NULL DEFAULT 'cron' COMMENT '触发方式：cron | manual | misfire',
    `worker_id`    int(10) UNSIGNED DEFAULT NULL COMMENT '实际执行节点（认领后填写）',
    `executor`     varchar(64)    DEFAULT NULL COMMENT '执行体标识（instance_id:goroutine）',
    `status`       varchar(16)    NOT NULL DEFAULT 'pending' COMMENT '状态：pending | running | orphan | success | fail | cancel',
    `heartbeat_at` bigint         DEFAULT NULL COMMENT '最近心跳（Unix秒，孤儿检测依据）',
    `started_at`   bigint         DEFAULT NULL COMMENT '开始执行时间（Unix秒）',
    `finished_at`  bigint         DEFAULT NULL COMMENT '结束时间（Unix秒）',
    `retry_count`  int(11)        NOT NULL DEFAULT 0 COMMENT '连接重试次数',
    `last_error`   text           COMMENT '最近一次错误（终态失败时填充）',
    `created_at`   bigint         DEFAULT NULL COMMENT '创建时间',
    `updated_at`   bigint         DEFAULT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_task_trigger` (`task_id`, `trigger_at`),
    KEY `idx_status_heartbeat` (`status`, `heartbeat_at`),
    KEY `idx_task_id` (`task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务执行意图';


-- ------------------------------------------------
-- P2 迁移：route_type 拆分 + misfire 策略列
-- 存量数据：exec_type = shell/http → route_type = single（原"非 group 即单机"语义）
-- ------------------------------------------------
ALTER TABLE `task`
    ADD COLUMN `route_type` varchar(16) NOT NULL DEFAULT 'single' COMMENT '路由方式：group | single' AFTER `worker_type`,
    ADD COLUMN `misfire` varchar(16) NOT NULL DEFAULT 'latest' COMMENT '错过补偿策略：skip | immediate | latest' AFTER `max_exec_time`;

UPDATE `task` SET `route_type` = 'single' WHERE `exec_type` IN ('shell', 'http');
UPDATE `task` SET `route_type` = 'group' WHERE `exec_type` = 'group';


-- ------------------------------------------------
-- App 模式迁移：应用隔离 + worker 授权
--   - task/task_group/execution/log：加 app 列（所属应用，X-App-Key）
--   - worker/worker_group：全局（不加 app 列），新增 worker_app 授权表（哪些 app 可用）
-- 存量：app 填 'manager'（当前唯一调用方）；worker_app 不迁移（无授权记录 = 所有 app 可用）
-- ------------------------------------------------
ALTER TABLE `task`
    ADD COLUMN `app` varchar(32) NOT NULL DEFAULT '' COMMENT '所属应用（X-App-Key）' AFTER `id`,
    ADD KEY `idx_app` (`app`);

ALTER TABLE `task_group`
    ADD COLUMN `app` varchar(32) NOT NULL DEFAULT '' COMMENT '所属应用（X-App-Key）' AFTER `id`,
    ADD KEY `idx_app` (`app`);

ALTER TABLE `execution`
    ADD COLUMN `app` varchar(32) NOT NULL DEFAULT '' COMMENT '所属应用（X-App-Key）' AFTER `id`,
    ADD KEY `idx_app` (`app`);

ALTER TABLE `log`
    ADD COLUMN `app` varchar(32) NOT NULL DEFAULT '' COMMENT '所属应用（X-App-Key）' AFTER `id`,
    ADD KEY `idx_app` (`app`);

UPDATE `task` SET `app` = 'manager' WHERE `app` = '';
UPDATE `task_group` SET `app` = 'manager' WHERE `app` = '';

DROP TABLE IF EXISTS `worker_app`;
CREATE TABLE IF NOT EXISTS `worker_app` (
    `id`         int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `worker_id`  int(10) UNSIGNED NOT NULL COMMENT '节点ID',
    `app`        varchar(32) NOT NULL COMMENT '可用应用（X-App-Key）',
    `created_at` bigint DEFAULT NULL COMMENT '创建时间',
    `updated_at` bigint DEFAULT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_worker_app` (`worker_id`, `app`),
    KEY `idx_app` (`app`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='节点可用应用授权';


-- ------------------------------------------------
-- 节点实时探测：最后一次探测时间 / 最后有效时间 / 最后错误
-- （leader 每 30s 探测 status=true 的节点；在线判定 = last_heartbeat 在 90s 窗口内）
-- ------------------------------------------------
ALTER TABLE `worker`
    ADD COLUMN `last_check_at`  bigint DEFAULT NULL COMMENT '最后一次探测时间（含失败，Unix秒）' AFTER `status`,
    ADD COLUMN `last_heartbeat` bigint DEFAULT NULL COMMENT '最后一次有效探测时间（成功，Unix秒；NULL=从未上线）' AFTER `last_check_at`,
    ADD COLUMN `last_error`     text   DEFAULT NULL COMMENT '最后一次探测失败原因（成功时清空）' AFTER `last_heartbeat`;

-- ------------------------------------------------
-- 执行历史：execution 关联执行日志（log.uuid）
-- ------------------------------------------------
ALTER TABLE `execution`
    ADD COLUMN `log_uuid` varchar(64) DEFAULT NULL COMMENT '关联执行日志（log.uuid），执行时回填' AFTER `worker_id`;
