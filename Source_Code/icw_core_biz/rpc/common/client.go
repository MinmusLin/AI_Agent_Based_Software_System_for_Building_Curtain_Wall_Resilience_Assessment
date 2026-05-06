package common

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"icw_common/consts"
	"icw_common/utils"
)

// Client 通用 gRPC Client
type Client struct {
	psm  string
	addr string
	conn *grpc.ClientConn
}

// NewClient 创建通用 gRPC Client
func NewClient(psm, addr string) (*Client, error) {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return nil, errors.New("grpc addr is required")
	}
	psm = strings.TrimSpace(psm)
	if psm == "" {
		return nil, errors.New("grpc psm is required")
	}
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(utils.GRPCUnaryClientInterceptor(consts.LogScopeRPC, psm)),
	)
	if err != nil {
		return nil, err
	}
	return &Client{
		psm:  psm,
		addr: addr,
		conn: conn,
	}, nil
}

// PSM 获取 gRPC 服务标识
func (c *Client) PSM() string {
	if c == nil || c.psm == "" {
		return "unknown"
	}
	return c.psm
}

// Addr 获取 gRPC 服务地址
func (c *Client) Addr() string {
	if c == nil {
		return ""
	}
	return c.addr
}

// Conn 获取 gRPC Client 连接
func (c *Client) Conn() grpc.ClientConnInterface {
	if c == nil {
		return nil
	}
	return c.conn
}

// Close 关闭 gRPC Client 连接
func (c *Client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

// Ready 判断 gRPC Client 连接是否可用
func (c *Client) Ready() bool {
	return c != nil && c.conn != nil
}

// Context 获取 gRPC 调用上下文
func Context(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return utils.AppendRequestIdToOutgoingContext(ctx, utils.RequestIdFromGRPCContext(ctx))
}
