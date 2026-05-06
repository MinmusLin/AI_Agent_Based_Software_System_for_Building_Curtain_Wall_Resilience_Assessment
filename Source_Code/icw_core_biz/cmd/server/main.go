package main

import (
	"context"
	"database/sql"
	"net"

	goredis "github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	"icw_common/consts"
	"icw_common/env"
	"icw_common/utils"
	"icw_core_biz/configs"
	"icw_core_biz/internal/cronjobs"
	cronjobCommon "icw_core_biz/internal/cronjobs/common"
	"icw_core_biz/internal/services"
	serviceCommon "icw_core_biz/internal/services/common"
	"icw_core_biz/repositories/minio"
	"icw_core_biz/repositories/mysql"
	"icw_core_biz/repositories/redis"
	"icw_core_biz/repositories/rocketmq"
	"icw_core_biz/repositories/smtp"
	"icw_core_biz/rpc/rpc_activity_classification"
	"icw_core_biz/rpc/rpc_activity_reasoning"
	"icw_core_biz/rpc/rpc_activity_summary"
)

// main icw.core.biz 服务入口
func main() {
	utils.LogInfo(consts.LogScopeInit, "", "Initializing service %s...", consts.CoreBizPSM)

	ctx := context.Background()

	// 加载服务配置
	if err := env.LoadDotEnv(".env"); err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to load .env file: %v", err)
	}
	cfg, err := configs.Load()
	if err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to load config: %v", err)
	}
	utils.LogInfo(consts.LogScopeInit, "", "Config initialized successfully:\n%s", env.FormatEnvConfig(cfg))

	// 初始化 MySQL
	dataMySQL, err := sql.Open("mysql", mysql.MySQLDSN(cfg))
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

	// 初始化 icw.activity.classification 服务
	activityClassificationClient, err := rpc_activity_classification.NewClient(cfg.ActivityClassificationAddr)
	if err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to initialize icw.activity.classification service: %v", err)
	}
	utils.LogInfo(consts.LogScopeRPC, consts.LogColorBoldGreen, "RPC service icw.activity.classification initialized successfully")

	// 初始化 icw.activity.reasoning 服务
	activityReasoningClient, err := rpc_activity_reasoning.NewClient(cfg.ActivityReasoningAddr)
	if err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to initialize icw.activity.reasoning service: %v", err)
	}
	utils.LogInfo(consts.LogScopeRPC, consts.LogColorBoldGreen, "RPC service icw.activity.reasoning initialized successfully")

	// 初始化 icw.activity.summary 服务
	activitySummaryClient, err := rpc_activity_summary.NewClient(cfg.ActivitySummaryAddr)
	if err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to initialize icw.activity.summary service: %v", err)
	}
	utils.LogInfo(consts.LogScopeRPC, consts.LogColorBoldGreen, "RPC service icw.activity.summary initialized successfully")

	defer func() {
		_ = activityClassificationClient.Close()
		_ = activityReasoningClient.Close()
		_ = activitySummaryClient.Close()
	}()

	// 启动定时任务
	cronjobs.Start(ctx, cronjobCommon.NewDeps(
		cfg,
		mysql.NewRepository(dataMySQL),
		redis.NewRepository(dataRedis),
		rocketmq.NewRepository(dataRocketMQ, cfg.RocketMQProjectEventTopic),
		minio.NewRepository(dataMinIO, cfg.MinIOBucket),
	))

	// 注册 gRPC 服务
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(utils.GRPCUnaryServerInterceptor(consts.LogScopeRPC, consts.LogColorBoldGreen)))
	services.RegisterRPCServices(ctx, serviceCommon.NewDeps(
		cfg,
		mysql.NewRepository(dataMySQL),
		redis.NewRepository(dataRedis),
		rocketmq.NewRepository(dataRocketMQ, cfg.RocketMQProjectEventTopic),
		minio.NewRepository(dataMinIO, cfg.MinIOBucket),
		smtp.NewRepository(cfg),
		activityClassificationClient,
		activityReasoningClient,
		activitySummaryClient,
	), grpcServer)

	// 运行 icw.core.biz 服务
	serviceCommon.RpcInfo("icw.core.biz service starts running on %s", cfg.CoreBizAddr)
	listener, err := net.Listen("tcp", cfg.CoreBizAddr)
	if err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to run icw.core.biz service: %v", err)
	}
	if err := grpcServer.Serve(listener); err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to run icw.core.biz service: %v", err)
	}
}
