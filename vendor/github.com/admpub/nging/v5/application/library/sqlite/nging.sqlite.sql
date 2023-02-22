CREATE TABLE `nging_alert_recipient` (
  `id` integer PRIMARY KEY ,
  `name` varchar(200) COLLATE NOCASE NOT NULL DEFAULT '',
  `account` varchar(255) COLLATE NOCASE NOT NULL DEFAULT '',
  `extra` text COLLATE NOCASE NOT NULL,
  `type` char(7) COLLATE NOCASE NOT NULL DEFAULT 'email',
  `platform` varchar(30) COLLATE NOCASE NOT NULL DEFAULT '',
  `description` varchar(500) COLLATE NOCASE NOT NULL DEFAULT '',
  `disabled` char(1) COLLATE NOCASE NOT NULL DEFAULT 'N',
  `created` int NOT NULL DEFAULT '0',
  `updated` int NOT NULL DEFAULT '0'
);
CREATE TABLE `nging_alert_topic` (
  `id` integer PRIMARY KEY ,
  `topic` varchar(100) COLLATE NOCASE NOT NULL DEFAULT '',
  `recipient_id` int NOT NULL DEFAULT '0',
  `disabled` char(1) COLLATE NOCASE NOT NULL DEFAULT 'N',
  `created` int NOT NULL DEFAULT '0',
  `updated` int NOT NULL DEFAULT '0'
);
CREATE INDEX IF NOT EXISTS `IDX_nging_alert_topic_alert_topic_recipient_id` ON `nging_alert_topic`(`recipient_id`);
CREATE INDEX IF NOT EXISTS `IDX_nging_alert_topic_alert_topic_topic_disabled` ON `nging_alert_topic`(`topic`,`disabled`);
CREATE TABLE `nging_cloud_backup` (
  `id` integer PRIMARY KEY ,
  `name` varchar(100) COLLATE NOCASE NOT NULL DEFAULT '',
  `source_path` varchar(200) COLLATE NOCASE NOT NULL,
  `ignore_rule` varchar(200) COLLATE NOCASE NOT NULL DEFAULT '',
  `wait_fill_completed` char(1) COLLATE NOCASE NOT NULL DEFAULT 'N',
  `ignore_wait_rule` varchar(200) COLLATE NOCASE NOT NULL DEFAULT '',
  `delay` int NOT NULL DEFAULT '0',
  `dest_storage` int NOT NULL,
  `dest_path` varchar(200) COLLATE NOCASE NOT NULL,
  `result` varchar(255) COLLATE NOCASE NOT NULL,
  `last_executed` int NOT NULL DEFAULT '0',
  `status` char(7) COLLATE NOCASE NOT NULL DEFAULT 'idle',
  `disabled` char(1) COLLATE NOCASE NOT NULL DEFAULT 'N',
  `created` int NOT NULL DEFAULT '0',
  `updated` int NOT NULL DEFAULT '0'
);
CREATE TABLE `nging_cloud_storage` (
  `id` integer PRIMARY KEY ,
  `name` varchar(200) COLLATE NOCASE NOT NULL DEFAULT '',
  `type` varchar(30) COLLATE NOCASE NOT NULL DEFAULT 'aws',
  `key` varchar(128) COLLATE NOCASE NOT NULL DEFAULT '',
  `secret` varchar(200) COLLATE NOCASE NOT NULL DEFAULT '',
  `bucket` varchar(200) COLLATE NOCASE NOT NULL DEFAULT '',
  `endpoint` varchar(200) COLLATE NOCASE NOT NULL DEFAULT '',
  `region` varchar(100) COLLATE NOCASE NOT NULL DEFAULT '',
  `secure` char(1) COLLATE NOCASE NOT NULL DEFAULT 'Y',
  `baseurl` varchar(200) COLLATE NOCASE NOT NULL DEFAULT '',
  `created` int NOT NULL DEFAULT '0',
  `updated` int NOT NULL DEFAULT '0'
);
CREATE TABLE `nging_code_invitation` (
  `id` integer PRIMARY KEY ,
  `uid` int NOT NULL DEFAULT '0',
  `recv_uid` int NOT NULL DEFAULT '0',
  `code` varchar(40) COLLATE NOCASE NOT NULL,
  `created` int NOT NULL,
  `used` int NOT NULL DEFAULT '0',
  `start` int NOT NULL DEFAULT '0',
  `end` int NOT NULL DEFAULT '0',
  `disabled` char(1) COLLATE NOCASE NOT NULL DEFAULT 'N',
  `role_ids` text COLLATE NOCASE NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS `UNQ_nging_code_invitation_code_invitation_code` ON `nging_code_invitation`(`code`);
CREATE TABLE `nging_code_verification` (
  `id` integer PRIMARY KEY ,
  `code` varchar(40) COLLATE NOCASE NOT NULL,
  `created` int NOT NULL,
  `owner_id` bigint NOT NULL DEFAULT '0',
  `owner_type` char(8) COLLATE NOCASE NOT NULL DEFAULT 'user',
  `used` int NOT NULL DEFAULT '0',
  `purpose` varchar(40) COLLATE NOCASE NOT NULL,
  `start` int NOT NULL DEFAULT '0',
  `end` int NOT NULL DEFAULT '0',
  `disabled` char(1) COLLATE NOCASE NOT NULL DEFAULT 'N',
  `send_method` varchar(60) COLLATE NOCASE NOT NULL DEFAULT 'mobile',
  `send_to` varchar(120) COLLATE NOCASE NOT NULL DEFAULT ''
);
CREATE TABLE `nging_config` (
  `key` varchar(60) COLLATE NOCASE NOT NULL,
  `group` varchar(60) COLLATE NOCASE NOT NULL DEFAULT '',
  `label` varchar(90) COLLATE NOCASE NOT NULL DEFAULT '',
  `value` varchar(5000) COLLATE NOCASE NOT NULL DEFAULT '',
  `description` varchar(255) COLLATE NOCASE NOT NULL DEFAULT '',
  `type` char(5) COLLATE NOCASE NOT NULL DEFAULT 'text',
  `sort` int NOT NULL DEFAULT '0',
  `disabled` char(1) COLLATE NOCASE NOT NULL DEFAULT 'N',
  `encrypted` char(1) COLLATE NOCASE NOT NULL DEFAULT 'N'
);
CREATE INDEX IF NOT EXISTS `IDX_nging_config_config_group` ON `nging_config`(`group`);
CREATE TABLE `nging_file` (
  `id` integer PRIMARY KEY ,
  `owner_type` char(8) COLLATE NOCASE NOT NULL DEFAULT 'user',
  `owner_id` bigint NOT NULL DEFAULT '0',
  `name` varchar(150) COLLATE NOCASE NOT NULL DEFAULT '',
  `save_name` varchar(100) COLLATE NOCASE NOT NULL DEFAULT '',
  `save_path` varchar(200) COLLATE NOCASE NOT NULL DEFAULT '',
  `view_url` varchar(200) COLLATE NOCASE NOT NULL DEFAULT '',
  `ext` varchar(5) COLLATE NOCASE NOT NULL DEFAULT '',
  `mime` varchar(40) COLLATE NOCASE NOT NULL DEFAULT '',
  `type` varchar(30) COLLATE NOCASE NOT NULL DEFAULT 'image',
  `size` bigint NOT NULL DEFAULT '0',
  `width` int NOT NULL DEFAULT '0',
  `height` int NOT NULL DEFAULT '0',
  `dpi` int NOT NULL DEFAULT '0',
  `md5` char(32) COLLATE NOCASE NOT NULL DEFAULT '',
  `storer_name` varchar(30) COLLATE NOCASE NOT NULL DEFAULT '',
  `storer_id` varchar(32) COLLATE NOCASE NOT NULL DEFAULT '',
  `created` int NOT NULL DEFAULT '0',
  `updated` int NOT NULL DEFAULT '0',
  `sort` bigint NOT NULL DEFAULT '0',
  `status` integer NOT NULL DEFAULT '0',
  `category_id` int NOT NULL DEFAULT '0',
  `tags` varchar(200) COLLATE NOCASE NOT NULL DEFAULT '',
  `subdir` varchar(60) COLLATE NOCASE NOT NULL DEFAULT '',
  `used_times` int NOT NULL DEFAULT '0'
);
CREATE INDEX IF NOT EXISTS `IDX_nging_file_file_category_id` ON `nging_file`(`category_id`);
CREATE INDEX IF NOT EXISTS `IDX_nging_file_file_view_url` ON `nging_file`(`view_url`);
CREATE INDEX IF NOT EXISTS `IDX_nging_file_file_owner_id_and_type` ON `nging_file`(`owner_id`,`owner_type`);
CREATE INDEX IF NOT EXISTS `IDX_nging_file_file_subdir` ON `nging_file`(`subdir`);
CREATE TABLE `nging_file_embedded` (
  `id` integer PRIMARY KEY ,
  `project` varchar(50) COLLATE NOCASE NOT NULL,
  `table_id` varchar(60) COLLATE NOCASE NOT NULL DEFAULT '0',
  `table_name` varchar(60) COLLATE NOCASE NOT NULL,
  `field_name` varchar(50) COLLATE NOCASE NOT NULL,
  `file_ids` varchar(1000) COLLATE NOCASE NOT NULL,
  `embedded` char(1) COLLATE NOCASE NOT NULL DEFAULT 'Y'
);
CREATE UNIQUE INDEX IF NOT EXISTS `UNQ_nging_file_embedded_file_embedded_table_id_field_table` ON `nging_file_embedded`(`table_id`,`field_name`,`table_name`);
CREATE TABLE `nging_file_moved` (
  `id` integer PRIMARY KEY ,
  `file_id` bigint NOT NULL,
  `from` varchar(200) COLLATE NOCASE NOT NULL,
  `to` varchar(200) COLLATE NOCASE NOT NULL,
  `thumb_id` bigint NOT NULL DEFAULT '0',
  `created` int NOT NULL DEFAULT '0'
);
CREATE UNIQUE INDEX IF NOT EXISTS `UNQ_nging_file_moved_file_moved_from` ON `nging_file_moved`(`from`);
CREATE TABLE `nging_file_thumb` (
  `id` integer PRIMARY KEY ,
  `file_id` bigint NOT NULL,
  `size` bigint NOT NULL,
  `width` int NOT NULL,
  `height` int NOT NULL,
  `dpi` int NOT NULL DEFAULT '0',
  `save_name` varchar(100) COLLATE NOCASE NOT NULL,
  `save_path` varchar(200) COLLATE NOCASE NOT NULL,
  `view_url` varchar(200) COLLATE NOCASE NOT NULL,
  `used_times` int NOT NULL DEFAULT '0',
  `md5` char(32) COLLATE NOCASE NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS `IDX_nging_file_thumb_file_thumb_view_url` ON `nging_file_thumb`(`view_url`);
CREATE UNIQUE INDEX IF NOT EXISTS `UNQ_nging_file_thumb_file_thumb_save_path` ON `nging_file_thumb`(`save_path`);
CREATE UNIQUE INDEX IF NOT EXISTS `UNQ_nging_file_thumb_file_thumb_file_id_size_flag` ON `nging_file_thumb`(`file_id`,`size`);
CREATE TABLE `nging_kv` (
  `id` integer PRIMARY KEY ,
  `key` varchar(100) COLLATE NOCASE NOT NULL,
  `value` varchar(255) COLLATE NOCASE NOT NULL,
  `description` varchar(100) COLLATE NOCASE NOT NULL DEFAULT '',
  `type` varchar(50) COLLATE NOCASE NOT NULL,
  `sort` int NOT NULL DEFAULT '0',
  `updated` int NOT NULL DEFAULT '0',
  `child_key_type` varchar(10) COLLATE NOCASE NOT NULL DEFAULT 'text'
);
CREATE UNIQUE INDEX IF NOT EXISTS `UNQ_nging_kv_kv_key_type` ON `nging_kv`(`key`,`type`);
CREATE TABLE `nging_login_log` (
  `owner_type` char(8) COLLATE NOCASE NOT NULL DEFAULT 'user',
  `owner_id` bigint NOT NULL DEFAULT '0',
  `session_id` varchar(64) COLLATE NOCASE NOT NULL DEFAULT '',
  `username` varchar(60) COLLATE NOCASE NOT NULL DEFAULT '',
  `errpwd` varchar(100) COLLATE NOCASE NOT NULL DEFAULT '',
  `ip_address` varchar(46) COLLATE NOCASE NOT NULL DEFAULT '',
  `ip_location` varchar(200) COLLATE NOCASE NOT NULL DEFAULT '',
  `user_agent` varchar(255) COLLATE NOCASE NOT NULL DEFAULT '',
  `success` char(1) COLLATE NOCASE NOT NULL DEFAULT 'N',
  `failmsg` varchar(100) COLLATE NOCASE NOT NULL DEFAULT '',
  `day` int NOT NULL DEFAULT '0',
  `created` int NOT NULL DEFAULT '0'
);
CREATE INDEX IF NOT EXISTS `IDX_nging_login_log_login_log_ip_address` ON `nging_login_log`(`ip_address`,`day`);
CREATE INDEX IF NOT EXISTS `IDX_nging_login_log_login_log_owner` ON `nging_login_log`(`owner_type`,`owner_id`,`session_id`);
CREATE INDEX IF NOT EXISTS `IDX_nging_login_log_login_log_created` ON `nging_login_log`(`created` DESC);
CREATE INDEX IF NOT EXISTS `IDX_nging_login_log_login_log_success` ON `nging_login_log`(`success`);
CREATE TABLE `nging_sending_log` (
  `id` integer PRIMARY KEY ,
  `created` int NOT NULL,
  `sent_at` int NOT NULL,
  `source_id` bigint NOT NULL DEFAULT '0',
  `source_type` varchar(30) COLLATE NOCASE NOT NULL DEFAULT 'user',
  `disabled` char(1) COLLATE NOCASE NOT NULL DEFAULT 'N',
  `method` varchar(60) COLLATE NOCASE NOT NULL DEFAULT 'mobile',
  `to` varchar(120) COLLATE NOCASE NOT NULL DEFAULT '',
  `provider` varchar(60) COLLATE NOCASE NOT NULL DEFAULT '',
  `result` varchar(255) COLLATE NOCASE NOT NULL DEFAULT '',
  `status` char(7) COLLATE NOCASE NOT NULL DEFAULT 'waiting',
  `retries` int NOT NULL DEFAULT '0',
  `content` varchar(255) COLLATE NOCASE NOT NULL DEFAULT '',
  `params` varchar(255) COLLATE NOCASE NOT NULL DEFAULT '',
  `appointment_time` int NOT NULL DEFAULT '0'
);
CREATE TABLE `nging_task` (
  `id` integer PRIMARY KEY ,
  `uid` int NOT NULL DEFAULT '0',
  `group_id` int NOT NULL DEFAULT '0',
  `name` varchar(50) COLLATE NOCASE NOT NULL DEFAULT '',
  `type` tinyint NOT NULL DEFAULT '0',
  `description` varchar(200) COLLATE NOCASE NOT NULL DEFAULT '',
  `cron_spec` varchar(100) COLLATE NOCASE NOT NULL DEFAULT '',
  `concurrent` tinyint NOT NULL DEFAULT '0',
  `command` text COLLATE NOCASE NOT NULL,
  `work_directory` varchar(255) COLLATE NOCASE NOT NULL DEFAULT '',
  `env` text COLLATE NOCASE NOT NULL,
  `disabled` char(1) COLLATE NOCASE NOT NULL DEFAULT 'N',
  `enable_notify` tinyint NOT NULL DEFAULT '0',
  `notify_email` text COLLATE NOCASE NOT NULL,
  `timeout` bigint NOT NULL DEFAULT '0',
  `execute_times` int NOT NULL DEFAULT '0',
  `prev_time` int NOT NULL DEFAULT '0',
  `created` int NOT NULL DEFAULT '0',
  `updated` int NOT NULL DEFAULT '0',
  `closed_log` char(1) COLLATE NOCASE NOT NULL DEFAULT 'N'
);
CREATE INDEX IF NOT EXISTS `IDX_nging_task_task_uid` ON `nging_task`(`uid`);
CREATE INDEX IF NOT EXISTS `IDX_nging_task_task_group_id` ON `nging_task`(`group_id`);
CREATE TABLE `nging_task_group` (
  `id` integer PRIMARY KEY ,
  `uid` int NOT NULL DEFAULT '0',
  `name` varchar(60) COLLATE NOCASE NOT NULL,
  `description` varchar(255) COLLATE NOCASE NOT NULL DEFAULT '',
  `created` int NOT NULL DEFAULT '0',
  `updated` int NOT NULL DEFAULT '0',
  `cmd_prefix` varchar(255) COLLATE NOCASE NOT NULL DEFAULT '',
  `cmd_suffix` varchar(255) COLLATE NOCASE NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS `IDX_nging_task_group_task_group_uid` ON `nging_task_group`(`uid`);
CREATE TABLE `nging_task_log` (
  `id` integer PRIMARY KEY ,
  `task_id` int NOT NULL DEFAULT '0',
  `output` mediumtext COLLATE NOCASE NOT NULL,
  `error` text COLLATE NOCASE NOT NULL,
  `status` char(7) COLLATE NOCASE NOT NULL DEFAULT 'success',
  `elapsed` int NOT NULL DEFAULT '0',
  `created` int NOT NULL DEFAULT '0'
);
CREATE INDEX IF NOT EXISTS `IDX_nging_task_log_task_log_task_id_created` ON `nging_task_log`(`task_id`,`created`);
CREATE TABLE `nging_user` (
  `id` integer PRIMARY KEY ,
  `username` varchar(30) COLLATE NOCASE NOT NULL DEFAULT '',
  `email` varchar(50) COLLATE NOCASE NOT NULL DEFAULT '',
  `mobile` varchar(15) COLLATE NOCASE NOT NULL DEFAULT '',
  `password` char(64) COLLATE NOCASE NOT NULL DEFAULT '',
  `salt` char(64) COLLATE NOCASE NOT NULL DEFAULT '',
  `safe_pwd` char(64) COLLATE NOCASE NOT NULL DEFAULT '',
  `session_id` char(64) COLLATE NOCASE NOT NULL DEFAULT '',
  `avatar` varchar(200) COLLATE NOCASE NOT NULL DEFAULT '',
  `gender` char(6) COLLATE NOCASE NOT NULL DEFAULT 'secret',
  `last_login` int NOT NULL DEFAULT '0',
  `last_ip` varchar(150) COLLATE NOCASE NOT NULL DEFAULT '',
  `login_fails` int NOT NULL DEFAULT '0',
  `disabled` char(1) COLLATE NOCASE NOT NULL DEFAULT 'N',
  `online` char(1) COLLATE NOCASE NOT NULL DEFAULT 'N',
  `role_ids` text COLLATE NOCASE NOT NULL,
  `created` int NOT NULL DEFAULT '0',
  `updated` int NOT NULL DEFAULT '0',
  `file_size` bigint NOT NULL DEFAULT '0',
  `file_num` bigint NOT NULL DEFAULT '0'
);
CREATE UNIQUE INDEX IF NOT EXISTS `UNQ_nging_user_user_username` ON `nging_user`(`username`);
CREATE TABLE `nging_user_role` (
  `id` integer PRIMARY KEY ,
  `name` varchar(60) COLLATE NOCASE NOT NULL,
  `description` tinytext COLLATE NOCASE NOT NULL,
  `created` int NOT NULL,
  `updated` int NOT NULL DEFAULT '0',
  `disabled` char(1) COLLATE NOCASE NOT NULL DEFAULT 'N',
  `parent_id` int NOT NULL DEFAULT '0'
);
CREATE TABLE `nging_user_role_permission` (
  `role_id` int NOT NULL,
  `type` varchar(30) COLLATE NOCASE NOT NULL,
  `permission` text COLLATE NOCASE NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS `UNQ_nging_user_role_permission_user_role_permission_uniqid` ON `nging_user_role_permission`(`role_id`,`type`);
CREATE TABLE `nging_user_u2f` (
  `id` integer PRIMARY KEY ,
  `uid` int NOT NULL,
  `name` varchar(100) COLLATE NOCASE NOT NULL DEFAULT '',
  `token` varchar(255) COLLATE NOCASE NOT NULL,
  `type` varchar(30) COLLATE NOCASE NOT NULL,
  `extra` text COLLATE NOCASE NOT NULL,
  `step` tinyint NOT NULL DEFAULT '2',
  `created` int NOT NULL DEFAULT '0'
);
CREATE INDEX IF NOT EXISTS `IDX_nging_user_u2f_user_u2f_step` ON `nging_user_u2f`(`step`);
CREATE UNIQUE INDEX IF NOT EXISTS `UNQ_nging_user_u2f_user_u2f_uid_type` ON `nging_user_u2f`(`uid`,`type`);
