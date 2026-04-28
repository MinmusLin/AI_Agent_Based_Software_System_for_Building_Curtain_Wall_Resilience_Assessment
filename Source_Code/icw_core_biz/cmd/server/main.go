package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/rpc"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"

	"icw_core_biz/configs"
	"icw_core_biz/internal/auth"
)

func main() {
	// Load configs
	configs.LoadDotEnv(".env")
	cfg := configs.Load()

	// Initialize MySQL
	dataMySQL, err := sql.Open("mysql", cfg.MySQLDSN)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	if err := dataMySQL.Ping(); err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	// Initialize Redis
	dataRedis := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	if err := dataRedis.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize RPC service
	authService := auth.NewService(cfg, dataMySQL, dataRedis)
	if err := rpc.RegisterName("AuthService", authService); err != nil {
		log.Fatalf("Failed to register auth rpc service: %v", err)
	}

	// Run icw.core.biz service
	log.Printf("icw.core.biz service starts running on %s", cfg.CoreBizAddr)
	listener, err := net.Listen("tcp", cfg.CoreBizAddr)
	if err != nil {
		log.Fatalf("Failed to run icw.core.biz service: %v", err)
	}
	rpc.Accept(listener)
}
