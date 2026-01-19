-- schema.sql
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  /* 重複登録（二重送信）をDBレベルでも防ぐためのユニーク制約 */
  /* 論理削除（deleted_at）を含めた複合ユニークにすることで、削除済みなら同じ名前で再登録可能 */
  UNIQUE KEY `uk_name_deleted_at` (`name`, `deleted_at`),
  PRIMARY KEY (`id`),
  INDEX `idx_users_deleted_at` (`deleted_at`)
);
