package icw_core_biz

import (
	"icw_common/consts"
	"icw_common/gen/core/biz"
	"icw_common/rpc"
)

// Client icw.core.biz gRPC Client
type Client struct {
	*rpc.Client
	auth             bizpb.AuthServiceClient
	projectAssets    bizpb.ProjectAssetsServiceClient
	projectCore      bizpb.ProjectCoreServiceClient
	projectDetection bizpb.ProjectDetectionServiceClient
	projectProfile   bizpb.ProjectProfileServiceClient
	projectReport    bizpb.ProjectReportServiceClient
	projectReview    bizpb.ProjectReviewServiceClient
	socket           bizpb.SocketServiceClient
	user             bizpb.UserServiceClient
}

// NewClient 创建 icw.core.biz gRPC Client
func NewClient(address string) (*Client, error) {
	baseClient, err := rpc.NewClient(consts.CoreBizPSM, address)
	if err != nil {
		return nil, err
	}
	conn := baseClient.Conn()
	return &Client{
		Client:           baseClient,
		auth:             bizpb.NewAuthServiceClient(conn),
		user:             bizpb.NewUserServiceClient(conn),
		socket:           bizpb.NewSocketServiceClient(conn),
		projectCore:      bizpb.NewProjectCoreServiceClient(conn),
		projectProfile:   bizpb.NewProjectProfileServiceClient(conn),
		projectAssets:    bizpb.NewProjectAssetsServiceClient(conn),
		projectDetection: bizpb.NewProjectDetectionServiceClient(conn),
		projectReview:    bizpb.NewProjectReviewServiceClient(conn),
		projectReport:    bizpb.NewProjectReportServiceClient(conn),
	}, nil
}

// Auth 获取登录鉴权服务 gRPC Client
func (c *Client) Auth() bizpb.AuthServiceClient {
	if c == nil {
		return nil
	}
	return c.auth
}

// ProjectAssets 获取图像资产服务 gRPC Client
func (c *Client) ProjectAssets() bizpb.ProjectAssetsServiceClient {
	if c == nil {
		return nil
	}
	return c.projectAssets
}

// ProjectCore 获取项目核心服务 gRPC Client
func (c *Client) ProjectCore() bizpb.ProjectCoreServiceClient {
	if c == nil {
		return nil
	}
	return c.projectCore
}

// ProjectDetection 获取智能检测服务 gRPC Client
func (c *Client) ProjectDetection() bizpb.ProjectDetectionServiceClient {
	if c == nil {
		return nil
	}
	return c.projectDetection
}

// ProjectProfile 获取基础信息服务 gRPC Client
func (c *Client) ProjectProfile() bizpb.ProjectProfileServiceClient {
	if c == nil {
		return nil
	}
	return c.projectProfile
}

// ProjectReport 获取评估报告服务 gRPC Client
func (c *Client) ProjectReport() bizpb.ProjectReportServiceClient {
	if c == nil {
		return nil
	}
	return c.projectReport
}

// ProjectReview 获取人工复核服务 gRPC Client
func (c *Client) ProjectReview() bizpb.ProjectReviewServiceClient {
	if c == nil {
		return nil
	}
	return c.projectReview
}

// Socket 获取 WebSocket 连接票据服务 gRPC Client
func (c *Client) Socket() bizpb.SocketServiceClient {
	if c == nil {
		return nil
	}
	return c.socket
}

// User 获取用户业务服务 gRPC Client
func (c *Client) User() bizpb.UserServiceClient {
	if c == nil {
		return nil
	}
	return c.user
}
