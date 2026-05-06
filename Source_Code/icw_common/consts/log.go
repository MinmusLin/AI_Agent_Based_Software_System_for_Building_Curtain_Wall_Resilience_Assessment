package consts

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
	// LogColorBoldBlue ANSI 终端颜色码：蓝色 [CBK INFO]
	LogColorBoldBlue = "\033[1;34m"
	// LogColorBoldPurple ANSI 终端颜色码：紫色 [WS INFO | CRON INFO | CLS INFO | RSN INFO]
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
	// LogScopeWS WebSocket 日志域
	LogScopeWS = "WS"
	// LogScopeCron 定时任务日志域
	LogScopeCron = "CRON"
	// LogScopeClassification 分类能力日志域
	LogScopeClassification = "CLS"
	// LogScopeReasoning 检测能力日志域
	LogScopeReasoning = "RSN"
	// LogScopeCallback 回调日志域
	LogScopeCallback = "CBK"
)
