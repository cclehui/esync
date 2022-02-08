```

CREATE TABLE `esync_event_default` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `event_date` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '消息日期, 20210914',
  `event_type` varchar(64) NOT NULL DEFAULT '' COMMENT '消息类型',
  `uniq_key` varchar(128) NOT NULL DEFAULT '' COMMENT '一种事件类型的唯一ID',
  `uniq_key_crc32` bigint(20) NOT NULL DEFAULT '0' COMMENT '一种事件类型的唯一ID crc32 hash',
  `event_option` varchar(500) NOT NULL DEFAULT '' COMMENT '事件的配置项',
  `event_data` text COMMENT '事件的数据',
  `e_status` tinyint(3) NOT NULL DEFAULT '1' COMMENT '消息处理状态，1:新建,5:成功,9:失败',
  `handler_info` varchar(3000) NOT NULL DEFAULT '' COMMENT '消息handler处理结果和状态跟进',
  `created_at` datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_event_date` (`event_date`),
  KEY `idx_uniq_key_crc32` (`uniq_key_crc32`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='事件持久化数据'

```
