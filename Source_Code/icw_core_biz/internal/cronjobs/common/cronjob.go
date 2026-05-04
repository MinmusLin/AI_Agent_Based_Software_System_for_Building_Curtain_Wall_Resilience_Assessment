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
	RocketMQ *rocketmq.Repository
	MinIO    *minio.Repository
}

func NewDeps(config configs.Config, MySQL *mysql.Repository, Redis *redis.Repository, RocketMQ *rocketmq.Repository, MinIO *minio.Repository) *Deps {
	return &Deps{
		Config:   config,
		MySQL:    MySQL,
		Redis:    Redis,
		RocketMQ: RocketMQ,
		MinIO:    MinIO,
	}
}

// CronJob 提供所有定时任务共享的基础依赖
type CronJob struct {
	ctx  context.Context
	deps *Deps
}

func NewCronJob(ctx context.Context, deps *Deps) *CronJob {
	if ctx == nil {
		ctx = context.Background()
	}
	if deps == nil {
		deps = &Deps{}
	}
	return &CronJob{
		ctx:  ctx,
		deps: deps,
	}
}

// CronJobFunc 定时任务执行函数
type CronJobFunc func(*CronJob) error

// Start 启动定时任务
func Start(ctx context.Context, deps *Deps, name, expression string, start CronJobFunc) error {
	if name == "" {
		CronFault("Failed to start cron job: %v", errors.New("name is required"))
	}
	if expression == "" {
		CronFault("Failed to start cron job: %v", errors.New("cron expression is required"))
	}
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err := parser.Parse(expression)
	if err != nil {
		CronFault("Failed to start cron job: %v", err.Error())
	}
	if deps == nil || deps.MySQL == nil || deps.Redis == nil || deps.MinIO == nil || deps.RocketMQ == nil {
		CronFault("Failed to start cron job: %v", errors.New("dependencies are required"))
		return nil
	}

	// 按 Cron 表达式执行定时任务
	job := NewCronJob(ctx, deps)
	return job.schedule(name, expression, func() error {
		if start == nil {
			return nil
		}
		return start(job)
	})
}

// schedule 按 Cron 表达式执行定时任务
func (c *CronJob) schedule(name, expression string, fn func() error) error {
	scheduler := cron.New(
		cron.WithSeconds(),
		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)),
	)

	if _, err := scheduler.AddFunc(expression, func() {
		_ = c.run(name, fn)
	}); err != nil {
		return err
	}

	scheduler.Start()

	go func() {
		<-c.Ctx().Done()
		stopCtx := scheduler.Stop()
		<-stopCtx.Done()
	}()

	return nil
}

// run 执行定时任务
func (c *CronJob) run(name string, fn func() error) (err error) {
	start := time.Now()
	defer func() {
		cronLog(name, start, err)
	}()

	if fn == nil {
		return nil
	}
	return fn()
}

// Ctx 获取上下文
func (c *CronJob) Ctx() context.Context {
	if c == nil || c.ctx == nil {
		return context.Background()
	}
	return c.ctx
}

// Config 获取服务配置
func (c *CronJob) Config() configs.Config {
	return c.deps.Config
}

// MySQL 获取 MySQL 服务
func (c *CronJob) MySQL() *mysql.Repository {
	return c.deps.MySQL
}

// Redis 获取 Redis 服务
func (c *CronJob) Redis() *redis.Repository {
	return c.deps.Redis
}

// RocketMQ 获取 RocketMQ 服务
func (c *CronJob) RocketMQ() *rocketmq.Repository {
	return c.deps.RocketMQ
}

// MinIO 获取 MinIO 服务
func (c *CronJob) MinIO() *minio.Repository {
	return c.deps.MinIO
}
