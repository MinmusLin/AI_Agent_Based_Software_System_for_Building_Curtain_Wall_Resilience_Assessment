import { http } from '@/api/http';
import type { ApiEnvelope } from '@/constants/common';
import type {
  GetImageDetectionResultRequest,
  GetImageDetectionResultResponse,
  GetProjectDetectionTasksResponse,
  RetryProjectDetectionResponse,
  StartProjectDetectionResponse,
} from '@/gen/core/api/project_detection';

// 获取图像检测结果
// @router /project/detection/result [GET]
export async function getImageDetectionResult(
  request: GetImageDetectionResultRequest,
): Promise<GetImageDetectionResultResponse> {
  const { data } = await http.get<ApiEnvelope<GetImageDetectionResultResponse>>('/project/detection/result', {
    params: {
      image_uuid: request.image_uuid,
      project_id: request.project_id,
    },
  });
  return data.data;
}

// 获取项目智能检测任务列表
// @router /project/detection/list [GET]
export async function getProjectDetectionTasks(projectId: string): Promise<GetProjectDetectionTasksResponse> {
  const { data } = await http.get<ApiEnvelope<GetProjectDetectionTasksResponse>>('/project/detection/list', {
    params: {
      project_id: projectId,
    },
  });
  return data.data;
}

// 启动项目智能检测
// @router /project/detection/start [POST]
export async function startProjectDetection(projectId: string): Promise<StartProjectDetectionResponse> {
  const { data } = await http.post<ApiEnvelope<StartProjectDetectionResponse>>('/project/detection/start', {
    project_id: projectId,
  });
  return data.data;
}

// 重试项目智能检测
// @router /project/detection/retry [POST]
export async function retryProjectDetection(projectId: string): Promise<RetryProjectDetectionResponse> {
  const { data } = await http.post<ApiEnvelope<RetryProjectDetectionResponse>>('/project/detection/retry', {
    project_id: projectId,
  });
  return data.data;
}
