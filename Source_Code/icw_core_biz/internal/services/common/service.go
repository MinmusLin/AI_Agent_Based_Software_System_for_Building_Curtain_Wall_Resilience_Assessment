package common

import (
	"context"
	"reflect"
	"time"

	"icw_core_biz/configs"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/repositories/minio"
	"icw_core_biz/repositories/mysql"
	"icw_core_biz/repositories/redis"
	"icw_core_biz/repositories/smtp"
)

// Deps RPC Service 的公共依赖集合
type Deps struct {
	Config configs.Config
	MySQL  *mysql.Repository
	Redis  *redis.Repository
	SMTP   *smtp.Repository
	MinIO  *minio.Repository
}

func NewDeps(Config configs.Config, MySQL *mysql.Repository, Redis *redis.Repository, SMTP *smtp.Repository, MinIO *minio.Repository) *Deps {
	return &Deps{
		Config: Config,
		MySQL:  MySQL,
		Redis:  Redis,
		SMTP:   SMTP,
		MinIO:  MinIO,
	}
}

// BaseService 提供所有 RPC Service 共享的基础依赖
type BaseService struct {
	deps *Deps
	ctx  context.Context
}

func NewBaseService(ctx context.Context, deps *Deps) *BaseService {
	if deps == nil {
		deps = &Deps{}
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return &BaseService{
		deps: deps,
		ctx:  ctx,
	}
}

// CallRPC RPC 服务通用调用
func (s *BaseService) CallRPC(method string, req interface{}, resp interface{}, fn func() error) (err error) {
	start := time.Now()
	defer func() {
		rpcLog(method, req, resp, start, err)
	}()

	if req == nil {
		return rpc_err.BadRequestDefault("request is nil")
	}

	value := reflect.ValueOf(req)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		if value.IsNil() {
			return rpc_err.BadRequestDefault("request is nil")
		}
	default:
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
	return s.deps.Config
}

// MySQL 获取 MySQL 服务
func (s *BaseService) MySQL() *mysql.Repository {
	return s.deps.MySQL
}

// Redis 获取 Redis 服务
func (s *BaseService) Redis() *redis.Repository {
	return s.deps.Redis
}

// SMTP 获取 SMTP 服务
func (s *BaseService) SMTP() *smtp.Repository {
	return s.deps.SMTP
}

// MinIO 获取 MinIO 服务
func (s *BaseService) MinIO() *minio.Repository {
	return s.deps.MinIO
}
