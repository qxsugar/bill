CREATE TABLE IF NOT EXISTS `bill_users` (
  `id`         BIGINT       NOT NULL AUTO_INCREMENT,
  `openid`     VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '微信 openid',
  `nickname`   VARCHAR(50)  NOT NULL DEFAULT '' COMMENT '昵称',
  `avatar`     VARCHAR(255) NOT NULL DEFAULT '' COMMENT '头像 URL',
  `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_openid` (`openid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

CREATE TABLE IF NOT EXISTS `bill_rooms` (
  `id`          BIGINT      NOT NULL AUTO_INCREMENT,
  `code`        VARCHAR(10) NOT NULL DEFAULT '' COMMENT '房间码',
  `owner_id`    BIGINT      NOT NULL DEFAULT 0 COMMENT '房主 user_id',
  `status`      TINYINT     NOT NULL DEFAULT 0 COMMENT '0=活跃 1=已结算',
  `settled_at`  DATETIME    NULL     DEFAULT NULL COMMENT '结算时间',
  `created_at`  DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_owner_id` (`owner_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='房间表';

CREATE TABLE IF NOT EXISTS `bill_room_members` (
  `id`        BIGINT         NOT NULL AUTO_INCREMENT,
  `room_id`   BIGINT         NOT NULL DEFAULT 0,
  `user_id`   BIGINT         NOT NULL DEFAULT 0,
  `balance`   DECIMAL(10,2)  NOT NULL DEFAULT 0.00 COMMENT '积分余额（可正可负）',
  `joined_at` DATETIME       NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `left_at`   DATETIME       NULL     DEFAULT NULL COMMENT 'NULL=在房间中',
  PRIMARY KEY (`id`),
  KEY `idx_room_id` (`room_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='房间成员表';

CREATE TABLE IF NOT EXISTS `bill_transactions` (
  `id`           BIGINT        NOT NULL AUTO_INCREMENT,
  `room_id`      BIGINT        NOT NULL DEFAULT 0,
  `from_user_id` BIGINT        NOT NULL DEFAULT 0 COMMENT '支出方',
  `to_user_id`   BIGINT        NOT NULL DEFAULT 0 COMMENT '收入方',
  `amount`       DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  `status`       TINYINT       NOT NULL DEFAULT 0 COMMENT '0=有效 1=已撤销',
  `thanked`      TINYINT       NOT NULL DEFAULT 0 COMMENT '0=未感谢 1=已感谢',
  `created_at`   DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `revoked_at`   DATETIME      NULL     DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_room_id` (`room_id`),
  KEY `idx_from_user_id` (`from_user_id`),
  KEY `idx_to_user_id` (`to_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='交易记录表';

CREATE TABLE IF NOT EXISTS `bill_room_logs` (
  `id`         BIGINT       NOT NULL AUTO_INCREMENT,
  `room_id`    BIGINT       NOT NULL DEFAULT 0,
  `user_id`    BIGINT       NULL     DEFAULT NULL COMMENT 'NULL=系统消息',
  `content`    VARCHAR(255) NOT NULL DEFAULT '' COMMENT '日志文本',
  `log_type`   VARCHAR(30)  NOT NULL DEFAULT '' COMMENT 'join/leave/transfer/revoke/thanks/settle',
  `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_room_id` (`room_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='房间日志表';
