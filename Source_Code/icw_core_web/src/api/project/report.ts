import { http } from '@/api/http';
import type { ApiEnvelope } from '@/constants/common';
import type {
  GetProjectReportResponse,
  RetryProjectReportRequest,
  RetryProjectReportResponse,
} from '@/gen/core/api/project_report';

// 获取项目评估报告
// @router /project/report/detail [GET]
export async function getProjectReport(projectId: string): Promise<GetProjectReportResponse> {
  const { data } = await http.get<ApiEnvelope<GetProjectReportResponse>>('/project/report/detail', {
    params: {
      project_id: projectId,
    },
  });
  return data.data;
}

// 重试项目评估报告生成
// @router /project/report/retry [POST]
export async function retryProjectReport(payload: RetryProjectReportRequest): Promise<RetryProjectReportResponse> {
  const { data } = await http.post<ApiEnvelope<RetryProjectReportResponse>>('/project/report/retry', payload);
  return data.data;
}
