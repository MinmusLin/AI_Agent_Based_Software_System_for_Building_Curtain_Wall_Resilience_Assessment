package cronjobs

import (
	"context"

	"icw_core_biz/internal/cronjobs/common"
	"icw_core_biz/internal/cronjobs/jobs"
)

// Start 启动定时任务
func Start(ctx context.Context, deps *common.Deps) {
	if deps == nil {
		deps = &common.Deps{}
	}

	cronJobs := registry(deps)
	common.CronInfo("Cron jobs registered, waiting for schedule:\n%s", common.FormatRegistryTable(cronJobs))
	for _, item := range cronJobs {
		if err := common.Start(ctx, deps, item.PSM, item.Cron, item.Start); err != nil {
			common.CronFatal("Failed to start cron job %s: %v", item.PSM, err)
		}
	}
}

// registry 定时任务注册表
func registry(deps *common.Deps) []common.CronJobMeta {
	return []common.CronJobMeta{
		{
			PSM:         "icw.cron.pending_image_timeout",
			Description: "上传中图像超时失败任务",
			Cron:        deps.Config.PendingImageTimeoutJobCron,
			Start:       jobs.NewPendingImageTimeoutJob,
		},
		{
			PSM:         "icw.cron.minio_dirty_object_cleanup",
			Description: "MinIO 脏对象清理任务",
			Cron:        deps.Config.MinIODirtyObjectCleanupJobCron,
			Start:       jobs.NewMinIODirtyObjectCleanupJob,
		},
	}
}
