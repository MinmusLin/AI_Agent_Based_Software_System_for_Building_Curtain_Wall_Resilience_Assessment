package common

import (
	"net/rpc"

	"icw_core_api/configs"
)

// Deps API Handler 的公共依赖集合
type Deps struct {
	Config        configs.Config
	CoreBizClient *rpc.Client
}

func NewDeps(cfg configs.Config, coreBizClient *rpc.Client) *Deps {
	return &Deps{
		Config:        cfg,
		CoreBizClient: coreBizClient,
	}
}

// BaseHandler 提供所有 API Handler 共享的基础依赖
type BaseHandler struct {
	deps *Deps
}

func NewBaseHandler(deps *Deps) *BaseHandler {
	if deps == nil {
		return &BaseHandler{}
	}
	return &BaseHandler{
		deps: deps,
	}
}

// Config 获取服务配置
func (h *BaseHandler) Config() configs.Config {
	return h.deps.Config
}

// CoreBizClient 获取 icw.core.biz RPC Client
func (h *BaseHandler) CoreBizClient() *rpc.Client {
	return h.deps.CoreBizClient
}
