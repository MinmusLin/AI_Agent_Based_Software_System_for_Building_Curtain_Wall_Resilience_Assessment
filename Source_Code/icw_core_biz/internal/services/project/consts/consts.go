package consts

const (
	// DefaultProjectName 默认项目名称
	DefaultProjectName = "新项目"
)

const (
	// ProgressInitializationFinished 项目初始化完成，当前项目基础信息阶段
	ProgressInitializationFinished uint8 = 0
	// ProgressProfileFinished 项目基础信息完成，当前图像资产构建阶段
	ProgressProfileFinished uint8 = 1
	// ProgressAssetsFinished 图像资产构建完成，当前 Agent 智能检测阶段
	ProgressAssetsFinished uint8 = 2
	// ProgressDetectFinished Agent 智能检测完成，当前人工复核确认阶段
	ProgressDetectFinished uint8 = 3
	// ProgressReviewFinished 人工复核确认完成，当前评估报告生成阶段
	ProgressReviewFinished uint8 = 4
	// ProgressReportFinished 评估报告生成完成，当前项目已完成
	ProgressReportFinished uint8 = 5
)
