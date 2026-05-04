package cronjobs

import (
	"context"

	"icw_core_biz/internal/cronjobs/common"
	"icw_core_biz/internal/cronjobs/jobs"
)

// job 定时任务配置
type job struct {
	name  string
	cron  string
	start func(context.Context, *common.Deps, string) error
}

// Start 启动定时任务
func Start(ctx context.Context, deps *common.Deps) {
	if deps == nil {
		deps = &common.Deps{}
	}
	for _, item := range register(deps) {
		common.CronInfo("Start cronjob: %s, cron: %s", item.name, item.cron)
		if err := item.start(ctx, deps, item.cron); err != nil {
			common.CronFault("Failed to start cronjob %s: %v", item.name, err)
		}
	}
}

// register 注册定时任务
func register(deps *common.Deps) []job {
	return []job{
		{
			name:  "icw.cron.pending_image_timeout",
			cron:  deps.Config.PendingImageTimeoutJobCron,
			start: jobs.StartPendingImageTimeoutJob,
		},
	}
}
