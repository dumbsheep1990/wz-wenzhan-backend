-- 万知文站数据库表结构

-- 用户表
CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `email` varchar(100) NOT NULL COMMENT '邮箱',
  `password` varchar(255) NOT NULL COMMENT '密码',
  `avatar` varchar(255) DEFAULT '' COMMENT '头像',
  `nickname` varchar(50) DEFAULT '' COMMENT '昵称',
  `status` tinyint DEFAULT '1' COMMENT '状态：1-激活，0-禁用',
  `last_login` datetime DEFAULT NULL COMMENT '最后登录时间',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`),
  UNIQUE KEY `idx_email` (`email`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 文件夹表
CREATE TABLE IF NOT EXISTS `folders` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '文件夹名称',
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `parent_id` bigint unsigned DEFAULT NULL COMMENT '父文件夹ID',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_deleted_at` (`deleted_at`),
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件夹表';

-- 文档表
CREATE TABLE IF NOT EXISTS `documents` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL COMMENT '文档标题',
  `content` longtext COMMENT '文档内容',
  `type` varchar(20) NOT NULL COMMENT '文档类型：word,excel,mindmap,note,ai_draft,imported',
  `status` tinyint DEFAULT '1' COMMENT '状态：1-草稿，2-已发布，3-已归档',
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `folder_id` bigint unsigned DEFAULT NULL COMMENT '文件夹ID',
  `tags` varchar(500) DEFAULT '' COMMENT '标签，逗号分隔',
  `size` bigint DEFAULT '0' COMMENT '文件大小（字节）',
  `view_count` int DEFAULT '0' COMMENT '查看次数',
  `is_shared` tinyint(1) DEFAULT '0' COMMENT '是否分享',
  `share_token` varchar(32) DEFAULT '' COMMENT '分享令牌',
  `share_expiry` datetime DEFAULT NULL COMMENT '分享过期时间',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_folder_id` (`folder_id`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`),
  KEY `idx_share_token` (`share_token`),
  KEY `idx_deleted_at` (`deleted_at`),
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  FOREIGN KEY (`folder_id`) REFERENCES `folders` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文档表';

-- 活动记录表
CREATE TABLE IF NOT EXISTS `activities` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `type` varchar(20) NOT NULL COMMENT '活动类型：create,update,delete,view,share,copy,move',
  `resource_type` varchar(20) NOT NULL COMMENT '资源类型：document,folder',
  `resource_id` bigint unsigned NOT NULL COMMENT '资源ID',
  `resource_name` varchar(255) NOT NULL COMMENT '资源名称',
  `description` varchar(500) DEFAULT '' COMMENT '描述',
  `ip_address` varchar(45) DEFAULT '' COMMENT 'IP地址',
  `user_agent` varchar(500) DEFAULT '' COMMENT '用户代理',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_type` (`type`),
  KEY `idx_resource_type` (`resource_type`),
  KEY `idx_resource_id` (`resource_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_deleted_at` (`deleted_at`),
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='活动记录表';

-- 回收站表
CREATE TABLE IF NOT EXISTS `recycle_items` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `resource_type` varchar(20) NOT NULL COMMENT '资源类型：document,folder',
  `resource_id` bigint unsigned NOT NULL COMMENT '资源ID',
  `resource_name` varchar(255) NOT NULL COMMENT '资源名称',
  `original_path` varchar(500) DEFAULT '' COMMENT '原始路径',
  `delete_reason` varchar(255) DEFAULT '' COMMENT '删除原因',
  `auto_delete` datetime DEFAULT NULL COMMENT '自动删除时间',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_resource_type` (`resource_type`),
  KEY `idx_resource_id` (`resource_id`),
  KEY `idx_auto_delete` (`auto_delete`),
  KEY `idx_deleted_at` (`deleted_at`),
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='回收站表';

-- 导航分类表
CREATE TABLE IF NOT EXISTS `categories` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL COMMENT '分类名称',
  `code` varchar(50) NOT NULL COMMENT '分类唯一代码',
  `parent_id` bigint unsigned DEFAULT NULL COMMENT '父分类ID',
  `level` int NOT NULL DEFAULT '1' COMMENT '分类层级',
  `icon` varchar(255) DEFAULT '' COMMENT '图标',
  `sort_order` int NOT NULL DEFAULT '0' COMMENT '排序',
  `is_visible` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否可见',
  `description` varchar(255) DEFAULT '' COMMENT '描述',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`code`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_sort_order` (`sort_order`),
  KEY `idx_is_visible` (`is_visible`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='导航分类表';

-- 文档分类关联表
CREATE TABLE IF NOT EXISTS `document_categories` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `document_id` bigint unsigned NOT NULL COMMENT '文档ID',
  `category_id` bigint unsigned NOT NULL COMMENT '分类ID',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_doc_cat` (`document_id`, `category_id`),
  KEY `idx_document_id` (`document_id`),
  KEY `idx_category_id` (`category_id`),
  FOREIGN KEY (`document_id`) REFERENCES `documents` (`id`) ON DELETE CASCADE,
  FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文档分类关联表';

-- 文件表
CREATE TABLE IF NOT EXISTS `files` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `filename` varchar(255) NOT NULL COMMENT '文件名',
  `path` varchar(500) NOT NULL COMMENT '存储路径',
  `size` bigint NOT NULL COMMENT '文件大小（字节）',
  `mime_type` varchar(100) NOT NULL COMMENT '文件类型',
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `folder_id` bigint unsigned DEFAULT NULL COMMENT '文件夹ID',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_folder_id` (`folder_id`),
  KEY `idx_deleted_at` (`deleted_at`),
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  FOREIGN KEY (`folder_id`) REFERENCES `folders` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件表';

-- 插入测试数据（可选）
INSERT INTO `users` (`username`, `email`, `password`, `nickname`, `status`) VALUES
('admin', 'admin@wenzhan.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '管理员', 1),
('demo', 'demo@wenzhan.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '演示用户', 1);
