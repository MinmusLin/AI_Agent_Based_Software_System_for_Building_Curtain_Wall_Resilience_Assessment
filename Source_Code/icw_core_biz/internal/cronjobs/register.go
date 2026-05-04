package cronjobs

import (
	"context"

	"icw_core_biz/internal/cronjobs/common"
	"icw_core_biz/internal/cronjobs/jobs"
)

// CronJob 定时任务配置
type CronJob struct {
	name  string
	cron  string
	start common.CronJobFunc
}

// Start 启动定时任务
func Start(ctx context.Context, deps *common.Deps) {
	if deps == nil {
		deps = &common.Deps{}
	}
	for _, item := range register(deps) {
		if err := common.Start(ctx, deps, item.name, item.cron, item.start); err != nil {
			common.CronFault("Failed to register cron job %s: %v", item.name, err)
		}
	}
}

// register 注册定时任务
func register(deps *common.Deps) []CronJob {
	return []CronJob{
		{
			name:  "icw.cron.pending_image_timeout",
			cron:  deps.Config.PendingImageTimeoutJobCron,
			start: jobs.StartPendingImageTimeoutJob,
		},
	}
}
