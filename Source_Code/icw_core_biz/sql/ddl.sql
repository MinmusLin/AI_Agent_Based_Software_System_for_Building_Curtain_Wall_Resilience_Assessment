-- users 用户表
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户 ID',
  `email` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '邮箱地址（统一存储小写字母）',
  `password_hash` char(60) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码哈希',
  `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名称',
  `last_login_at` datetime(3) DEFAULT NULL COMMENT '最近登录时间',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_users_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- refresh_tokens Refresh Token 表
CREATE TABLE `refresh_tokens` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '记录 ID',
  `token_id` char(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Refresh Token ID',
  `user_id` bigint unsigned NOT NULL COMMENT '用户 ID',
  `token_hash` char(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Refresh Token 哈希',
  `expires_at` datetime(3) NOT NULL COMMENT '过期时间',
  `revoked_at` datetime(3) DEFAULT NULL COMMENT '吊销时间',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `replaced_by_token_id` char(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '新 Refresh Token ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_refresh_tokens_token_id` (`token_id`),
  UNIQUE KEY `uk_refresh_tokens_token_hash` (`token_hash`),
  KEY `idx_refresh_tokens_expires_at` (`expires_at`),
  KEY `idx_refresh_tokens_user_id_revoked_at` (`user_id`,`revoked_at`),
  CONSTRAINT `fk_refresh_tokens_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Refresh Token 表';

-- email_send_logs 邮件发送记录表
CREATE TABLE `email_send_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '记录 ID',
  `receiver_email` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '收件邮箱地址',
  `sender_email` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '发件邮箱地址',
  `scene` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '邮箱验证码业务场景：register|login|reset',
  `email_code` char(6) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '邮箱验证码',
  `status` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '邮件发送状态：success|failed',
  `error_message` varchar(2048) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '邮件发送失败原因',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_email_send_logs_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='邮件发送记录表';

-- projects 项目表
CREATE TABLE `projects` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '项目 ID',
  `user_id` bigint unsigned NOT NULL COMMENT '用户 ID',
  `name` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '项目名称',
  `building_name` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '建筑名称',
  `building_location` varchar(1024) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '建筑地址',
  `built_year` smallint unsigned DEFAULT NULL COMMENT '建筑建成年份',
  `building_description` text COLLATE utf8mb4_unicode_ci COMMENT '建筑描述',
  `known_issues` text COLLATE utf8mb4_unicode_ci COMMENT '已知问题或人工先验描述',
  `assessment_goal` text COLLATE utf8mb4_unicode_ci COMMENT '评估目标或重点关注方向',
  `progress` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '项目进度',
  `status` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'active' COMMENT '项目状态：active|completed|deleted',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_projects_user_id_status_created_at` (`user_id`,`status`,`created_at`),
  CONSTRAINT `fk_projects_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='项目表';

-- project_groups 项目图像组表
CREATE TABLE `project_groups` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '图像组 ID',
  `project_id` bigint unsigned NOT NULL COMMENT '项目 ID',
  `user_id` bigint unsigned NOT NULL COMMENT '用户 ID',
  `name` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '图像组名称',
  `sort_order` decimal(20,10) NOT NULL DEFAULT '0.0000000000' COMMENT '顺序优先级',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_project_groups_project_id_user_id_name` (`project_id`,`user_id`,`name`),
  KEY `idx_project_groups_user_id_project_id_sort_order_created_at` (`user_id`,`project_id`,`sort_order`,`created_at`),
  CONSTRAINT `fk_project_groups_project_id` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_project_groups_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='项目图像组表';

-- project_group_images 项目图像表
CREATE TABLE `project_group_images` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '图像 ID',
  `group_id` bigint unsigned NOT NULL COMMENT '图像组 ID',
  `project_id` bigint unsigned NOT NULL COMMENT '项目 ID',
  `user_id` bigint unsigned NOT NULL COMMENT '用户 ID',
  `uuid` char(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '图像 UUID',
  `file_name` varchar(2048) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '原图文件名',
  `content_type` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '原图 MIME 类型',
  `size_bytes` bigint unsigned NOT NULL DEFAULT '0' COMMENT '原图文件大小',
  `width` int unsigned NOT NULL DEFAULT '0' COMMENT '图像宽度',
  `height` int unsigned NOT NULL DEFAULT '0' COMMENT '图像高度',
  `metadata` json NOT NULL DEFAULT (json_object()) COMMENT '图像元数据',
  `status` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending' COMMENT '图像状态：pending|uploaded|failed',
  `uploaded_at` datetime(3) DEFAULT NULL COMMENT '上传时间',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_project_group_images_uuid` (`uuid`),
  KEY `fk_project_group_images_group_id` (`group_id`),
  KEY `fk_project_group_images_project_id` (`project_id`),
  KEY `idx_project_group_images_user_id_project_id_created_at` (`user_id`,`project_id`,`created_at`),
  KEY `idx_project_group_images_user_id_project_id_group_id_created_at` (`user_id`,`project_id`,`group_id`,`created_at`),
  KEY `idx_project_group_images_status_created_at` (`status`,`created_at`),
  CONSTRAINT `fk_project_group_images_group_id` FOREIGN KEY (`group_id`) REFERENCES `project_groups` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `fk_project_group_images_project_id` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_project_group_images_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='项目图像表';

