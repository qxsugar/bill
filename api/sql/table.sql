-- bill 数据库初始化建表脚本
-- 库名：bill；所有业务表加前缀 bill_

CREATE TABLE IF NOT EXISTS `bill_users` (
  `id`         bigint       NOT NULL AUTO_INCREMENT,
  `openid`     varchar(64)  NOT NULL DEFAULT '' COMMENT '微信 openid',
  `nickname`   varchar(50)  NOT NULL DEFAULT '' COMMENT '昵称（首次随机生成）',
  `avatar`     varchar(255) NOT NULL DEFAULT '' COMMENT '头像 URL',
  `created_at` timestamp    NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp    NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_bill_users_openid` (`openid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

CREATE TABLE IF NOT EXISTS `bill_rooms` (
  `id`         bigint      NOT NULL AUTO_INCREMENT,
  `code`       varchar(10) NOT NULL DEFAULT '' COMMENT '房间码（4位数优先，两两连号）',
  `owner_id`   bigint      NOT NULL DEFAULT 0 COMMENT '房主 user_id',
  `status`     tinyint     NOT NULL DEFAULT 0 COMMENT '0=活跃 1=已结算',
  `settled_at` datetime    NULL DEFAULT NULL COMMENT '结算时间',
  `created_at` timestamp   NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_bill_rooms_code` (`code`),
  KEY `idx_bill_rooms_owner_id` (`owner_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='房间表';

CREATE TABLE IF NOT EXISTS `bill_room_members` (
  `id`        bigint        NOT NULL AUTO_INCREMENT,
  `room_id`   bigint        NOT NULL DEFAULT 0,
  `user_id`   bigint        NOT NULL DEFAULT 0,
  `balance`   decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '当前积分余额（可正可负）',
  `joined_at` timestamp     NULL DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',
  `left_at`   datetime      NULL DEFAULT NULL COMMENT '离开时间（NULL=在房间中）',
  PRIMARY KEY (`id`),
  KEY `idx_bill_room_members_room_id` (`room_id`),
  KEY `idx_bill_room_members_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='房间成员表';

CREATE TABLE IF NOT EXISTS `bill_transactions` (
  `id`           bigint        NOT NULL AUTO_INCREMENT,
  `room_id`      bigint        NOT NULL DEFAULT 0,
  `from_user_id` bigint        NOT NULL DEFAULT 0 COMMENT '支出方',
  `to_user_id`   bigint        NOT NULL DEFAULT 0 COMMENT '收入方',
  `amount`       decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '金额',
  `status`       tinyint       NOT NULL DEFAULT 0 COMMENT '0=有效 1=已撤销',
  `thanked`      tinyint       NOT NULL DEFAULT 0 COMMENT '0=未感谢 1=已感谢',
  `created_at`   timestamp     NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `revoked_at`   datetime      NULL DEFAULT NULL COMMENT '撤销时间',
  PRIMARY KEY (`id`),
  KEY `idx_bill_transactions_room_id` (`room_id`),
  KEY `idx_bill_transactions_from_user_id` (`from_user_id`),
  KEY `idx_bill_transactions_to_user_id` (`to_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='交易记录表';

CREATE TABLE IF NOT EXISTS `bill_room_logs` (
  `id`         bigint       NOT NULL AUTO_INCREMENT,
  `room_id`    bigint       NOT NULL DEFAULT 0,
  `user_id`    bigint       NULL DEFAULT NULL COMMENT '操作用户（NULL=系统消息）',
  `content`    varchar(255) NOT NULL DEFAULT '' COMMENT '日志文本',
  `log_type`   varchar(30)  NOT NULL DEFAULT '' COMMENT 'join/leave/transfer/revoke/thanks/settle',
  `created_at` timestamp    NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_bill_room_logs_room_id` (`room_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='房间日志表';
