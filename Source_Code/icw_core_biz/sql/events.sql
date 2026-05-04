-- 开启 MySQL Event Scheduler。
-- 注意：SET GLOBAL 只对当前 MySQL 实例运行期生效；如果需要重启后仍生效，需要在 MySQL 配置文件中配置 event_scheduler=ON。
SET GLOBAL event_scheduler = ON;

-- 为 pending 图像超时扫描补充索引。
-- 使用动态 SQL 做幂等保护，避免重复执行脚本时因为索引已存在而中断。
SET @index_exists = (
  SELECT COUNT(1)
  FROM information_schema.statistics
  WHERE table_schema = DATABASE()
    AND table_name = 'project_group_images'
    AND index_name = 'idx_project_group_images_status_created_at'
);
SET @create_index_sql = IF(
  @index_exists = 0,
  'ALTER TABLE project_group_images ADD INDEX idx_project_group_images_status_created_at (status, created_at)',
  'SELECT ''idx_project_group_images_status_created_at already exists'''
);
PREPARE create_index_stmt FROM @create_index_sql;
EXECUTE create_index_stmt;
DEALLOCATE PREPARE create_index_stmt;

-- 重新创建 pending 图像超时失败事件。
DROP EVENT IF EXISTS ev_project_group_images_pending_to_failed;

DELIMITER $$

CREATE EVENT ev_project_group_images_pending_to_failed
ON SCHEDULE EVERY 1 MINUTE
ON COMPLETION PRESERVE
ENABLE
DO
BEGIN
  UPDATE project_group_images
  SET status = 'failed'
  WHERE status = 'pending'
    AND created_at < NOW(3) - INTERVAL 5 MINUTE;
END$$

DELIMITER ;
