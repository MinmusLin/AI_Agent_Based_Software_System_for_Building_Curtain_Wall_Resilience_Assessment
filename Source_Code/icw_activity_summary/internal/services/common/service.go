package common

import (
	"context"

	"icw_common/rpc/error"
	"icw_common/utils"

	"icw_activity_summary/configs"
	"icw_activity_summary/rpc/icw_core_biz"
)

// Deps RPC Service 的公共依赖集合
type Deps struct {
	Config                    configs.Config
	CoreBiz                   *icw_core_biz.Client
	DetectionSummarySemaphore chan struct{}
	ProjectSummarySemaphore   chan struct{}
}

// NewDeps 创建 RPC Service 的公共依赖集合
func NewDeps(config configs.Config, coreBiz *icw_core_biz.Client) *Deps {
	return &Deps{
		Config:                    config,
		CoreBiz:                   coreBiz,
		DetectionSummarySemaphore: make(chan struct{}, config.DetectionSummaryTaskMaxConcurrency),
		ProjectSummarySemaphore:   make(chan struct{}, config.ProjectSummaryTaskMaxConcurrency),
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

// CoreBizClient 获取 icw.core.biz RPC Client
func (s *BaseService) CoreBizClient() *icw_core_biz.Client {
	if s == nil || s.deps == nil {
		return nil
	}
	return s.deps.CoreBiz
}

// AcquireDetectionSummary 获取图像检测总结并发执行额度
func (s *BaseService) AcquireDetectionSummary() {
	if s == nil || s.deps == nil || s.deps.DetectionSummarySemaphore == nil {
		return
	}
	s.deps.DetectionSummarySemaphore <- struct{}{}
}

// ReleaseDetectionSummary 释放图像检测总结并发执行额度
func (s *BaseService) ReleaseDetectionSummary() {
	if s == nil || s.deps == nil || s.deps.DetectionSummarySemaphore == nil {
		return
	}
	<-s.deps.DetectionSummarySemaphore
}

// AcquireProjectSummary 获取项目总结并发执行额度
func (s *BaseService) AcquireProjectSummary() {
	if s == nil || s.deps == nil || s.deps.ProjectSummarySemaphore == nil {
		return
	}
	s.deps.ProjectSummarySemaphore <- struct{}{}
}

// ReleaseProjectSummary 释放项目总结并发执行额度
func (s *BaseService) ReleaseProjectSummary() {
	if s == nil || s.deps == nil || s.deps.ProjectSummarySemaphore == nil {
		return
	}
	<-s.deps.ProjectSummarySemaphore
}
