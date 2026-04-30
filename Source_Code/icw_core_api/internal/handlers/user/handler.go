package user

import (
	"net/rpc"

	"icw_core_api/configs"
)

// Handler 用户业务 Handler
type Handler struct {
	cfg configs.Config
	rpc *rpc.Client
}

func NewHandler(cfg configs.Config, rpcClient *rpc.Client) *Handler {
	return &Handler{
		cfg: cfg,
		rpc: rpcClient,
	}
}
