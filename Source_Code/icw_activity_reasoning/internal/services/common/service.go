package common

import (
	"context"

	"icw_common/rpc/error"
	"icw_common/utils"

	"icw_activity_reasoning/configs"
	"icw_activity_reasoning/internal/detectors/common"
	"icw_activity_reasoning/rpc/icw_core_biz"
)

// Deps RPC Service 的公共依赖集合
type Deps struct {
	Config    configs.Config
	Registry  *common.Registry
	CoreBiz   *icw_core_biz.Client
	Semaphore chan struct{}
}

// NewDeps 创建 RPC Service 的公共依赖集合
func NewDeps(config configs.Config, registry *common.Registry, coreBiz *icw_core_biz.Client) *Deps {
	return &Deps{
		Config:    config,
		Registry:  registry,
		CoreBiz:   coreBiz,
		Semaphore: make(chan struct{}, config.ReasoningTaskMaxConcurrency),
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
		return rpc_error.BadRequestDefault("request is nil")
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

// Registry 获取原子检测能力注册表
func (s *BaseService) Registry() *common.Registry {
	if s == nil || s.deps == nil {
		return nil
	}
	return s.deps.Registry
}

// CoreBizClient 获取 icw.core.biz RPC Client
func (s *BaseService) CoreBizClient() *icw_core_biz.Client {
	if s == nil || s.deps == nil {
		return nil
	}
	return s.deps.CoreBiz
}

// Acquire 获取并发执行额度
func (s *BaseService) Acquire() {
	if s == nil || s.deps == nil || s.deps.Semaphore == nil {
		return
	}
	s.deps.Semaphore <- struct{}{}
}

// Release 释放并发执行额度
func (s *BaseService) Release() {
	if s == nil || s.deps == nil || s.deps.Semaphore == nil {
		return
	}
	<-s.deps.Semaphore
}
