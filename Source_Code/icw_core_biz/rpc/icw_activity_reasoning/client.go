package icw_activity_reasoning

import (
	"icw_common/consts"
	"icw_common/gen/activity/reasoning"
	"icw_core_biz/rpc/common"
)

// 编译期接口实现校验
var _ reasoningpb.ReasoningServiceClient = (*Client)(nil)

// Client icw.activity.reasoning gRPC Client
type Client struct {
	*common.Client
	reasoningpb.ReasoningServiceClient
}

// NewClient 创建 icw.activity.reasoning gRPC Client
func NewClient(addr string) (*Client, error) {
	baseClient, err := common.NewClient(consts.ActivityReasoningPSM, addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		Client:                 baseClient,
		ReasoningServiceClient: reasoningpb.NewReasoningServiceClient(baseClient.Conn()),
	}, nil
}
