package icw_activity_summary

import (
	"icw_common/consts"
	"icw_common/gen/activity/summary"
	"icw_common/rpc"
)

// 编译期接口实现校验
var _ summarypb.SummaryServiceClient = (*Client)(nil)

// Client icw.activity.summary gRPC Client
type Client struct {
	*rpc.Client
	summarypb.SummaryServiceClient
}

// NewClient 创建 icw.activity.summary gRPC Client
func NewClient(addr string) (*Client, error) {
	baseClient, err := rpc.NewClient(consts.ActivitySummaryPSM, addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		Client:               baseClient,
		SummaryServiceClient: summarypb.NewSummaryServiceClient(baseClient.Conn()),
	}, nil
}
