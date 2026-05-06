package icw_core_biz

import (
	"icw_common/consts"
	"icw_common/gen/core/biz"
	"icw_common/rpc"
)

// 编译期接口实现校验
var _ bizpb.ProjectDetectionServiceClient = (*Client)(nil)

// Client icw.core.biz gRPC Client
type Client struct {
	*rpc.Client
	bizpb.ProjectDetectionServiceClient
}

// NewClient 创建 icw.core.biz gRPC Client
func NewClient(addr string) (*Client, error) {
	baseClient, err := rpc.NewClient(consts.CoreBizPSM, addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		Client:                        baseClient,
		ProjectDetectionServiceClient: bizpb.NewProjectDetectionServiceClient(baseClient.Conn()),
	}, nil
}
