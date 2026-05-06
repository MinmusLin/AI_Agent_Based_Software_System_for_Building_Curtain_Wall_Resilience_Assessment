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
	// SocketScopeProjectAssets 图像资产 WebSocket 连接范围
	SocketScopeProjectAssets = "ws_project_assets"
	// EventTypeProjectImageStatusChanged 项目图像状态变化事件类型
	EventTypeProjectImageStatusChanged = "project_image_status_changed"
	// EventTagProjectImageStatusChanged 项目图像状态变化事件 Tag
	EventTagProjectImageStatusChanged = "PROJECT_IMAGE_STATUS_CHANGED"
)

const (
	// LogColorReset ANSI 终端颜色重置码
	LogColorReset = "\033[0m"
	// LogColorBoldRed ANSI 终端颜色码：红色 [ERROR]
	LogColorBoldRed = "\033[1;31m"
	// LogColorBoldYellow ANSI 终端颜色码：黄色 [WARN]
	LogColorBoldYellow = "\033[1;33m"
	// LogColorBoldGreen ANSI 终端颜色码：绿色 [HTTP INFO | RPC INFO]
	LogColorBoldGreen = "\033[1;32m"
	// LogColorBoldCyan ANSI 终端颜色码：青色 [MQ INFO]
	LogColorBoldCyan = "\033[1;36m"
	// LogColorBoldBlue ANSI 终端颜色码：蓝色 [CALLBACK INFO]
	LogColorBoldBlue = "\033[1;34m"
	// LogColorBoldPurple ANSI 终端颜色码：紫色 [WS INFO | CRON INFO | REASONING INFO]
	LogColorBoldPurple = "\033[1;35m"
	// LogColorBoldWhiteOnRed ANSI 终端颜色码：白色（红色背景）[FATAL]
	LogColorBoldWhiteOnRed = "\033[1;37;41m"
	// LogColorBoldBlackOnWhite ANSI 终端颜色码：黑色（白色背景）[COST]
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
	// LogScopeReasoning 原子检测能力日志域
	LogScopeReasoning = "REASONING"
	// LogScopeCallback 回调日志域
	LogScopeCallback = "CALLBACK"
)
