package consts

import (
	"icw_common/gen/activity/classification"
	"icw_common/gen/activity/reasoning"
	"icw_common/gen/activity/summary"
	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
)

var (
	// CoreApiPSM icw.core.api 服务标识
	CoreApiPSM = string(apipb.File_core_api_proto.Package())
	// CoreBizPSM icw.core.biz 服务标识
	CoreBizPSM = string(bizpb.File_core_biz_proto.Package())
	// ActivityClassificationPSM icw.activity.classification 服务标识
	ActivityClassificationPSM = string(classificationpb.File_activity_classification_proto.Package())
	// ActivityReasoningPSM icw.activity.reasoning 服务标识
	ActivityReasoningPSM = string(reasoningpb.File_activity_reasoning_proto.Package())
	// ActivitySummaryPSM icw.activity.summary 服务标识
	ActivitySummaryPSM = string(summarypb.File_activity_summary_proto.Package())
)

const (
	// GRPCMetadataRequestId gRPC metadata 中透传请求 ID 的 Key
	GRPCMetadataRequestId = "x-request-id"
)

// ProjectProgress 项目进度枚举
type ProjectProgress uint8

const (
	// ProjectProgressInitializationFinished 项目初始化完成，当前项目基础信息阶段
	ProjectProgressInitializationFinished ProjectProgress = 0
	// ProjectProgressProfileFinished 项目基础信息完成，当前图像资产构建阶段
	ProjectProgressProfileFinished ProjectProgress = 1
	// ProjectProgressAssetsFinished 图像资产构建完成，当前 Agent 智能检测阶段
	ProjectProgressAssetsFinished ProjectProgress = 2
	// ProjectProgressDetectionFinished Agent 智能检测完成，当前人工复核确认阶段
	ProjectProgressDetectionFinished ProjectProgress = 3
	// ProjectProgressReviewFinished 人工复核确认完成，当前评估报告生成阶段
	ProjectProgressReviewFinished ProjectProgress = 4
	// ProjectProgressReportFinished 评估报告生成完成，当前项目已完成
	ProjectProgressReportFinished ProjectProgress = 5
)

// Uint8 将项目进度枚举转换为 uint8
func (p ProjectProgress) Uint8() uint8 {
	return uint8(p)
}

// ParseProjectProgress 将外部输入转换为项目进度枚举
func ParseProjectProgress(value uint8) ProjectProgress {
	switch progress := ProjectProgress(value); progress {
	case ProjectProgressInitializationFinished,
		ProjectProgressProfileFinished,
		ProjectProgressAssetsFinished,
		ProjectProgressDetectionFinished,
		ProjectProgressReviewFinished,
		ProjectProgressReportFinished:
		return progress
	default:
		return ProjectProgressInitializationFinished
	}
}

const (
	// EventTypeProjectImageStatusChanged 项目图像状态变化事件类型
	EventTypeProjectImageStatusChanged = "project_image_status_changed"
	// EventTagProjectImageStatusChanged 项目图像状态变化事件 Tag
	EventTagProjectImageStatusChanged = "PROJECT_IMAGE_STATUS_CHANGED"
	// SocketScopeProjectAssets 图像资产 WebSocket 连接范围
	SocketScopeProjectAssets = "ws_project_assets"
)

const (
	// LogColorReset ANSI 终端颜色重置码
	LogColorReset = "\033[0m"
	// LogColorBoldRed ANSI 终端颜色码：红色
	LogColorBoldRed = "\033[1;31m"
	// LogColorBoldYellow ANSI 终端颜色码：黄色
	LogColorBoldYellow = "\033[1;33m"
	// LogColorBoldGreen ANSI 终端颜色码：绿色
	LogColorBoldGreen = "\033[1;32m"
	// LogColorBoldCyan ANSI 终端颜色码：青色
	LogColorBoldCyan = "\033[1;36m"
	// LogColorBoldBlue ANSI 终端颜色码：蓝色
	LogColorBoldBlue = "\033[1;34m"
	// LogColorBoldPurple ANSI 终端颜色码：紫色
	LogColorBoldPurple = "\033[1;35m"
	// LogColorBoldPink ANSI 终端颜色码：粉色
	LogColorBoldPink = "\033[1;95m"
	// LogColorBoldWhiteOnRed ANSI 终端颜色码：白色（红色背景）
	LogColorBoldWhiteOnRed = "\033[1;37;41m"
	// LogColorBoldBlackOnWhite ANSI 终端颜色码：黑色（白色背景）
	LogColorBoldBlackOnWhite = "\033[30;47m"
)

const (
	// LogScopeInit 服务初始化日志域
	LogScopeInit = "INIT"
	// LogScopeHTTP HTTP 服务日志域
	LogScopeHTTP = "HTTP"
	// LogScopeRPC RPC 服务日志域
	LogScopeRPC = "RPC"
	// LogScopeMQ 消息队列日志域
	LogScopeMQ = "MQ"
	// LogScopeCron 定时任务日志域
	LogScopeCron = "CRON"
	// LogScopeWS WebSocket 日志域
	LogScopeWS = "WS"
	// LogScopeReasoning 推理能力日志域
	LogScopeReasoning = "REASONING"
	// LogScopeCallback 回调日志域
	LogScopeCallback = "CALLBACK"
)
