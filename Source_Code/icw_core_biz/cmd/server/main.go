package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/rpc"

	"github.com/redis/go-redis/v9"

	"icw_core_biz/configs"
	"icw_core_biz/internal/auth"
	"icw_core_biz/internal/user"
	"icw_core_biz/repositories/minio"
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
	dataRedis := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	if err := dataRedis.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// 初始化 MinIO
	dataMinIO, err := minio.NewRepository(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MinIO: %v", err)
	}
	if err := dataMinIO.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to connect to MinIO: %v", err)
	}

	// 注册 RPC 服务
	userService := user.NewService(cfg, dataMinIO)
	if err := rpc.RegisterName("UserService", userService); err != nil {
		log.Fatalf("Failed to register user rpc service: %v", err)
	}
	authService := auth.NewService(cfg, dataMySQL, dataRedis, userService)
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
