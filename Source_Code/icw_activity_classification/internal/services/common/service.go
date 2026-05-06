package common

import (
	"context"

	"icw_activity_classification/configs"
	"icw_activity_classification/rpc/icw_core_biz"
	"icw_common/rpc_err"
	"icw_common/utils"
)

// Deps RPC Service 的公共依赖集合
type Deps struct {
	Config  configs.Config
	CoreBiz *icw_core_biz.Client
}

// NewDeps 创建 RPC Service 的公共依赖集合
func NewDeps(config configs.Config, coreBiz *icw_core_biz.Client) *Deps {
	return &Deps{
		Config:  config,
		CoreBiz: coreBiz,
	}
}

// BaseService 所有 RPC Service 共享的基础依赖
type BaseService struct {
	deps *Deps
	ctx  context.Context
}

// NewBaseService 创建所有 RPC Service 共享的基础依赖
func NewBaseService(ctx context.Context, deps *Deps) *BaseService {
	if ctx == nil {
		ctx = context.Background()
	}
	if deps == nil {
		deps = &Deps{}
	}
	return &BaseService{
		deps: deps,
		ctx:  ctx,
	}
}

// CallRPC RPC 服务通用调用
func (s *BaseService) CallRPC(ctx context.Context, req interface{}, fn func() error) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if utils.IsNil(req) {
		return rpc_err.BadRequestDefault("request is nil")
	}
	if fn == nil {
		return nil
	}
	return fn()
}

// Ctx 获取上下文
func (s *BaseService) Ctx() context.Context {
	if s == nil || s.ctx == nil {
		return context.Background()
	}
	return s.ctx
}

// Config 获取服务配置
func (s *BaseService) Config() configs.Config {
	if s == nil || s.deps == nil {
		return configs.Config{}
	}
	return s.deps.Config
}

// CoreBizClient 获取 icw.core.biz RPC Client
func (s *BaseService) CoreBizClient() *icw_core_biz.Client {
	if s == nil || s.deps == nil {
		return nil
	}
	return s.deps.CoreBiz
}
