package cronjobs

import (
	"context"
	"fmt"
	"strings"

	"icw_core_biz/internal/cronjobs/common"
	"icw_core_biz/internal/cronjobs/jobs"
	"icw_core_biz/utils"
)

// cronJob 定时任务配置
type cronJob struct {
	name        string
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
		if err := common.Start(ctx, deps, item.name, item.cron, item.start); err != nil {
			common.CronFatal("Failed to start cron job %s: %v", item.name, err)
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
			name:        "icw.cron.pending_image_timeout",
			description: "上传中图像超时失败任务",
			cron:        deps.Config.PendingImageTimeoutJobCron,
			start:       jobs.NewPendingImageTimeoutJob,
		},
	}
}

// formatRegistryTable 将定时任务注册表格式化为表格
func formatRegistryTable(cronJobs []cronJob) string {
	const (
		psmHeader            = "PSM"
		descriptionHeader    = "description"
		cronExpressionHeader = "cronExpression"
	)
	psmWidth := len(psmHeader)
	descriptionWidth := len(descriptionHeader)
	cronExpressionWidth := len(cronExpressionHeader)
	for _, item := range cronJobs {
		psmWidth = max(psmWidth, utils.DisplayWidth(item.name))
		descriptionWidth = max(descriptionWidth, utils.DisplayWidth(item.description))
		cronExpressionWidth = max(cronExpressionWidth, utils.DisplayWidth(item.cron))
	}
	border := fmt.Sprintf(
		"+-%s-+-%s-+-%s-+",
		strings.Repeat("-", psmWidth),
		strings.Repeat("-", descriptionWidth),
		strings.Repeat("-", cronExpressionWidth),
	)
	var builder strings.Builder
	builder.WriteString(border)
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("| %s | %s | %s |\n", utils.PadRight(psmHeader, psmWidth), utils.PadRight(descriptionHeader, descriptionWidth), utils.PadRight(cronExpressionHeader, cronExpressionWidth)))
	builder.WriteString(border)
	builder.WriteString("\n")
	for _, item := range cronJobs {
		builder.WriteString(fmt.Sprintf("| %s | %s | %s |\n", utils.PadRight(item.name, psmWidth), utils.PadRight(item.description, descriptionWidth), utils.PadRight(item.cron, cronExpressionWidth)))
	}
	builder.WriteString(border)
	return builder.String()
}
