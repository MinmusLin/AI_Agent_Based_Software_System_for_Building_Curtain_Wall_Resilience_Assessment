package main

import (
	"context"
	"database/sql"
	"net"
	"net/rpc"

	goredis "github.com/redis/go-redis/v9"

	"icw_core_biz/configs"
	"icw_core_biz/consts"
	"icw_core_biz/internal/cronjobs"
	cronjobCommon "icw_core_biz/internal/cronjobs/common"
	"icw_core_biz/internal/services"
	serviceCommon "icw_core_biz/internal/services/common"
	"icw_core_biz/repositories/minio"
	"icw_core_biz/repositories/mysql"
	"icw_core_biz/repositories/redis"
	"icw_core_biz/repositories/rocketmq"
	"icw_core_biz/repositories/smtp"
	"icw_core_biz/utils"
)

// main icw.core.biz 服务入口
func main() {
	utils.LogInfo(consts.LogScopeInit, "", "Initializing service %s...", consts.CoreBizPSM)

	ctx := context.Background()

	// 加载服务配置
	configs.LoadDotEnv(".env")
	cfg, err := configs.Load()
	if err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to load config: %v", err)
	}
	utils.LogInfo(consts.LogScopeInit, "", "Config initialized successfully:\n%s", utils.FormatEnvConfig(cfg))

	// 初始化 MySQL
	dataMySQL, err := sql.Open("mysql", cfg.MySQLDSN)
	if err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to connect to MySQL: %v", err)
	}
	if err := dataMySQL.PingContext(ctx); err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to connect to MySQL: %v", err)
	}
	utils.LogInfo(consts.LogScopeInit, "", "MySQL initialized successfully")

	// 初始化 Redis
	dataRedis := goredis.NewClient(&goredis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	if err := dataRedis.Ping(ctx).Err(); err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to connect to Redis: %v", err)
	}
	utils.LogInfo(consts.LogScopeInit, "", "Redis initialized successfully")

	// 初始化 MinIO
	dataMinIO, err := minio.NewClient(cfg)
	if err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to connect to MinIO: %v", err)
	}
	bucketExists, err := dataMinIO.BucketExists(ctx, cfg.MinIOBucket)
	if err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to connect to MinIO: %v", err)
	}
	if !bucketExists {
		utils.LogFatal(consts.LogScopeInit, "Failed to find MinIO bucket: %s", cfg.MinIOBucket)
	}
	utils.LogInfo(consts.LogScopeInit, "", "MinIO initialized successfully")

	// 初始化 RocketMQ 生产者
	dataRocketMQ, err := rocketmq.NewProducer(cfg)
	if err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to connect to RocketMQ: %v", err)
	}
	defer func() {
		_ = dataRocketMQ.Shutdown()
	}()
	rocketmq.MQInfo("RocketMQ producer starts running")

	// 注册 RPC 服务
	services.RegisterRPCServices(ctx, serviceCommon.NewDeps(
		cfg,
		mysql.NewRepository(dataMySQL),
		redis.NewRepository(dataRedis),
		rocketmq.NewRepository(dataRocketMQ, cfg.RocketMQProjectEventTopic),
		minio.NewRepository(dataMinIO, cfg.MinIOBucket),
		smtp.NewRepository(cfg),
	))

	// 启动定时任务
	cronjobs.Start(ctx, cronjobCommon.NewDeps(
		cfg,
		mysql.NewRepository(dataMySQL),
		redis.NewRepository(dataRedis),
		rocketmq.NewRepository(dataRocketMQ, cfg.RocketMQProjectEventTopic),
		minio.NewRepository(dataMinIO, cfg.MinIOBucket),
	))

	// 运行 icw.core.biz 服务
	serviceCommon.RpcInfo("icw.core.biz service starts running on %s", cfg.CoreBizAddr)
	listener, err := net.Listen("tcp", cfg.CoreBizAddr)
	if err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to run icw.core.biz service: %v", err)
	}
	rpc.Accept(listener)
}
