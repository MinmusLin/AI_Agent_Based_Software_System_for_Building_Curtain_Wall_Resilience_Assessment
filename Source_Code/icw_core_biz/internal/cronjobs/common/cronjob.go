package common

import (
	"context"
	"errors"
	"time"

	"github.com/robfig/cron/v3"

	"icw_core_biz/configs"
	"icw_core_biz/repositories/minio"
	"icw_core_biz/repositories/mysql"
	"icw_core_biz/repositories/redis"
	"icw_core_biz/repositories/rocketmq"
)

// Deps 定时任务的公共依赖集合
type Deps struct {
	Config   configs.Config
	MySQL    *mysql.Repository
	Redis    *redis.Repository
	RocketMQ *rocketmq.Producer
	MinIO    *minio.Repository
}

func NewDeps(config configs.Config, MySQL *mysql.Repository, Redis *redis.Repository, RocketMQ *rocketmq.Producer, MinIO *minio.Repository) *Deps {
	return &Deps{
		Config:   config,
		MySQL:    MySQL,
		Redis:    Redis,
		RocketMQ: RocketMQ,
		MinIO:    MinIO,
	}
}

// BaseCronJob 提供所有定时任务共享的基础依赖
type BaseCronJob struct {
	ctx  context.Context
	deps *Deps
}

func NewBaseCronJob(ctx context.Context, deps *Deps) *BaseCronJob {
	if ctx == nil {
		ctx = context.Background()
	}
	if deps == nil {
		deps = &Deps{}
	}
	return &BaseCronJob{
		ctx:  ctx,
		deps: deps,
	}
}

// CronJob 定时任务执行接口
type CronJob interface {
	Start() (interface{}, error)
}

// JobFactory 定时任务构造函数
type JobFactory func(*BaseCronJob) CronJob

// Start 启动定时任务
func Start(ctx context.Context, deps *Deps, name, expression string, factory JobFactory) error {
	if name == "" {
		CronFatal("Failed to start cron job: %v", errors.New("name is required"))
	}
	if expression == "" {
		CronFatal("Failed to start cron job: %v", errors.New("cron expression is required"))
	}
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err := parser.Parse(expression)
	if err != nil {
		CronFatal("Failed to start cron job: %v", err.Error())
	}
	if deps == nil || deps.MySQL == nil || deps.Redis == nil || deps.MinIO == nil || deps.RocketMQ == nil {
		CronFatal("Failed to start cron job: %v", errors.New("dependencies are required"))
	}

	// 按 Cron 表达式执行定时任务
	baseJob := NewBaseCronJob(ctx, deps)
	job := CronJob(nil)
	if factory != nil {
		job = factory(baseJob)
	}
	return baseJob.schedule(name, expression, func() (interface{}, error) {
		if job == nil {
			return nil, nil
		}
		return job.Start()
	})
}

// schedule 按 Cron 表达式执行定时任务
func (j *BaseCronJob) schedule(name, expression string, fn func() (interface{}, error)) error {
	scheduler := cron.New(
		cron.WithSeconds(),
		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)),
	)

	if _, err := scheduler.AddFunc(expression, func() {
		_, _ = j.run(name, fn)
	}); err != nil {
		return err
	}

	scheduler.Start()

	go func() {
		<-j.Ctx().Done()
		stopCtx := scheduler.Stop()
		<-stopCtx.Done()
	}()

	return nil
}

// run 执行定时任务
func (j *BaseCronJob) run(name string, fn func() (interface{}, error)) (result interface{}, err error) {
	start := time.Now()
	defer func() {
		cronLog(name, start, result, err)
	}()

	if fn == nil {
		return nil, nil
	}
	return fn()
}

// Ctx 获取上下文
func (j *BaseCronJob) Ctx() context.Context {
	if j == nil || j.ctx == nil {
		return context.Background()
	}
	return j.ctx
}

// Config 获取服务配置
func (j *BaseCronJob) Config() configs.Config {
	if j == nil || j.deps == nil {
		return configs.Config{}
	}
	return j.deps.Config
}

// MySQL 获取 MySQL 服务
func (j *BaseCronJob) MySQL() *mysql.Repository {
	if j == nil || j.deps == nil {
		return nil
	}
	return j.deps.MySQL
}

// Redis 获取 Redis 服务
func (j *BaseCronJob) Redis() *redis.Repository {
	if j == nil || j.deps == nil {
		return nil
	}
	return j.deps.Redis
}

// RocketMQ 获取 RocketMQ 服务
func (j *BaseCronJob) RocketMQ() *rocketmq.Producer {
	if j == nil || j.deps == nil {
		return nil
	}
	return j.deps.RocketMQ
}

// MinIO 获取 MinIO 服务
func (j *BaseCronJob) MinIO() *minio.Repository {
	if j == nil || j.deps == nil {
		return nil
	}
	return j.deps.MinIO
}