-- project_detection_tasks 项目图像检测主任务表
CREATE TABLE `project_detection_tasks` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主任务 ID',
  `uuid` char(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主任务 UUID',
  `user_id` bigint unsigned NOT NULL COMMENT '用户 ID',
  `project_id` bigint unsigned NOT NULL COMMENT '项目 ID',
  `image_id` bigint unsigned NOT NULL COMMENT '图像 ID',
  `image_uuid` char(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '图像 UUID',
  `status` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending' COMMENT '主任务状态：pending|classifying|detecting|summarizing|succeeded|failed',
  `corrosion_should_execute` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否执行金属锈蚀检测',
  `corrosion_task_id` bigint unsigned DEFAULT NULL COMMENT '金属锈蚀检测子任务 ID',
  `crack_should_execute` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否执行石材裂缝检测',
  `crack_task_id` bigint unsigned DEFAULT NULL COMMENT '石材裂缝检测子任务 ID',
  `stain_should_execute` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否执行石材污渍检测',
  `stain_task_id` bigint unsigned DEFAULT NULL COMMENT '石材污渍检测子任务 ID',
  `flatness_should_execute` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否执行玻璃平整度检测',
  `flatness_task_id` bigint unsigned DEFAULT NULL COMMENT '玻璃平整度检测子任务 ID',
  `spalling_should_execute` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否执行玻璃爆裂检测',
  `spalling_task_id` bigint unsigned DEFAULT NULL COMMENT '玻璃爆裂检测子任务 ID',
  `started_at` datetime(3) DEFAULT NULL COMMENT '开始时间',
  `finished_at` datetime(3) DEFAULT NULL COMMENT '完成时间',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_project_detection_tasks_uuid` (`uuid`),
  UNIQUE KEY `uk_project_detection_tasks_image_uuid` (`image_uuid`),
  CONSTRAINT `fk_project_detection_tasks_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_project_detection_tasks_project_id` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_project_detection_tasks_image_id` FOREIGN KEY (`image_id`) REFERENCES `project_group_images` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='项目图像检测主任务表';

-- project_detection_sub_tasks 项目图像检测子任务表
CREATE TABLE `project_detection_sub_tasks` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '子任务 ID',
  `uuid` char(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '子任务 UUID',
  `main_task_id` bigint unsigned NOT NULL COMMENT '主任务 ID',
  `main_task_uuid` char(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主任务 UUID',
  `user_id` bigint unsigned NOT NULL COMMENT '用户 ID',
  `project_id` bigint unsigned NOT NULL COMMENT '项目 ID',
  `image_id` bigint unsigned NOT NULL COMMENT '图像 ID',
  `image_uuid` char(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '图像 UUID',
  `status` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending' COMMENT '子任务状态：pending|running|succeeded|failed',
  `started_at` datetime(3) DEFAULT NULL COMMENT '开始时间',
  `finished_at` datetime(3) DEFAULT NULL COMMENT '完成时间',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_project_detection_sub_tasks_uuid` (`uuid`),
  UNIQUE KEY `uk_project_detection_sub_tasks_main_task_uuid` (`main_task_uuid`),
  UNIQUE KEY `uk_project_detection_sub_tasks_image_uuid` (`image_uuid`),
  CONSTRAINT `fk_project_detection_sub_tasks_main_task_id` FOREIGN KEY (`main_task_id`) REFERENCES `project_detection_tasks` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_project_detection_sub_tasks_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_project_detection_sub_tasks_project_id` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_project_detection_sub_tasks_image_id` FOREIGN KEY (`image_id`) REFERENCES `project_group_images` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='项目图像检测子任务表';
