package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/rpc"

	goredis "github.com/redis/go-redis/v9"

	"icw_core_biz/configs"
	"icw_core_biz/internal/services/auth"
	"icw_core_biz/internal/services/common"
	"icw_core_biz/internal/services/user"
	"icw_core_biz/repositories/minio"
	"icw_core_biz/repositories/mysql"
	"icw_core_biz/repositories/redis"
	"icw_core_biz/repositories/smtp"
)

// main icw.core.biz 服务入口
func main() {
	// 加载服务配置
	configs.LoadDotEnv(".env")
	cfg, err := configs.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化 MySQL
	dataMySQL, err := sql.Open("mysql", cfg.MySQLDSN)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	if err := dataMySQL.Ping(); err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	// 初始化 Redis
	dataRedis := goredis.NewClient(&goredis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	if err := dataRedis.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// 初始化 MinIO
	dataMinIO, err := minio.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MinIO: %v", err)
	}
	bucketExists, err := dataMinIO.BucketExists(context.Background(), cfg.MinIOBucket)
	if err != nil {
		log.Fatalf("Failed to connect to MinIO: %v", err)
	}
	if !bucketExists {
		log.Fatalf("Failed to find MinIO bucket: %s", cfg.MinIOBucket)
	}

	// 创建 RPC Service 的公共依赖集合
	serviceDeps := common.NewDeps(
		cfg,
		mysql.NewRepository(dataMySQL),
		redis.NewRepository(dataRedis),
		smtp.NewRepository(cfg),
		minio.NewRepository(dataMinIO, cfg.MinIOBucket),
	)

	// 注册 RPC 服务
	userService := user.NewService(serviceDeps)
	if err := rpc.RegisterName("UserService", userService); err != nil {
		log.Fatalf("Failed to register user rpc service: %v", err)
	}
	authService := auth.NewService(serviceDeps)
	if err := rpc.RegisterName("AuthService", authService); err != nil {
		log.Fatalf("Failed to register auth rpc service: %v", err)
	}

	// 运行 icw.core.biz 服务
	log.Printf("icw.core.biz service starts running on %s", cfg.CoreBizAddr)
	listener, err := net.Listen("tcp", cfg.CoreBizAddr)
	if err != nil {
		log.Fatalf("Failed to run icw.core.biz service: %v", err)
	}
	rpc.Accept(listener)
}
