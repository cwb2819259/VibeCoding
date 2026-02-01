-- C 端用户表（仅社区用户；登录方式：手机号+验证码。无外键）
CREATE TABLE IF NOT EXISTS `users` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `phone` VARCHAR(20) NOT NULL DEFAULT '' COMMENT '手机号，C 端登录',
  `nickname` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '昵称',
  `avatar_url` VARCHAR(512) NOT NULL DEFAULT '' COMMENT '头像 URL',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='C 端用户表';

-- B 端管理员表（仅后台登录与鉴权，与 users 隔离。无外键）
CREATE TABLE IF NOT EXISTS `admins` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `username` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '登录账号',
  `password_hash` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '密码哈希，如 bcrypt',
  `nickname` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '显示名',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='B 端管理员表';

-- 帖子表（user_id 逻辑关联 users.id，不在库中建外键）
CREATE TABLE IF NOT EXISTS `posts` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '发帖人，逻辑关联 users.id',
  `content` TEXT NOT NULL COMMENT '正文',
  `type` ENUM('text','image','video') NOT NULL DEFAULT 'text' COMMENT '帖子类型',
  `status` ENUM('normal','hidden','flagged') NOT NULL DEFAULT 'normal' COMMENT 'normal 可见，hidden 隐藏，flagged 标记',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status_created` (`status`,`created_at` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='帖子表';

-- 帖子媒体表（post_id 逻辑关联 posts.id；单帖最多 9 图或 1 视频。无外键）
CREATE TABLE IF NOT EXISTS `post_media` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `post_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '帖子 ID，逻辑关联 posts.id',
  `type` ENUM('image','video') NOT NULL DEFAULT 'image' COMMENT '媒体类型',
  `url` VARCHAR(512) NOT NULL DEFAULT '' COMMENT '本地存储相对路径或 URL',
  `sort_order` SMALLINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '排序，从 0 起',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_post_id` (`post_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='帖子媒体表';

-- 点赞表（user_id 逻辑关联 users.id，target 为帖子或评论。无外键）
CREATE TABLE IF NOT EXISTS `likes` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '点赞人，逻辑关联 users.id',
  `target_type` ENUM('post','comment') NOT NULL DEFAULT 'post' COMMENT '目标类型',
  `target_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '目标 ID',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_target` (`user_id`,`target_type`,`target_id`),
  KEY `idx_target` (`target_type`,`target_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='点赞表';

-- 评论表（post_id 逻辑关联 posts.id，user_id 逻辑关联 users.id。无外键）
CREATE TABLE IF NOT EXISTS `comments` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `post_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '帖子 ID，逻辑关联 posts.id',
  `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '评论人，逻辑关联 users.id',
  `content` TEXT NOT NULL COMMENT '评论内容',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_post_id` (`post_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论表';

-- 通知表（user_id 逻辑关联 users.id，接收人均为 C 端用户。无外键）
CREATE TABLE IF NOT EXISTS `notifications` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '接收人，逻辑关联 users.id',
  `type` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '类型，如 like、comment、system',
  `related_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '关联 ID，如帖子/评论 ID，0 表示无',
  `payload` JSON NOT NULL COMMENT '扩展数据，默认空对象',
  `read_at` DATETIME NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '已读时间，未读时为 1970-01-01 00:00:00',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_read_created` (`user_id`,`read_at`,`created_at` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='消息通知表';

-- 话题/标签表（无外键）
CREATE TABLE IF NOT EXISTS `topics` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '话题名',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='话题标签表';

-- 帖子-话题多对多（post_id 逻辑关联 posts.id，topic_id 逻辑关联 topics.id。无外键）
CREATE TABLE IF NOT EXISTS `post_topics` (
  `post_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '帖子 ID，逻辑关联 posts.id',
  `topic_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '话题 ID，逻辑关联 topics.id',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`post_id`,`topic_id`),
  KEY `idx_topic_id` (`topic_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='帖子-话题关联表';
