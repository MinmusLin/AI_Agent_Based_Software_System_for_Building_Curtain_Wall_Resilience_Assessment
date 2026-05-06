package common

import (
	"icw_common/utils"
)

// CronJobMeta 定时任务配置
type CronJobMeta struct {
	PSM         string
	Description string
	Cron        string
	Start       JobFactory
}

// FormatRegistryTable 将定时任务注册表格式化为表格
func FormatRegistryTable(cronJobs []CronJobMeta) string {
	psmValues := make([]string, 0, len(cronJobs))
	descriptionValues := make([]string, 0, len(cronJobs))
	cronExpressionValues := make([]string, 0, len(cronJobs))
	for _, item := range cronJobs {
		psmValues = append(psmValues, item.PSM)
		descriptionValues = append(descriptionValues, item.Description)
		cronExpressionValues = append(cronExpressionValues, item.Cron)
	}
	return utils.FormatTable([]*utils.TableColumn{
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
