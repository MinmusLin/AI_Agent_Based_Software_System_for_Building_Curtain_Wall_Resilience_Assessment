package cronjobs

import (
	"context"
	"fmt"
	"strings"

	"icw_core_biz/internal/cronjobs/common"
	"icw_core_biz/internal/cronjobs/jobs"
)

// CronJob 定时任务配置
type CronJob struct {
	name  string
	cron  string
	start common.JobFactory
}

// Start 启动定时任务
func Start(ctx context.Context, deps *common.Deps) {
	if deps == nil {
		deps = &common.Deps{}
	}

	for _, item := range register(deps) {
		if err := common.Start(ctx, deps, item.name, item.cron, item.start); err != nil {
			common.CronFault("Failed to start cron job %s: %v", item.name, err)
		}
	}
}

// register 注册定时任务
func register(deps *common.Deps) []CronJob {
	cronJobs := registry(deps)
	common.CronInfo("Cron job registered, waiting for schedule:\n%s", formatRegistryTable(cronJobs))
	return cronJobs
}

// registry 定时任务注册表
func registry(deps *common.Deps) []CronJob {
	return []CronJob{
		{
			name:  "icw.cron.pending_image_timeout",
			cron:  deps.Config.PendingImageTimeoutJobCron,
			start: jobs.NewPendingImageTimeoutJob,
		},
	}
}

// formatRegistryTable 将定时任务注册表格式化为表格
func formatRegistryTable(cronJobs []CronJob) string {
	const (
		nameHeader       = "job name"
		expressionHeader = "cron expression"
	)
	nameWidth := len(nameHeader)
	expressionWidth := len(expressionHeader)
	for _, item := range cronJobs {
		nameWidth = max(nameWidth, len(item.name))
		expressionWidth = max(expressionWidth, len(item.cron))
	}
	lineIndent := strings.Repeat(" ", 20)
	border := fmt.Sprintf("+-%s-+-%s-+", strings.Repeat("-", nameWidth), strings.Repeat("-", expressionWidth))
	var builder strings.Builder
	builder.WriteString(lineIndent)
	builder.WriteString(border)
	builder.WriteString("\n")
	builder.WriteString(lineIndent)
	builder.WriteString(fmt.Sprintf("| %-*s | %-*s |\n", nameWidth, nameHeader, expressionWidth, expressionHeader))
	builder.WriteString(lineIndent)
	builder.WriteString(border)
	builder.WriteString("\n")
	for _, item := range cronJobs {
		builder.WriteString(lineIndent)
		builder.WriteString(fmt.Sprintf("| %-*s | %-*s |\n", nameWidth, item.name, expressionWidth, item.cron))
	}
	builder.WriteString(lineIndent)
	builder.WriteString(border)
	return builder.String()
}
