package icw_core_biz

import (
	"icw_activity_reasoning/rpc/common"
	"icw_common/consts"
	"icw_common/gen/core/biz"
)

// Client icw.core.biz gRPC Client
type Client struct {
	*common.Client
	bizpb.ProjectDetectionServiceClient
}

// NewClient 创建 icw.core.biz gRPC Client
func NewClient(addr string) (*Client, error) {
	baseClient, err := common.NewClient(consts.CoreBizPSM, addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		Client:                        baseClient,
		ProjectDetectionServiceClient: bizpb.NewProjectDetectionServiceClient(baseClient.Conn()),
	}, nil
}
