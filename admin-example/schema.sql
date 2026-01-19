-- schema.sql
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` longtext,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL, -- 論理削除用
  PRIMARY KEY (`id`),
  INDEX `idx_users_deleted_at` (`deleted_at`) -- 削除されていないデータを探しやすくする
);