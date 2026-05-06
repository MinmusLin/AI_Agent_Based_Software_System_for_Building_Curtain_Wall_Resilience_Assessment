package common

import (
	"context"
	"icw_common/rpc_err"
	"icw_common/utils"
	"icw_core_biz/configs"
	"icw_core_biz/repositories/minio"
	"icw_core_biz/repositories/mysql"
	"icw_core_biz/repositories/redis"
	"icw_core_biz/repositories/rocketmq"
	"icw_core_biz/repositories/smtp"
	"icw_core_biz/rpc/icw_activity_classification"
	"icw_core_biz/rpc/icw_activity_reasoning"
	"icw_core_biz/rpc/icw_activity_summary"
)

// Deps RPC Service 的公共依赖集合
type Deps struct {
	Config                       configs.Config
	MySQL                        *mysql.Repository
	Redis                        *redis.Repository
	RocketMQ                     *rocketmq.Repository
	MinIO                        *minio.Repository
	SMTP                         *smtp.Repository
	ActivityClassificationClient *icw_activity_classification.Client
	ActivityReasoningClient      *icw_activity_reasoning.Client
	ActivitySummaryClient        *icw_activity_summary.Client
}

// NewDeps 创建 RPC Service 的公共依赖集合
func NewDeps(
	Config configs.Config,
	MySQL *mysql.Repository,
	Redis *redis.Repository,
	RocketMQ *rocketmq.Repository,
	MinIO *minio.Repository,
	SMTP *smtp.Repository,
	ActivityClassificationClient *icw_activity_classification.Client,
	ActivityReasoningClient *icw_activity_reasoning.Client,
	ActivitySummaryClient *icw_activity_summary.Client,
) *Deps {
	return &Deps{
		Config:                       Config,
		MySQL:                        MySQL,
		Redis:                        Redis,
		RocketMQ:                     RocketMQ,
		MinIO:                        MinIO,
		SMTP:                         SMTP,
		ActivityClassificationClient: ActivityClassificationClient,
		ActivityReasoningClient:      ActivityReasoningClient,
		ActivitySummaryClient:        ActivitySummaryClient,
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

// MySQL 获取 MySQL 服务
func (s *BaseService) MySQL() *mysql.Repository {
	if s == nil || s.deps == nil {
		return nil
	}
	return s.deps.MySQL
}

// Redis 获取 Redis 服务
func (s *BaseService) Redis() *redis.Repository {
	if s == nil || s.deps == nil {
		return nil
	}
	return s.deps.Redis
}

// RocketMQ 获取 RocketMQ 服务
func (s *BaseService) RocketMQ() *rocketmq.Repository {
	if s == nil || s.deps == nil {
		return nil
	}
	return s.deps.RocketMQ
}

// MinIO 获取 MinIO 服务
func (s *BaseService) MinIO() *minio.Repository {
	if s == nil || s.deps == nil {
		return nil
	}
	return s.deps.MinIO
}

// SMTP 获取 SMTP 服务
func (s *BaseService) SMTP() *smtp.Repository {
	if s == nil || s.deps == nil {
		return nil
	}
	return s.deps.SMTP
}

// ActivityClassificationClient 获取 icw.activity.classification RPC Client
func (s *BaseService) ActivityClassificationClient() *icw_activity_classification.Client {
	if s == nil || s.deps == nil {
		return nil
	}
	return s.deps.ActivityClassificationClient
}

// ActivityReasoningClient 获取 icw.activity.reasoning RPC Client
func (s *BaseService) ActivityReasoningClient() *icw_activity_reasoning.Client {
	if s == nil || s.deps == nil {
		return nil
	}
	return s.deps.ActivityReasoningClient
}

// ActivitySummaryClient 获取 icw.activity.summary RPC Client
func (s *BaseService) ActivitySummaryClient() *icw_activity_summary.Client {
	if s == nil || s.deps == nil {
		return nil
	}
	return s.deps.ActivitySummaryClient
}
