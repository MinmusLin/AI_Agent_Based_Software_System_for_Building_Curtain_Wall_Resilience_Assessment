package project_report

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_common/rpc"

	"icw_core_api/rpc/icw_core_biz"
	"icw_core_api/utils"
)

// GetProjectReport 获取项目评估报告
func GetProjectReport(ctx context.Context, client *icw_core_biz.Client, req *bizpb.GetProjectReportRequest, resp *bizpb.GetProjectReportResponse) error {
	return rpc.CallGRPC[bizpb.GetProjectReportRequest, bizpb.GetProjectReportResponse](ctx, client, req, resp, client.ProjectReport().GetProjectReport, rpc.WithRequestIdResolver(utils.GetXRequestId))
}
