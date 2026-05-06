package common

import (
	"icw_core_api/configs"
	"icw_core_api/rpc/icw_core_biz"
)

// Deps API Handler 的公共依赖集合
type Deps struct {
	Config        configs.Config
	CoreBizClient *icw_core_biz.Client
}

// NewDeps 创建 API Handler 的公共依赖集合
func NewDeps(cfg configs.Config, coreBizClient *icw_core_biz.Client) *Deps {
	return &Deps{
		Config:        cfg,
		CoreBizClient: coreBizClient,
	}
}

// BaseHandler 所有 API Handler 共享的基础依赖
type BaseHandler struct {
	deps *Deps
}

// NewBaseHandler 创建所有 API Handler 共享的基础依赖
func NewBaseHandler(deps *Deps) *BaseHandler {
	if deps == nil {
		deps = &Deps{}
	}
	return &BaseHandler{
		deps: deps,
	}
}

// Config 获取服务配置
func (h *BaseHandler) Config() configs.Config {
	if h == nil || h.deps == nil {
		return configs.Config{}
	}
	return h.deps.Config
}

// CoreBizClient 获取 icw.core.biz gRPC Client
func (h *BaseHandler) CoreBizClient() *icw_core_biz.Client {
	if h == nil || h.deps == nil {
		return nil
	}
	return h.deps.CoreBizClient
}
