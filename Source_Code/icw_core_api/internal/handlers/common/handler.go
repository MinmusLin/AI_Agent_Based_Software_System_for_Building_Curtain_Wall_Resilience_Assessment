package common

import (
	"context"
	"net/rpc"
	"reflect"
	"strings"

	"icw_core_api/configs"
	"icw_core_api/utils"
	"icw_core_biz/consts"
	"icw_core_biz/pkg/dto"
	"icw_core_biz/pkg/rpc_err"
	bizUtils "icw_core_biz/utils"
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
func (h *BaseHandler) CoreBizCall(ctx context.Context, method string, req, resp interface{}) error {
	if h == nil || h.deps == nil {
		return nil
	}
	return CallRPC(ctx, h.deps.CoreBizClient, method, req, resp)
}

// CallRPC RPC 服务通用调用
func CallRPC(ctx context.Context, client *RPCClient, method string, req, resp interface{}) error {
	if ctx == nil {
		ctx = context.Background()
	}

	// 从请求上下文中获取请求 ID 并设置 RPC 元数据
	requestId := utils.GetRequestId(ctx)
	setRPCMeta(req, requestId)

	psm := client.PSM()
	if client.Raw() == nil {
		err := rpc_err.InternalErrorDefault("rpc client is nil")
		rpcErrorLog(requestId, psm, method, req, resp, err)
		return err
	}

	call := client.Raw().Go(method, req, resp, make(chan *rpc.Call, 1))
	select {
	case <-ctx.Done():
		ctxErr := ctx.Err()
		if ctxErr == nil {
			ctxErr = context.Canceled
		}
		err := rpc_err.InternalErrorDefault(ctxErr.Error())
		rpcErrorLog(requestId, psm, method, req, resp, err)
		return err
	case done := <-call.Done:
		if done.Error != nil {
			rpcErrorLog(requestId, psm, method, req, resp, done.Error)
			return done.Error
		}
	}

	if resp == nil {
		err := rpc_err.InternalErrorDefault("rpc response is nil")
		rpcErrorLog(requestId, psm, method, req, resp, err)
		return err
	}

	return nil
}

// setRPCMeta 设置 RPC 元数据
func setRPCMeta(req interface{}, requestId string) {
	requestId = strings.TrimSpace(requestId)
	if req == nil || requestId == "" {
		return
	}

	reqValue := reflect.ValueOf(req)
	if reqValue.Kind() != reflect.Ptr || reqValue.IsNil() {
		return
	}

	elem := reqValue.Elem()
	if elem.Kind() != reflect.Struct {
		return
	}

	metaField := elem.FieldByName("Meta")
	if !metaField.IsValid() || !metaField.CanSet() {
		return
	}

	meta := &dto.Meta{
		RequestId: requestId,
	}
	if metaField.Kind() == reflect.Ptr && metaField.Type() == reflect.TypeOf(meta) {
		metaField.Set(reflect.ValueOf(meta))
	}
}

// rpcErrorLog 输出 RPC 调用失败日志
func rpcErrorLog(requestId, psm, method string, req, resp, err interface{}) {
	bizUtils.LogError(consts.LogScopeHTTP, "[%s] Call %s %s failed, req: %s, resp: %s, err: %s",
		requestId,
		psm,
		method,
		bizUtils.JSONF(req),
		bizUtils.JSONF(resp),
		bizUtils.FormatErrorLog(err),
	)
}
