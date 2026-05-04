-- ev_project_group_images_pending_to_failed 项目图像上传超时状态流转任务
CREATE DEFINER=`mysql-user`@`%` EVENT `ev_project_group_images_pending_to_failed` ON SCHEDULE EVERY 1 MINUTE STARTS '2026-05-04 12:53:27' ON COMPLETION PRESERVE ENABLE DO BEGIN
  UPDATE project_group_images
  SET status = 'failed'
  WHERE status = 'pending'
    AND created_at < NOW(3) - INTERVAL 5 MINUTE;
END;
