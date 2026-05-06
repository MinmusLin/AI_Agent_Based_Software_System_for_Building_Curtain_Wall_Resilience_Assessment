package icw_activity_classification

import (
	"icw_common/consts"
	"icw_common/gen/activity/classification"
	"icw_core_biz/rpc/common"
)

// 编译期接口实现校验
var _ classificationpb.ClassificationServiceClient = (*Client)(nil)

// Client icw.activity.classification gRPC Client
type Client struct {
	*common.Client
	classificationpb.ClassificationServiceClient
}

// NewClient 创建 icw.activity.classification gRPC Client
func NewClient(addr string) (*Client, error) {
	baseClient, err := common.NewClient(consts.ActivityClassificationPSM, addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		Client:                      baseClient,
		ClassificationServiceClient: classificationpb.NewClassificationServiceClient(baseClient.Conn()),
	}, nil
}
