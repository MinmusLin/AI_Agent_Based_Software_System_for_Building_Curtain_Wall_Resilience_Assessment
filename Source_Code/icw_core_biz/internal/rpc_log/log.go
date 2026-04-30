package rpc_log

import (
	"log"
	"time"

	"icw_core_biz/utils"
)

// Start 记录 RPC 请求入口日志
func Start(method string, req interface{}) time.Time {
	start := time.Now()
	log.Printf("[RPC] [%s] Call", method)
	return start
}

// Finish 记录 RPC 请求完成日志
func Finish(method string, req interface{}, resp interface{}, start time.Time, err error) {
	if err != nil {
		log.Printf("[RPC] [%s] cost=%s req=%s resp=%s err=%v", method, time.Since(start), utils.JSONF(req), utils.JSONF(resp), err)
		return
	}
	log.Printf("[RPC] [%s] cost=%s req=%s resp=%s err=<nil>", method, time.Since(start), utils.JSONF(req), utils.JSONF(resp))
}
