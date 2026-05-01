CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户 ID',
  `email` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '邮箱地址（统一存储小写字母）',
  `password_hash` char(60) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码哈希',
  `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名称',
  `last_login_at` datetime(3) DEFAULT NULL COMMENT '最近登录时间',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_users_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表'

CREATE TABLE `refresh_tokens` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '记录 ID',
  `token_id` char(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Refresh Token ID',
  `user_id` bigint unsigned NOT NULL COMMENT '用户 ID',
  `token_hash` char(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Refresh Token 哈希',
  `expires_at` datetime(3) NOT NULL COMMENT '过期时间',
  `revoked_at` datetime(3) DEFAULT NULL COMMENT '吊销时间',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `replaced_by_token_id` char(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '新 Refresh Token ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_refresh_tokens_token_id` (`token_id`),
  UNIQUE KEY `uk_refresh_tokens_token_hash` (`token_hash`),
  KEY `idx_refresh_tokens_user_id` (`user_id`),
  CONSTRAINT `fk_refresh_tokens_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Refresh Token 表'

CREATE TABLE `email_send_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '记录 ID',
  `receiver_email` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '收件邮箱地址',
  `sender_email` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '发件邮箱地址',
  `scene` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '邮箱验证码业务场景：register|login|reset',
  `email_code` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '邮箱验证码',
  `status` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '邮件发送状态：success|failed',
  `error_message` varchar(1024) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '邮件发送失败原因',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='邮件发送记录表'



CREATE TABLE IF NOT EXISTS projects (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '项目 ID',
  user_id BIGINT UNSIGNED NOT NULL COMMENT '用户 ID',
  name VARCHAR(100) NOT NULL COMMENT '项目名称',
  building_name VARCHAR(100) NOT NULL DEFAULT '' COMMENT '建筑名称',
  building_location VARCHAR(255) NOT NULL DEFAULT '' COMMENT '建筑地址',
  building_description TEXT NULL COMMENT '建筑描述',
  progress_step TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '项目进度：0-初始化|1-基本信息|2-图像上传|3-模型检测|4-人工复核|5-报告生成',
  status VARCHAR(32) NOT NULL DEFAULT 'active' COMMENT '项目状态：active|archived',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (id),
  KEY idx_projects_user_id (user_id),
  CONSTRAINT fk_projects_user_id
    FOREIGN KEY (user_id) REFERENCES users (id)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='项目表';

CREATE TABLE IF NOT EXISTS project_categories (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '分类 ID',
  project_id BIGINT UNSIGNED NOT NULL COMMENT '所属项目 ID',
  name VARCHAR(50) NOT NULL COMMENT '分类名称，如东立面、西立面',
  is_default TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否默认分类',
  sort_order INT NOT NULL DEFAULT 1000 COMMENT '展示排序，数值越小越靠前',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (id),
  KEY idx_project_categories_project_id_sort (project_id, sort_order),
  UNIQUE KEY uk_project_categories_project_name (project_id, name),
  CONSTRAINT fk_project_categories_project_id
    FOREIGN KEY (project_id) REFERENCES projects (id)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='项目图像分类表';

CREATE TABLE IF NOT EXISTS project_images (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '图像 ID',
  project_id BIGINT UNSIGNED NOT NULL COMMENT '所属项目 ID',
  category_id BIGINT UNSIGNED NOT NULL COMMENT '所属分类 ID',
  image_uuid CHAR(36) NOT NULL COMMENT '图像 UUID，用于生成 OSS Key',
  original_key VARCHAR(1024) NOT NULL COMMENT '原图 OSS Key',
  thumbnail_key VARCHAR(1024) NOT NULL COMMENT '缩略图 OSS Key',
  file_name VARCHAR(255) NOT NULL DEFAULT '' COMMENT '原始文件名',
  content_type VARCHAR(100) NOT NULL DEFAULT '' COMMENT '原图 MIME 类型',
  size_bytes BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '原图文件大小',
  width INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '图像宽度',
  height INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '图像高度',
  sort_order INT NOT NULL DEFAULT 1000 COMMENT '分类内展示排序，数值越小越靠前',
  upload_status VARCHAR(32) NOT NULL DEFAULT 'pending' COMMENT '上传状态：pending|uploaded|failed',
  upload_expires_at DATETIME(3) NULL COMMENT '上传链接过期时间',
  uploaded_at DATETIME(3) NULL COMMENT '上传完成时间',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_project_images_image_uuid (image_uuid),
  KEY idx_project_images_project_category_sort (project_id, category_id, sort_order),
  KEY idx_project_images_project_upload_status (project_id, upload_status),
  CONSTRAINT fk_project_images_project_id
    FOREIGN KEY (project_id) REFERENCES projects (id)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT fk_project_images_category_id
    FOREIGN KEY (category_id) REFERENCES project_categories (id)
    ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='项目图像表';
