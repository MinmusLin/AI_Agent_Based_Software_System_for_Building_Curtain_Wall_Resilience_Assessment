package cronjobs

import (
	"context"

	"icw_core_biz/internal/cronjobs/common"
	"icw_core_biz/internal/cronjobs/jobs"
	"icw_core_biz/utils"
)

// cronJob 定时任务配置
type cronJob struct {
	psm         string
	description string
	cron        string
	start       common.JobFactory
}

// Start 启动定时任务
func Start(ctx context.Context, deps *common.Deps) {
	if deps == nil {
		deps = &common.Deps{}
	}

	for _, item := range register(deps) {
		if err := common.Start(ctx, deps, item.psm, item.cron, item.start); err != nil {
			common.CronFatal("Failed to start cron job %s: %v", item.psm, err)
		}
	}
}

// register 注册定时任务
func register(deps *common.Deps) []cronJob {
	cronJobs := registry(deps)
	common.CronInfo("Cron job registered, waiting for schedule:\n%s", formatRegistryTable(cronJobs))
	return cronJobs
}

// registry 定时任务注册表
func registry(deps *common.Deps) []cronJob {
	return []cronJob{
		{
			psm:         "icw.cron.pending_image_timeout",
			description: "上传中图像超时失败任务",
			cron:        deps.Config.PendingImageTimeoutJobCron,
			start:       jobs.NewPendingImageTimeoutJob,
		},
	}
}

// formatRegistryTable 将定时任务注册表格式化为表格
func formatRegistryTable(cronJobs []cronJob) string {
	psmValues := make([]string, 0, len(cronJobs))
	descriptionValues := make([]string, 0, len(cronJobs))
	cronExpressionValues := make([]string, 0, len(cronJobs))
	for _, item := range cronJobs {
		psmValues = append(psmValues, item.psm)
		descriptionValues = append(descriptionValues, item.description)
		cronExpressionValues = append(cronExpressionValues, item.cron)
	}
	return utils.FormatTable([]utils.TableColumn{
		{
			Header: "psm",
			Values: psmValues,
		},
		{
			Header: "description",
			Values: descriptionValues,
		},
		{
			Header: "cron expression",
			Values: cronExpressionValues,
		},
	})
}
