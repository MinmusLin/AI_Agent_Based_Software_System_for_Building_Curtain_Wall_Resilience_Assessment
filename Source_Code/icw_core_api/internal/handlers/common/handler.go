package common

import (
	"log"
	"net/rpc"
	"strings"

	"icw_core_api/configs"
	"icw_core_api/utils"
	"icw_core_biz/pkg/rpc_err"
)

const (
	// CoreBizPSM icw.core.biz 服务标识
	CoreBizPSM = "icw.core.biz"
)

// RPCClient 带服务标识的 RPC Client
type RPCClient struct {
	psm    string
	client *rpc.Client
}

// NewRPCClient 创建带服务标识的 RPC Client
func NewRPCClient(psm, address string) (*RPCClient, error) {
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &RPCClient{
		psm:    strings.TrimSpace(psm),
		client: client,
	}, nil
}

// PSM 获取 RPC 服务标识
func (c *RPCClient) PSM() string {
	if c == nil || c.psm == "" {
		return "unknown"
	}
	return c.psm
}

// Raw 获取原始 RPC Client
func (c *RPCClient) Raw() *rpc.Client {
	if c == nil {
		return nil
	}
	return c.client
}

// Close 关闭原始 RPC Client
func (c *RPCClient) Close() error {
	if c == nil || c.client == nil {
		return nil
	}
	return c.client.Close()
}

// Deps API Handler 的公共依赖集合
type Deps struct {
	Config        configs.Config
	CoreBizClient *RPCClient
}

// NewDeps 创建 API Handler 的公共依赖集合
func NewDeps(cfg configs.Config, coreBizClient *RPCClient) *Deps {
	return &Deps{
		Config:        cfg,
		CoreBizClient: coreBizClient,
	}
}

// BaseHandler 提供所有 API Handler 共享的基础依赖
type BaseHandler struct {
	deps *Deps
}

// NewBaseHandler 创建提供所有 API Handler 共享的基础依赖
func NewBaseHandler(deps *Deps) *BaseHandler {
	if deps == nil {
		deps = &Deps{}
	}
	return &BaseHandler{
		deps: deps,
	}
}

// Config 获取服务配置
func (h *BaseHandler) Config() configs.Config {
	return h.deps.Config
}

// CoreBizCall 调用 icw.core.biz RPC 服务
func (h *BaseHandler) CoreBizCall(method string, req interface{}, resp interface{}) error {
	if h == nil || h.deps == nil {
		return nil
	}
	return CallRPC(h.deps.CoreBizClient, method, req, resp)
}

// CallRPC RPC 服务通用调用
func CallRPC(client *RPCClient, method string, req interface{}, resp interface{}) error {
	psm := client.PSM()
	if client.Raw() == nil {
		err := rpc_err.InternalErrorDefault("rpc client is nil")
		log.Printf("[ERROR] Call %s %s failed, req: %s, resp: %s, err: %v", psm, method, utils.JSONF(req), utils.JSONF(resp), err)
		return err
	}
	err := client.Raw().Call(method, req, resp)
	if err != nil {
		log.Printf("[ERROR] Call %s %s failed, req: %s, resp: %s, err: %v", psm, method, utils.JSONF(req), utils.JSONF(resp), err)
		return err
	}
	if resp == nil {
		err = rpc_err.InternalErrorDefault("rpc response is nil")
		log.Printf("[ERROR] Call %s %s failed, req: %s, resp: %s, err: %v", psm, method, utils.JSONF(req), utils.JSONF(resp), err)
		return err
	}
	return nil
}
