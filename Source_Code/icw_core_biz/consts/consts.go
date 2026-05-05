package consts

const (
	// CoreBizPSM icw.core.biz 服务标识
	CoreBizPSM = "icw.core.biz"
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
)
